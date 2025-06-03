package models

import (
	"fmt"
	"rentacar/types"
	"rentacar/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var ProcessStatus *types.CloningWokgroupStatus

// *************************************************************************************************************************
// retorna lista de grupos
// *************************************************************************************************************************
func GetWorkgroupsForDatatable(request *types.DatatableParamsRequest, c *fiber.Ctx) (recordset *[]types.WorkgroupsResponse, err error) {

	// nao foi passado campo order by inicialmente
	if request.OrderBy == "" {
		request.OrderBy = "name"
		request.OrderDirection = "asc"
	}

	var fields string

	if request.Country == "usa" {
		fields =
			" id, name, ifnull(in_use, false) as in_use, ifnull(client_ip, '') as client_ip, ifnull(client_country, '') as client_country, " +
				" ifnull(client_city, '') as client_city, ifnull(client_loc, '') as client_loc, " +
				" date_format(updated_at, '%m/%d/%Y %l:%i - %p') as updated_at, " +
				" date_format(deleted_at, '%m/%d/%Y %l:%i - %p') as deleted_at, ifnull(active, false) as active "
	}

	if request.Country == "brazil" {
		fields =
			" id, name, ifnull(in_use, false) as in_use, ifnull(client_ip, '') as client_ip, ifnull(client_country, '') as client_country, " +
				" ifnull(client_city, '') as client_city, ifnull(client_loc, '') as client_loc, " +
				" date_format(updated_at, '%d/%m/%Y %H:%i') as updated_at, " +
				" date_format(deleted_at, '%d/%m/%Y %H:%i') as deleted_at, ifnull(active, false) as active "
	}

	var where string
	// se o grupo atual no front end nao é 'admin', exibe somente grupos disponiveis
	// nao deixa usuario comum ver grupos abertos por outras pessoas para que ele nao possa ingressar
	// usuario comum pode compartilhar o proprio grupo com outras pessoas
	if strings.ToLower(strings.TrimSpace(c.Params("workgroup"))) != "admin" {
		where = " ifnull(in_use, false) = false "
	}

	if request.SearchtTxt != "" {
		// usa todas as colunas visiveis no datatable do front end
		request.SearchtTxt = "%" + request.SearchtTxt + "%"
		where += utils.ConcatWhere(where, fmt.Sprintf("  (name like '%v') ", request.SearchtTxt))
	}
	if request.OnlyActiveOrInactiveRecords == "active" {
		// so ativos
		where += utils.ConcatWhere(where, " ifnull(active, false) = true ")
	}
	if request.OnlyActiveOrInactiveRecords == "inactive" {
		// so inativos
		where += utils.ConcatWhere(where, " ifnull(active, false) = false ")
	}

	request.OrderBy += " " + request.OrderDirection

	// busca somente os campos necessarios
	// necessario que seja unscoped (incluir registros apagados logicamente (deleted_at <> null))
	// porque qdo grupo é resetado, ele é apagado logicamente mas é mantido para log
	if err = Db.Model(&types.Workgroups{}).Unscoped().Order(request.OrderBy).Select(fields).Where(where).Scan(&recordset).Error; err != nil {
		return nil, err
	}

	return recordset, nil
}

// *************************************************************************************************************************
// 1. escolhe aleatoriamente um grupo dentro dos grupos disponiveis na tabelas 'workgroups' (apliacao executada 1a vez)
// ou
// 2. reseta os dado do grupo atual (usuario pediu para resetar dados do grupo)
// ou
// 3. desativa o grupo atual e escolhe outro grupo aleatoriamente (usuario pediu para sortear outro grupo)
// *************************************************************************************************************************
func GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother(c *fiber.Ctx) (workgroup string, err error) {

	var _workgroup_ *types.WorkgroupResponse

	var _destination_workgroupName string // nome do grupo que sera gerado ou cujos dados serao resetados
	var _old_workgroupName string         // nome do grupo que esta sendo substituido (whatToDo == "another")

	whatToDo := strings.ToLower(strings.TrimSpace(c.Params("action")))

	// se front end sendo executada pela 1a vez ou se usuario pediu para sortear outro grupo, sorteia um grupo e prepara seus dados/arquivos
	if whatToDo == "generate" || whatToDo == "another" {
		// rand() = sorteia um grupo que nao foi utilizado ainda
		if err = Db.Model(&types.Workgroups{}).Where(" ifnull(in_use, false) = false").Select("name, id").Limit(1).Order("RAND()").Scan(&_workgroup_).Error; err != nil {
			return "", err
		}

		// dados do cliente
		var _clientInfo types.ClientInfo
		c.BodyParser(&_clientInfo)

		// marca o grupo sorteado como indisponivel e atribui IP que o pegou
		if err = Db.Model(&types.Workgroups{}).Where("name = ? ", _workgroup_.Name).Updates(
			types.Workgroups{
				Active:         true,
				InUse:          1,
				ClientIp:       _clientInfo.Ip,
				ClientCity:     _clientInfo.City,
				ClientCountry:  _clientInfo.Country,
				ClientHostname: _clientInfo.Hostname,
				ClientLoc:      _clientInfo.Loc,
				ClientOrg:      _clientInfo.Org,
				ClientPostal:   _clientInfo.Postal,
				ClientRegion:   _clientInfo.Region,
				ClientTimezone: _clientInfo.Timezone,
			}).Error; err != nil {

			return "", err
		}

		_destination_workgroupName = _workgroup_.Name
		_old_workgroupName = c.Params("workgroup") // em caso de substituicao (another), os dados/arquivos deste grupo serao excluidos
	}

	// front end ja tem grupo definido/dados.. mas usuario pediu para resetar seus dados/arquivos ou usuario pediu para sortear outro grupo
	if whatToDo == "reset" {
		_destination_workgroupName = c.Params("workgroup")
	}

	// copia todos registros do grupo 'Admin' para novos registros vinculados somente ao grupo recem criado
	if err = CopyAdminRecordsToWorkgroup(_destination_workgroupName, whatToDo, _old_workgroupName); err != nil {
		return "", err
	}

	return "__success__|" + _destination_workgroupName, err
}

// *************************************************************************************************************************
// faz cópia de todos os registros do grupo 'Admin' para determinado grupo
// *************************************************************************************************************************
func CopyAdminRecordsToWorkgroup(_destination_workgroupName string, whatToDo string, _old_workgroupName string) (err error) {

	ProcessStatus = &types.CloningWokgroupStatus{
		Status:          "initiating_workgroup_records",
		PercentReady:    "10",
		ChosenWorkgroup: _destination_workgroupName,
	}

	// vai conter os arquivos de imagem (carro, logotipo fabricante) , que em caso de reset/substituicao, devem ser apagados
	var carFilesDelete []*types.FilesToDeleteInAWS
	carFilesToDelete := false

	var manufacturerFilesDelete []*types.FilesToDeleteInAWS
	manufacturerFilesToDelete := false

	var sql string

	// se der algo de errado, cancela todas as queries ja feitas
	trans := Db.Begin()

	time.Sleep(1000 * time.Millisecond)

	// se usuario pediu para resetar dados de um grupo ou pediu para sortear outro grupo, apaga fisicamente registros de carros, termos, etc
	// e apaga LOGICAMENTE (deleted_at) o registro do grupo e cria um novo - faz isso para manter log do IP que utilizou o grupo por ultimo

	if (whatToDo == "reset" || whatToDo == "another") && _destination_workgroupName != "Admin" {

		var where string
		var _previous_workgroupName string

		// se é um reset de dados, o grupo esta sendo substituido por ele mesmo
		if whatToDo == "reset" {
			_previous_workgroupName = _destination_workgroupName
		}

		// se é um sorteio de novo grupo, um grupo some (_old_workgroupName) e outro sera usado (_destination_workgroupName)
		if whatToDo == "another" {
			_previous_workgroupName = _old_workgroupName
		}

		// obtem nome dos arquivos (imagem) de carros para mais a frente exclui los via ftp
		sql = fmt.Sprintf(" concat('%v_car_', LPAD(id, 6, '0'), '.png') as filename ", _previous_workgroupName)
		where = fmt.Sprintf(" workgroup = '%v' ", _previous_workgroupName)

		if err = trans.Model(&types.Cars{}).Where(where).Select(sql).Scan(&carFilesDelete).Error; err != nil {
			trans.Rollback()
			return err
		}

		for range carFilesDelete {
			carFilesToDelete = true
		}

		// obtem nome dos arquivos (logotipos) de fabricantes para mais a frente exclui los via ftp
		sql = fmt.Sprintf(" concat('%v_manufacturer_', LPAD(id, 6, '0'), '.png') as filename ", _previous_workgroupName)
		where = fmt.Sprintf(" workgroup = '%v' ", _previous_workgroupName)

		if err = trans.Model(&types.Manufacturers{}).Where(where).Select(sql).Scan(&manufacturerFilesDelete).Error; err != nil {
			trans.Rollback()
			return err
		}
		for range manufacturerFilesDelete {
			manufacturerFilesToDelete = true
		}

		var _workgroup_ *types.WorkgroupResponse

		// obtem ID do grupo que sera resetado/substituido
		where = fmt.Sprintf("name='%v' ", _previous_workgroupName)
		if err = Db.Model(&types.Workgroups{}).Where(where).Select("name, id").Scan(&_workgroup_).Error; err != nil {
			trans.Rollback()
			return err
		}

		// apaga somente logicamente o registro antigo do grupo, para manter log do IP que usou por ultimo
		if err = trans.Exec(fmt.Sprintf("update workgroups set deleted_at=now() where id=%v ", _workgroup_.Id)).Error; err != nil {
			trans.Rollback()
			return err
		}

		// se pediu para resetar, clona registro do grupo mas com dados atualizados
		if whatToDo == "reset" {
			sql = fmt.Sprintf(
				"insert into workgroups(name, in_use, database_changes_amount, active,  created_at, updated_at, client_ip, client_city, client_region, client_country, "+
					" client_hostname, client_loc, client_org, client_postal, client_timezone)"+
					"select name, true, 0, true, now(), now(), client_ip, client_city, client_region, client_country, "+
					"				client_hostname, client_loc, client_org, client_postal, client_timezone "+
					"from workgroups where id = %v  ", _workgroup_.Id)

			if err = trans.Exec(sql).Error; err != nil {
				trans.Rollback()
				return err
			}
		}

		// se pediu para sortear outro grupo, faz copia do registro do grupo que foi excluido acima para deixa lo disponivel novamente (in_use= false)
		if whatToDo == "another" {
			sql = fmt.Sprintf(
				"insert into workgroups(name, in_use, database_changes_amount, active,  created_at, updated_at) "+
					"select '%v', false, 0, true, now(), now() ", _previous_workgroupName)

			if err = trans.Exec(sql).Error; err != nil {
				trans.Rollback()
				return err
			}
		}

		// apaga fisicamente os demais registros poruqe nao ha necessidade de manter log deles
		if err = trans.Exec(fmt.Sprintf("delete from cars where workgroup='%v' ", _previous_workgroupName)).Error; err != nil {
			trans.Rollback()
			return err
		}
		if err = trans.Exec(fmt.Sprintf("delete from manufacturers where workgroup='%v' ", _previous_workgroupName)).Error; err != nil {
			trans.Rollback()
			return err
		}
		if err = trans.Exec(fmt.Sprintf("delete from terms where workgroup='%v' ", _previous_workgroupName)).Error; err != nil {
			trans.Rollback()
			return err
		}
		if err = trans.Exec(fmt.Sprintf("delete from bookings where workgroup='%v' ", _previous_workgroupName)).Error; err != nil {
			trans.Rollback()
			return err
		}
		if err = trans.Exec(fmt.Sprintf("delete from notifications where workgroup='%v' ", _previous_workgroupName)).Error; err != nil {
			trans.Rollback()
			return err
		}

	}

	//****************************************************************************************
	// clona registros usando como base os dados do grupo 'Admin' (dados default)
	//****************************************************************************************

	//********************************************************
	// clona registros de carros
	//********************************************************

	ProcessStatus.Status = "preparing_car_records"
	ProcessStatus.PercentReady = "20"

	time.Sleep(1000 * time.Millisecond)

	sql = fmt.Sprintf(
		"insert into cars(active, country, year, manufacturer_id, name, odometer, mpg, cylinders, transmission_manual, hp, doors, cc, rental_price, created_at, "+
			"     updated_at, original_id, workgroup) "+
			"select active, country, year, manufacturer_id, name, odometer, mpg, cylinders, transmission_manual, hp, doors, cc, rental_price, now(), now(), id, '%v' "+
			"from cars where workgroup='Admin' and deleted_at is null ", _destination_workgroupName)

	if err = trans.Exec(sql).Error; err != nil {
		trans.Rollback()
		return err
	}

	ProcessStatus.Status = "preparing_car_records"
	ProcessStatus.PercentReady = "30"

	time.Sleep(1000 * time.Millisecond)

	//********************************************************
	// clona registros de fabricantes
	//********************************************************
	ProcessStatus.Status = "preparing_manufacturers_records"
	ProcessStatus.PercentReady = "40"

	time.Sleep(1000 * time.Millisecond)

	sql = fmt.Sprintf(
		"insert into manufacturers (active, name, created_at, updated_at, original_id, workgroup) "+
			"select active, name, now(), now(), id, '%v'  "+
			"from manufacturers where workgroup='Admin' and deleted_at is null", _destination_workgroupName)

	if err = trans.Exec(sql).Error; err != nil {
		trans.Rollback()
		return err
	}

	// atualiza/corrige ID do fabricante nos registros recem clonados de carros
	sql = fmt.Sprintf(
		"update cars, manufacturers "+
			"set cars.manufacturer_id = manufacturers.id "+
			"where cars.workgroup  = manufacturers.workgroup and cars.manufacturer_id = manufacturers.original_id and cars.workgroup = '%v' ", _destination_workgroupName)

	if err = trans.Exec(sql).Error; err != nil {
		trans.Rollback()
		return err
	}

	//**********************************************************************
	// clona arquivos de imagem de fabricantes (logos) e carros
	//**********************************************************************
	ProcessStatus.Status = "cloning_image_files"
	ProcessStatus.PercentReady = "41"
	time.Sleep(2000 * time.Millisecond)

	// relaciona arquivos de imagem de carros que devem ser copiados do grupo ADMIN para o grupo recem criado
	var filesToCopyCars []*types.FilesToCopyFromAdmin

	sql = fmt.Sprintf(" concat('Admin_car_', LPAD(original_id, 6, '0'), '.png') as AdminFilename , "+
		" concat('%v_car_', LPAD(id, 6, '0'), '.png') as NewGroupFilename ", _destination_workgroupName)
	where := fmt.Sprintf(" workgroup = '%v' ", _destination_workgroupName)

	if err = trans.Model(&types.Cars{}).Where(where).Select(sql).Scan(&filesToCopyCars).Error; err != nil {
		trans.Rollback()
		return err
	}

	// relaciona arquivos de imagem de carros que devem ser copiados do grupo ADMIN para o grupo recem criado
	var filesToCopyManufacturers []*types.FilesToCopyFromAdmin

	sql = fmt.Sprintf(" concat('Admin_manufacturer_', LPAD(original_id, 6, '0'), '.png') as AdminFilename , "+
		" concat('%v_manufacturer_', LPAD(id, 6, '0'), '.png') as NewGroupFilename ", _destination_workgroupName)
	where = fmt.Sprintf(" workgroup = '%v' ", _destination_workgroupName)

	if err = trans.Model(&types.Manufacturers{}).Where(where).Select(sql).Scan(&filesToCopyManufacturers).Error; err != nil {
		trans.Rollback()
		return err
	}

	// agrupa arquivos imagem de carros de de fabricantes (logotipos)
	var filesToCopy []*types.FilesToCopyFromAdmin
	for _, file := range filesToCopyCars {
		fileadd := types.FilesToCopyFromAdmin{AdminFilename: file.AdminFilename, NewGroupFilename: file.NewGroupFilename}
		filesToCopy = append(filesToCopy, &fileadd)
	}

	for _, file := range filesToCopyManufacturers {
		fileadd := types.FilesToCopyFromAdmin{AdminFilename: file.AdminFilename, NewGroupFilename: file.NewGroupFilename}
		filesToCopy = append(filesToCopy, &fileadd)
	}

	// copia os arquivos para o S3
	err = utils.CopyAdminFilesInAWS_S3(filesToCopy, ProcessStatus, 41, 60)
	if err != nil {
		trans.Rollback()
		return err
	}

	//********************************************************
	// termos/expressoes em ingles/portugues
	//********************************************************

	ProcessStatus.Status = "preparing_terms_records"
	ProcessStatus.PercentReady = "60"

	time.Sleep(1000 * time.Millisecond)

	sql = fmt.Sprintf(
		"insert into terms(item, portuguese, english, description, active,  workgroup, created_at, updated_at) "+
			"select item, portuguese, english, description, active, '%v', now(), now() "+
			"from terms where workgroup='Admin' and deleted_at is null", _destination_workgroupName)

	if err = trans.Exec(sql).Error; err != nil {
		trans.Rollback()
		return err
	}

	// transacao sql ocorrey ok, se o grupo foi resetado ou substituido, apaga arquivos de imagem (carros, etc) antigos
	ProcessStatus.Status = "deleting_image_old_files"
	ProcessStatus.PercentReady = "70"

	time.Sleep(1000 * time.Millisecond)

	// agrupa arquivos para apagar
	var filesToDelete []*types.FilesToDeleteInAWS

	if carFilesToDelete || manufacturerFilesToDelete {

		for _, file := range carFilesDelete {
			fileadd := types.FilesToDeleteInAWS{Filename: file.Filename}
			filesToDelete = append(filesToDelete, &fileadd)
		}

		for _, file := range manufacturerFilesDelete {
			fileadd := types.FilesToDeleteInAWS{Filename: file.Filename}
			filesToDelete = append(filesToDelete, &fileadd)
		}
	}

	utils.DeleteFilesInAWS_S3(filesToDelete)

	ProcessStatus.Status = "cloning_process_almost_done"
	ProcessStatus.PercentReady = "80"

	time.Sleep(1000 * time.Millisecond)

	// avisa que a clonagenm foi concluida
	ProcessStatus.PercentReady = "100"
	ProcessStatus.Status = "cloning_process_done"

	time.Sleep(3000 * time.Millisecond)

	trans.Commit()

	// apos 5 segundos, reseta status e percentual da clonagem para que nao informe errado o front end
	time.AfterFunc(time.Duration(5)*time.Second, func() {
		ProcessStatus.Status = ""
		ProcessStatus.PercentReady = "0"
	})

	return nil
}

// *************************************************************************************************************************
// verifica se um nome de grupo passado pelo usuario existe e esta em uso, consequentemente seus dados podem ser acessados
// *************************************************************************************************************************
func AccessAnotherGroupData(c *fiber.Ctx) (workgroup string, err error) {

	var _workgroup_ *types.WorkgroupResponse

	where := fmt.Sprintf(" ifnull(in_use, false) = true and trim(name)=trim('%v') ", c.Params("workgroup"))

	query := Db.Model(&types.Workgroups{}).Where(where).Select("name, id").Limit(1).Scan(&_workgroup_)

	if query.Error != nil {
		return "", query.Error

	} else {

		// o grupo requisitado nao existe ou nao esta disponivel
		if query.RowsAffected == 0 {
			return "none|", nil

			// grupo encontrado e disponivel
		} else {
			return "__success__|" + _workgroup_.Name + "|" + strconv.Itoa(_workgroup_.Id), nil
		}

	}
}

// *************************************************************************************************************************
// obtem qtde de alteracoes feitas na base pelo grupo
// *************************************************************************************************************************
func WorkgroupReport(c *fiber.Ctx) (workgroup string, err error) {

	var _workgroup_ *types.WorkgroupReport

	where := fmt.Sprintf(" trim(name)=trim('%v') ", c.Params("workgroup"))

	query := Db.Model(&types.Workgroups{}).Where(where).Select("ifnull(database_changes_amount, 0) as database_changes_amount").Scan(&_workgroup_)

	if query.Error != nil {
		return "", query.Error

	} else {

		return "__success__|" + strconv.Itoa(_workgroup_.DatabaseChangesAmount), nil

	}
}

// *************************************************************************************************************************
// contabiliza 1 operacao (inclusao, edicao, exclusao) a mais feita pelo workgroup atual
// isso é necessario porque se o usuario pedir para resetar os dados, o front end so permite se houver pelo menos 1
// alteracao feita na base, pelo grupo
// *************************************************************************************************************************
func WorkgroupAddLog(workgroup string) error {

	try := Db.Exec(fmt.Sprintf("update workgroups set database_changes_amount = ifnull(database_changes_amount, 0) + 1 where name='%v'", workgroup))

	return try.Error
}
