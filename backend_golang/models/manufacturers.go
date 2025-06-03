package models

import (
	"fmt"
	"rentacar/types"
	"rentacar/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// *************************************************************************************************************************
// retorna lista de fabricantes de carros
// *************************************************************************************************************************
func GetManufacturersForDatatable(request *types.DatatableParamsRequest, c *fiber.Ctx) (recordset *[]types.ManufacturerResponse, err error) {

	// nao foi passado campo order by inicialmente
	if request.OrderBy == "" {
		request.OrderBy = "name"
		request.OrderDirection = "asc"
	}

	fields := " id , name, ifnull(active, 0) as active, concat(workgroup, '_manufacturer_', LPAD(id, 6, '0'), '.png') as manufacturer_logo  "

	var where string
	if request.SearchtTxt != "" {
		// usa todas as colunas visiveis no datatable do front end
		request.SearchtTxt = "%" + request.SearchtTxt + "%"
		where += utils.ConcatWhere(where, fmt.Sprintf(" (name like '%v') ", request.SearchtTxt))
	}
	if request.OnlyActiveOrInactiveRecords == "active" {
		// so ativos
		where += utils.ConcatWhere(where, " ifnull(active, false) = true ")
	}
	if request.OnlyActiveOrInactiveRecords == "inactive" {
		// so inativos
		where += utils.ConcatWhere(where, " ifnull(active, false) = false ")
	}
	where += utils.ConcatWhere(where, fmt.Sprintf(" workgroup = '%v' ", c.Params("workgroup")))

	request.OrderBy += " " + request.OrderDirection

	// busca somente os campos necessarios
	err = Db.Model(&types.Manufacturers{}).Order(request.OrderBy).Select(fields).Where(where).Scan(&recordset).Error

	return recordset, err
}

// *************************************************************************************************************************
// salva registro de fabricante
// *************************************************************************************************************************
func SaveManufacturer(request *types.ManufacturerRequest, c *fiber.Ctx, manufacturerId int) string {

	var err error
	var result string
	var manufacturerId_use int

	// se der erro no upload da image, cancela alteracao/criacao do registro
	trans := Db.Begin()

	manufacturer := types.Manufacturers{
		Name:      request.Name,
		Active:    true,
		Workgroup: c.Params("workgroup"),
	}

	// novo registro
	if c.Method() == "POST" {
		if err = trans.Model(&types.Manufacturers{}).Create(&manufacturer).Error; err != nil {
			trans.Rollback()
			return err.Error()
		}

		manufacturerId_use = manufacturer.Id
		result = "__success__|" + strconv.Itoa(manufacturer.Id)
	}

	// edita registro
	if c.Method() == "PATCH" {
		if err = trans.Model(&types.Manufacturers{}).Where("id = ? ", manufacturerId).Updates(&manufacturer).Error; err != nil {
			trans.Rollback()
			return err.Error()
		}

		manufacturerId_use = manufacturerId
		result = "__success__"
	}

	// quando registro esta sendo editado, ja ha imagem gravada, a verificacao/upload
	// de nova imagem só é necessaria caso o usuario escolheu outra imagem
	// 'BypassImageUpload' é decidida pelo front end
	if !request.BypassImageUpload {
		// se der qq erro com o upload ou integridade da imagem enviada, cancela a insercao/update feito acima
		if err = utils.ValidateIAndUploadImageFile(c, manufacturerId_use); err != nil {
			trans.Rollback()
			return err.Error()
		}
	}

	// necessario efetivar a transacao aqui, caso contrario nao sera possivel registrar notificacao
	trans.Commit()

	// registra que a operacao acima deve ser notificada a todos os usuarios do grupo, exceto o usuario que efetuou a operacao (ip cliente)
	if err = LogNotification(c, manufacturerId_use, request.ClientIp); err != nil {
		return err.Error()
	}

	return result

}

// *************************************************************************************************************************
// retorna campos de determinado registro de fabricante
// *************************************************************************************************************************
func GetManufacturer(manufacturer_id int) (record *types.ManufacturerResponse, err error) {

	// logo_filename ficara parecido com: 	manufacturer_000002.png, manufacturer_000041.png, etc
	_select := " name, concat(workgroup, '_manufacturer_', LPAD(id, 6, '0'), '.png') as manufacturer_logo  "

	// retorna so os campos necessarios
	err = Db.Model(&types.Manufacturers{}).Where("id = ?", manufacturer_id).Select(_select).Scan(&record).Error

	return record, err
}
