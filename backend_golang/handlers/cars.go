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
// retorna lista de carros disponiveis para aluguel que serao exibidos na pagina /home
// *************************************************************************************************************************
func GetCarsForDatatable(c *fiber.Ctx) error {
	var err error
	var request types.DatatableParamsRequest

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if strings.TrimSpace(c.Params("workgroup")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}

	// verifica se campos minimos foram enviados (country, order_by, order_direction, etc)
	// FullURI exemplo ==> http://localhost:8070/admin/cars?country=brazil&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=
	_url := string(c.Request().URI().FullURI())

	if err = utils.UrlToStruct(_url, &request); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// monta lista de carros baseado no país selecionado atualmente no front end
	var recordset *[]types.CarsResponse

	if recordset, err = models.GetCarsForDatatable(&request, c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if recordset == nil {
		emptyJSON := make([]types.CarsResponse, 0)
		c.JSON(emptyJSON)
	} else {
		c.JSON(*recordset)
	}

	return nil
}

// *************************************************************************************************************************
// retorna cards de carros disponiveis
// *************************************************************************************************************************
func GetCarsForCards(c *fiber.Ctx) error {

	country := c.Params("country")
	workgroup := c.Params("workgroup")

	// lista de carros baseado no país/idioma selecionado no front end
	var recordset []*types.CarsResponse
	var err error
	if recordset, err = models.GetCarsForCards(country, workgroup); err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if recordset == nil {
		emptyJSON := make([]types.CarsResponse, 0)
		c.JSON(emptyJSON)
	} else {
		c.JSON(recordset)
	}

	return nil
}

// *************************************************************************************************************************
// lê registro de determinado carro
// *************************************************************************************************************************
func GetCar(c *fiber.Ctx) error {

	var carId int
	var err error
	var data *types.CarResponse

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if carId, err = strconv.Atoi(c.Params("id")); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if data, err = models.GetCar(carId, c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// se nenhum registro retornado
	if data == nil {
		emptyJSON := make([]types.CarResponse, 0)
		c.JSON(emptyJSON)

	} else {
		c.JSON(*data)
	}
	return nil
}

// *************************************************************************************************************************
// salva registro de carro
// *************************************************************************************************************************
func SaveCar(c *fiber.Ctx) error {

	// verifica se todos campos/seus tipos vieram ok
	var request types.CarRequest
	var err error

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if strings.TrimSpace(c.Params("workgroup")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}

	if err = c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("a" + err.Error())
	}

	// verifica se campos (text/textarea) obrigatorios vieram, tamanhos e quantidades min/maxima, caso alguem tente injetar campos "gigantes", estando logado
	validate := validator.New()

	if err = validate.Struct(&request); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusInternalServerError).SendString("b" + errors.Error())
	}

	// c.Method()  ==>  POST (novo reg) ou PATCH (editar reg)
	carId := 0 // POST
	if c.Method() == "PATCH" {
		if carId, err = strconv.Atoi(c.Params("id")); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	try := models.SaveCar(&request, c, carId)

	if !strings.Contains(try, "__success__") {
		return c.Status(fiber.StatusInternalServerError).SendString(try) // erro
	} else {
		return c.Status(fiber.StatusOK).SendString(try) // gravacao ok
	}

}

// *************************************************************************************************************************
// applicacao dando 'oi' ao ser invocada sem rota
// *************************************************************************************************************************

func HelloWorld(c *fiber.Ctx) error {

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return c.Status(fiber.StatusOK).SendString(" OK *** OK *** Ok  **  OK ***")

}
