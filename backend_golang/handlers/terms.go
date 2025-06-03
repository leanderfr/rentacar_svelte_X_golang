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
// retorna lista de expressoes ingles/portugues usadas no front end
// *************************************************************************************************************************
func GetTermsForDatatable(c *fiber.Ctx) error {
	var err error
	var request types.DatatableParamsRequest
	var data *[]types.TermsResponse

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if strings.TrimSpace(c.Params("workgroup")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}

	// verifica se campos minimos foram enviados (country, order_by, order_direction, etc)
	// FullURI exemplo ==> http://localhost:8070/cars?country=brazil&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=
	_url := string(c.Request().URI().FullURI())

	if err = utils.UrlToStruct(_url, &request); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// as palavras sao exibidas dependendo do país selecionado pelo usuario (country)

	if data, err = models.GetTermsForDatatable(&request, c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if data == nil {
		emptyJSON := make([]types.TermsResponse, 0)
		c.JSON(emptyJSON)
	} else {
		c.JSON(*data)
	}
	return nil
}

// *************************************************************************************************************************
// retorna expressoes do idioma escolhido pelo usuario no front end
// *************************************************************************************************************************
func GetTermsForFillingUpFrontEnd(c *fiber.Ctx) error {
	language := c.Params("language")
	workgroup := c.Params("workgroup")

	// qdo é feita a 1a carga no front end, o arquivo +page.js nao tem acesso ao workgroup ainda
	// nesse caso, pegar expressoes ing/port do grupo principal (Admin)
	if workgroup == "" {
		workgroup = "Admin"
	}

	// as palavras sao obtidas  dependendo do país selecionado pelo usuario (language_id)
	recordset, err := models.GetTermsForFillingFrontEnd(language, workgroup)

	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if recordset == nil {
		emptyJSON := make([]types.TermsResponse, 0)
		c.JSON(emptyJSON)
	} else {
		c.JSON(*recordset)
	}

	return nil
}

// *************************************************************************************************************************
// salva registro de frase em ingles/portugues
// *************************************************************************************************************************
func SaveTerm(c *fiber.Ctx) error {

	// verifica se todos campos/seus tipos vieram ok
	var request types.TermRequest
	var err error

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if strings.TrimSpace(c.Params("workgroup")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}

	if err = c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// verifica campos obrigatorios, tamanhos e quantidades min/maxima, caso a verificacao front end tenha falhado
	validate := validator.New()

	// Validate the User struct
	if err = validate.Struct(&request); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusInternalServerError).SendString(errors.Error())
	}

	// c.Method()  ==>  POST (novo reg) ou PATCH (editar reg)
	termId := 0 //
	if c.Method() == "PATCH" {
		if termId, err = strconv.Atoi(c.Params("id")); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	try := models.SaveTerm(&request, c, termId)

	if !strings.Contains(try, "__success__") {
		return c.Status(fiber.StatusInternalServerError).SendString(try) // erro
	} else {
		return c.Status(fiber.StatusOK).SendString(try) // gravacao ok
	}

}

// *************************************************************************************************************************
// lê  registro de frase em ingles/portugues solicitada pelo front end
// *************************************************************************************************************************
func GetTerm(c *fiber.Ctx) error {

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	var err error
	var expressionId int

	expressionId, err = strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	data, err := models.GetTerm(expressionId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// se nenhum registro retornado
	if data == nil {
		emptyJSON := make([]types.TermResponse, 0)
		c.JSON(emptyJSON)

	} else {
		c.JSON(*data)
	}
	return nil
}
