package handlers

import (
	"reflect"
	"rentacar/models"
	"rentacar/types"
	"rentacar/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// *************************************************************************************************************************
// lista de grupos
// *************************************************************************************************************************
func GetWorkgroupsForDatatable(c *fiber.Ctx) error {
	var err error
	var request types.DatatableParamsRequest

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

	// monta lista de grupos
	var recordset *[]types.WorkgroupsResponse

	if recordset, err = models.GetWorkgroupsForDatatable(&request, c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if recordset == nil {
		emptyJSON := make([]types.WorkgroupsResponse, 0)
		c.JSON(emptyJSON)
	} else {
		c.JSON(*recordset)
	}

	return nil
}

// *************************************************************************************************************************
// 1. escolhe aleatoriamente um grupo dentro dos grupos disponiveis na tabelas 'workgroups' (apliacao executada 1a vez)
// ou
// 2. reseta os dado do grupo atual (usuario pediu para resetar dados do grupo)
// ou
// 3. desativa o grupo atual e escolhe outro grupo aleatoriamente (usuario pediu para sortear outro grupo)
// *************************************************************************************************************************
func GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother(c *fiber.Ctx) error {

	var try string
	var err error
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// qdo aplicacao front sendo carregada 1a vez, workgroup= none, action= generate
	// qdo usuario pediu para resetr dados,  workgroup= <grupo atual>, action= reset
	// qdo usuario pediu para sortear outro grupo dados,  workgroup= <grupo atual>, action= another
	workgroup := strings.ToLower(strings.TrimSpace(c.Params("workgroup")))
	whatToDo := strings.ToLower(strings.TrimSpace(c.Params("action")))

	// impede grupo admin de ser resetado ou substituido por outro (sorteio)
	if (whatToDo == "reset" || whatToDo == "another") && (workgroup == "" || workgroup == "admin") {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}

	if whatToDo != "reset" && whatToDo != "generate" && whatToDo != "another" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing whatToDo info")
	}

	if try, err = models.GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if !strings.Contains(try, "__success__") {
		return c.Status(fiber.StatusInternalServerError).SendString(try)

	} else {
		return c.Status(fiber.StatusOK).SendString(try) // __success__ | _workgroup_.Name | _workgroup_.Id
	}
}

// *************************************************************************************************************************
// mantem front end informado de como esta a clonagem dos registros do workgroup recem criado
// *************************************************************************************************************************

func ShowCloningWorkgroupStatus(c *fiber.Ctx) error {

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// ao iniciar a aplicacao web, ela fica pedindo informacoes a cada 1/2 segundo de como esta o processo de clonagem
	// isso ocorre no inicio da app, uma unica vez, e talvez a clonagem nao iniciou ou nem sera necessario,
	// para nao ocorrer erro, necessario retornar string em branco
	if reflect.ValueOf(models.ProcessStatus).IsZero() {
		return c.Status(fiber.StatusOK).SendString("")

	} else {
		return c.Status(fiber.StatusOK).SendString(models.ProcessStatus.Status + "|" + models.ProcessStatus.PercentReady + "|" + models.ProcessStatus.ChosenWorkgroup)
	}

}

// *************************************************************************************************************************
// acessar dados de outro grupo
// *************************************************************************************************************************
func AccessAnotherGroupData(c *fiber.Ctx) error {

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	_workgroup := strings.ToLower(strings.TrimSpace(c.Params("workgroup")))
	if _workgroup == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}
	// nao permite acessar admin
	if _workgroup == "admin" {
		return c.Status(fiber.StatusOK).SendString("none")
	}

	var try string
	var err error

	if try, err = models.AccessAnotherGroupData(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// success= achou grupo,  none= grupo nao existe
	if !strings.Contains(try, "__success__") && !strings.Contains(try, "none") {
		return c.Status(fiber.StatusInternalServerError).SendString(try)

	} else {
		return c.Status(fiber.StatusOK).SendString(try) // __success__ | _workgroup_.Name | _workgroup_.Id
	}

}

// *************************************************************************************************************************
// obtem a qtde de alteracoes ja feitas na base pelo grupo indicado
// *************************************************************************************************************************
func WorkgroupReport(c *fiber.Ctx) error {

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if strings.TrimSpace(c.Params("workgroup")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}

	var try string
	var err error

	if try, err = models.WorkgroupReport(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if !strings.Contains(try, "__success__") {
		return c.Status(fiber.StatusInternalServerError).SendString(try)

	} else {
		return c.Status(fiber.StatusOK).SendString(try) // __success__ |qtde alteracoes feitas na base
	}

}
