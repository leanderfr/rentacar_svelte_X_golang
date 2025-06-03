package handlers

import (
	"fmt"
	"rentacar/models"
	"rentacar/types"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// *************************************************************************************************************************
// carrega itens (registros) para preenchimento de um autocomplete no front end
// *************************************************************************************************************************
func ItensForAutocomplete(c *fiber.Ctx) error {

	table_name := c.Params("table_name")
	workgroup := c.Params("workgroup")

	var itens []*models.ItensForAutoComplete
	var err error

	sql := fmt.Sprintf("SELECT id as item_id, name as item_name FROM %v  where ifnull(active, false) = true and deleted_at is null and workgroup = '%v' order by name", table_name, workgroup)

	if err = models.Db.Raw(sql).Scan(&itens).Error; err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// monta string com registros e colunas para autocomplete
	// ao inves de retornar como JSON, vai retornar uma string, para facilitar o processamento pelo jscript
	strItens := ""
	for _, item := range itens {

		if strItens != "" {
			strItens += "|" // separador de registros
		}

		strItens += item.ItemId + ";" + item.ItemName
	}

	fmt.Fprint(c.Response().BodyWriter(), strItens)

	return nil
}

// *************************************************************************************************************************
// lê metadados de determinado registro de qq tabela
// *************************************************************************************************************************
func GetRecordMetadata(c *fiber.Ctx) (err error) {

	// necessario fornecer o país selecionado no front end por causa do formato da data retornado (diferente Brasil/USA)
	recordId := c.Params("id")
	country := c.Params("country")
	tableName := c.Params("table_name")

	//var dateFormat, locale string
	var dateFormat string
	if country == "usa" {
		dateFormat = "%M / %d / %Y - %h:%i %p" // mm/dd/yyyy hh:mm am/pm
		//locale = "en_US"

	} else if country == "brazil" {
		dateFormat = "%d / %M / %Y - %H:%i" // dd/mm/yyyy HH:mm
		//locale = "pt_BR"
	}

	var record *models.RecordMetadataResponse

	var join, _fields string

	//	_fields = fmt.Sprintf(" date_format(%v.created_at, '%v', '%v') as created_at, date_format(%v.updated_at, '%v', '%v') as updated_at, "+
	//		" workgroup  ", tableName, dateFormat, locale, tableName, dateFormat, locale)

	_fields = fmt.Sprintf(" date_format(%v.created_at, '%v') as created_at, date_format(%v.updated_at, '%v') as updated_at, "+
		" workgroup  ", tableName, dateFormat, tableName, dateFormat)

	if tableName == "cars" {
		err = models.Db.Model(&types.Cars{}).Joins(join).Where("cars.id = ?", recordId).Select(_fields).Scan(&record).Error
	} else if tableName == "terms" {
		err = models.Db.Model(&types.Terms{}).Joins(join).Where("terms.id = ?", recordId).Select(_fields).Scan(&record).Error
	} else if tableName == "manufacturers" {
		err = models.Db.Model(&types.Manufacturers{}).Joins(join).Where("manufacturers.id = ?", recordId).Select(_fields).Scan(&record).Error
	} else if tableName == "bookings" {
		err = models.Db.Model(&types.Bookings{}).Joins(join).Where("bookings.id = ?", recordId).Select(_fields).Scan(&record).Error

	}

	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// se nenhum registro retornado
	if record == nil {
		emptyJSON := make([]models.RecordMetadataResponse, 0)
		c.JSON(emptyJSON)

		// se algum registro retornado
	} else {
		c.JSON(&record)
	}
	return nil
}

// *************************************************************************************************************************
// muda status (ativo/inativo) de determinado registro
// *************************************************************************************************************************
func ChangeRecordStatus(c *fiber.Ctx) error {

	tableName := c.Params("tablename")
	recordId := c.Params("id")

	try := models.ChangeRecordStatus(tableName, recordId)
	fmt.Fprint(c.Response().BodyWriter(), try)

	return nil
}

// *************************************************************************************************************************
// exclui registro de qq tabela
// *************************************************************************************************************************
func DeleteRecord(c *fiber.Ctx) error {

	tableName := c.Params("tablename")
	recordId := c.Params("id")

	try := models.DeleteRecord(tableName, recordId)
	fmt.Fprint(c.Response().BodyWriter(), try)

	return nil
}

// *************************************************************************************************************************
// exclui registros de qq tabela que foram selecionados no front end (datatable)
// *************************************************************************************************************************
func DeleteRecords(c *fiber.Ctx) error {

	// 'ids' contem IDs separados por virgula
	tableName := c.Params("tablename")
	recordsIds := c.Params("ids")

	try := models.DeleteRecords(tableName, recordsIds)
	fmt.Fprint(c.Response().BodyWriter(), try)

	return nil
}

//*************************************************************************************************************
// verifica se ha notificacoes a serem enviadas ao front end, dentro do grupo atual
//*************************************************************************************************************

func CheckNotifications(c *fiber.Ctx) error {

	var err error

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if strings.TrimSpace(c.Params("workgroup")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}
	if strings.TrimSpace(c.Params("client_ip")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing client ip")
	}

	var newNotifications *[]types.NotificationsResponse

	if newNotifications, err = models.CheckNotifications(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if newNotifications == nil {
		emptyJSON := make([]types.NotificationsResponse, 0)
		c.JSON(emptyJSON)
	} else {
		c.JSON(*newNotifications)
	}

	return nil

}
