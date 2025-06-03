package handlers

import (
	"rentacar/models"
	"rentacar/types"
	"rentacar/utils"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// *************************************************************************************************************************
// retorna lista de facricante
// *************************************************************************************************************************
func GetManufacturersForDatatable(c *fiber.Ctx) error {
	var err error
	var request types.DatatableParamsRequest

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if strings.TrimSpace(c.Params("workgroup")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}
	// verifica se campos minimos foram enviados (country, order_by, order_direction, etc)
	// FullURI exemplo ==> http://localhost:8070/manufacturers?country=brazil&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=
	_url := string(c.Request().URI().FullURI())

	if err = utils.UrlToStruct(_url, &request); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// monta lista com fabricantes
	car_records, _ := models.GetManufacturersForDatatable(&request, c)

	if car_records == nil {
		emptyJSON := make([]types.ManufacturerResponse, 0)
		c.JSON(emptyJSON)
	} else {
		c.JSON(*car_records)
	}

	return nil
}

// *************************************************************************************************************************
// salva registro de fabricante
// *************************************************************************************************************************
func SaveManufacturer(c *fiber.Ctx) error {

	// verifica se todos campos/seus tipos vieram ok
	var request types.ManufacturerRequest
	var err error

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if strings.TrimSpace(c.Params("workgroup")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}

	if err = c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// verifica se campos obrigatorios vieram, tamanhos e quantidades min/maxima, caso alguem tente injetar campos "gigantes", estando logado
	validate := validator.New()

	if err = validate.Struct(&request); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusInternalServerError).SendString(errors.Error())
	}

	// c.Method()  ==>  POST (novo reg) ou PATCH (editar reg)
	manufacturerId := 0 // POST
	if c.Method() == "PATCH" {
		if manufacturerId, err = strconv.Atoi(c.Params("id")); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	try := models.SaveManufacturer(&request, c, manufacturerId)

	if !strings.Contains(try, "__success__") {
		return c.Status(fiber.StatusInternalServerError).SendString(try) // erro
	} else {
		return c.Status(fiber.StatusOK).SendString(try) // gravacao ok
	}

}

// *************************************************************************************************************************
// lê registro de determinado fabricante de carros
// *************************************************************************************************************************
func GetManufacturer(c *fiber.Ctx) error {

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	var err error
	var manufacturer_id int

	manufacturer_id, err = strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	data, err := models.GetManufacturer(manufacturer_id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// se nenhum registro retornado
	if data == nil {
		emptyJSON := make([]types.ManufacturerResponse, 0)
		c.JSON(emptyJSON)

	} else {
		c.JSON(*data)
	}
	return nil
}
