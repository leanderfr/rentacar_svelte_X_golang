package models

import (
	"fmt"
	"strconv"

	"rentacar/types"
	"rentacar/utils"

	"github.com/gofiber/fiber/v2"
)

// *************************************************************************************************************************
// retorna dados minimos das frases em ingles/portugues, baseado no idioma selecionado pelo usuario no front end
// *************************************************************************************************************************
func GetTermsForFillingFrontEnd(language string, workgroup string) (expressions *[]types.TermsForPopulateFronEndResponse, err error) {

	// busca o campo dependendo de qual idioma requisitado
	var _select string
	if language == "portuguese" {
		_select = " portuguese as expression, item  "
	} else if language == "english" {
		_select = " english as expression, item  "
	}

	where := fmt.Sprintf(" workgroup= '%v' and ifnull(active, false)= true  ", workgroup)

	// retorna so os campos necessarios, expressoes ativas
	if err := Db.Model(&types.Terms{}).Select(_select).Where(where).Scan(&expressions).Error; err != nil {
		return nil, err
	}

	return expressions, err
}

// *************************************************************************************************************************
// retorna lista de todas as frases em ingles/portugues usadas no front end
// *************************************************************************************************************************
func GetTermsForDatatable(request *types.DatatableParamsRequest, c *fiber.Ctx) (recordset *[]types.TermsResponse, err error) {

	// nao foi passado campo order by inicialmente
	if request.OrderBy == "" {
		request.OrderBy = "name"
		request.OrderDirection = "asc"
	}

	// busca o campo dependendo de qual idioma requisitado
	_select := " id, portuguese, english, item, description, ifnull(active, 0) as active   "

	where := fmt.Sprintf(" workgroup= '%v'  ", c.Params("workgroup"))

	if request.SearchtTxt != "" {
		// usa todas as colunas visiveis no datatable do front end
		request.SearchtTxt = "%" + request.SearchtTxt + "%"
		where += utils.ConcatWhere(where, fmt.Sprintf(" (item like '%v' or portuguese like '%v' or english like '%v' or description like '%v') ",
			request.SearchtTxt, request.SearchtTxt, request.SearchtTxt, request.SearchtTxt))
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

	// retorna so os campos necessarios
	if where != "" {
		err = Db.Model(&types.Terms{}).Order(request.OrderBy).Where(where).Select(_select).Scan(&recordset).Error
	} else {
		err = Db.Model(&types.Terms{}).Order(request.OrderBy).Select(_select).Scan(&recordset).Error
	}
	return recordset, err
}

// *************************************************************************************************************************
// salva registro de frase em ingles/portugues
// *************************************************************************************************************************
func SaveTerm(request *types.TermRequest, c *fiber.Ctx, termId int) string {

	var err error
	var result string

	term := types.Terms{
		Item:        request.Item,
		Portuguese:  request.Portuguese,
		English:     request.English,
		Description: request.Description,
		Active:      true,
		Workgroup:   c.Params("workgroup"),
	}

	var termId_use int

	// novo registro
	if c.Method() == "POST" {
		if err = Db.Model(&types.Terms{}).Create(&term).Error; err != nil {
			return err.Error()
		}
		termId_use = term.Id
		result = "__success__|" + strconv.Itoa(term.Id)
	}

	// edita registro
	if c.Method() == "PATCH" {
		if err = Db.Model(&types.Terms{}).Where("id = ? ", termId).Updates(&term).Error; err != nil {
			return err.Error()
		}
		termId_use = termId
		result = "__success__"
	}

	// registra que a operacao acima deve ser notificada a todos os usuarios do grupo, exceto o usuario que efetuou a operacao (ip cliente)
	if err = LogNotification(c, termId_use, request.ClientIp); err != nil {
		return err.Error()
	}

	return result

}

// *************************************************************************************************************************
// retorna campos de determinado registro de expressao ingles/portugues
// *************************************************************************************************************************
func GetTerm(recordId int) (record *types.TermResponse, err error) {

	// busca o campo dependendo de qual idioma requisitado
	_select := " portuguese, english, item, description  "

	// retorna so os campos necessarios
	err = Db.Model(&types.Terms{}).Where("id = ?", recordId).Select(_select).Scan(&record).Error

	return record, err
}

// *************************************************************************************************************************
// retorna texto em ingles/portugues de determinado item, funcao auxiliar na gravacao de notificacao (LogNotification)
// exemplo, notificacao de gravacao de registro de carro: item=  notification_created_car ou notification_updated_car
// *************************************************************************************************************************
func GetTermByItem(item string) (record types.TermResponse, err error) {

	_select := " portuguese, english, item, description  "

	// retorna so os campos necessarios
	err = Db.Model(&types.Terms{}).Where("item = ?", item).Select(_select).Limit(1).Scan(&record).Error

	return record, err
}
