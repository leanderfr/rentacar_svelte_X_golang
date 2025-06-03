package main

import (
	"fmt"
	"os"
	"rentacar/handlers"
	"rentacar/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var App *fiber.App

var Ctx *fiber.Ctx

func main() {

	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")
	fmt.Println("CHEGOU")
	// conecta com base com base
	if err := models.DbSetup(); err != nil {
		return
	}

	App = fiber.New()

	App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Accept-Language, Content-Length",
	}))

	//*************************************************************************************************************************************
	// rotas que exigem como 1o parametro a informacao 'workgroup', serao vinculadas a 'groupRequest'
	groupRequest := App.Group("/:workgroup")
	//*************************************************************************************************************************************

	App.Get("/", handlers.HelloWorld)
	// carrega as expressoes em portugues/ingles, mas somente do idioma escolhido pelo usuario no front end
	// qdo chamado sem especificar workgroup, retorna os termos padrao da aplicacao (workgroup= Admin)
	// a requisicao abaixo é feita quando usuario carrega aplicacao pela 1a vez e nao ha workgroup definido ainda, para poder iniciar a aplicacao
	App.Get("/terms/:language", handlers.GetTermsForFillingUpFrontEnd)

	// o front end precisa ter um 'workgroup', informado em toda requisicao de leitura/gravacao de dados
	// os dados sao separados por workgroup, cada grupo de trabalho enxerga e altera somente seus dados

	// usuario pediu para acessar dados de outro grupo
	App.Get("/workgroup/change/:workgroup", handlers.AccessAnotherGroupData)

	// a rota abaixo serve para 3 situacoes:
	// 1- front end sendo executado pela 1a vez, foi feito pedido para sortear algum workgroup e preparar seus dados  (action= generate)
	// 2- front end ja tem workgroup/dados definidos e usuario pediu para resetar os dados do workgroup  (action= reset)
	// 3- front end ja tem workgroup/dados definidos e usuario pediu para sortear outro grupo  (action= another)
	App.Post("/workgroup/:action/:workgroup", handlers.GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother)

	// os dados sao separados por grupo de trabalho, cada grupo tem seus carros, fabricantes, etc
	// quando um grupo de trabalho é iniciado, ele recebe uma cópia dos registros padrao (carros, fabricantes, etc)
	App.Get("/cloning_status", handlers.ShowCloningWorkgroupStatus)

	// le metadados de registro de qq tabela
	// necessaria obter o país para retornar a data no respectivo formato
	App.Get("/:table_name/record_metadata/:id/:country", handlers.GetRecordMetadata)

	// altera status (ativo/inativo) de determinado registro de qq tabela
	// exemplo:  /cars/status/3  , /manufacturers/status/2  , ...
	App.Patch("/:tablename/status/:id", handlers.ChangeRecordStatus)

	// exclui registro de qq tabela
	App.Delete("/:tablename/delete/:id", handlers.DeleteRecord)

	// exclui registros (varios selecionados na datatable (front-end)) de qq tabela
	// ids separados por virgula
	App.Delete("/:tablename/batch_delete/:ids", handlers.DeleteRecords)

	// carrega as expressoes em portugues/ingles, mas somente do idioma escolhido pelo usuario no front end e do grupo 'Admin'
	App.Get("/terms/:language", handlers.GetTermsForFillingUpFrontEnd)

	// obtem a qtde de alteracoes ja feitas na base
	groupRequest.Get("/report", handlers.WorkgroupReport)

	// carrega as expressoes em portugues/ingles, mas somente do idioma escolhido pelo usuario no front end, no momento
	// e pertencentes ao grupo logado (workgroup)
	groupRequest.Get("/terms/:language", handlers.GetTermsForFillingUpFrontEnd)

	// cards com informacoes basicas de cada veiculo, baseado no país que esta selecionado pelo usuario no front end
	groupRequest.Get("/car_cards/:country", handlers.GetCarsForCards)

	// carrega itens (registros) para preenchimento de um autocomplete no front end
	groupRequest.Get("/:table_name/itens_for_autocomplete/", handlers.ItensForAutocomplete)

	// verifica se grupo atual (front end) possui notificacões a serem exibidas na barra de notificacoes
	// notificacao= alteracao de registro de carro, fabricante, nova reserva, etc que deva ser exibida no front
	groupRequest.Get("/notifications/:client_ip", handlers.CheckNotifications)

	/********************************************************************************************************
	  padrao usado:
	    GET= ler registro
	    POST= novo registro
	    PATCH= alterar registro
	    DELETE= excluir registro
	  ********************************************************************************************************/

	//********************************************************************************************************
	// datatable e crud de veiculos
	//********************************************************************************************************
	groupRequest.Get("/cars", handlers.GetCarsForDatatable)
	groupRequest.Post("/car", handlers.SaveCar)
	groupRequest.Patch("/car/:id", handlers.SaveCar)

	App.Get("/car/:id", handlers.GetCar)                      // busca somente detalhes do carro
	App.Get("/car/:id/:first_day/:last_day", handlers.GetCar) // busca detalhes do carro e os dias do mes/ano (dentro de um periodo) em que ele ja esta reservado

	//********************************************************************************************************
	// terms= expressoes em cada idioma (portugues/ingles)
	//********************************************************************************************************
	groupRequest.Get("/terms", handlers.GetTermsForDatatable)
	groupRequest.Post("/term", handlers.SaveTerm)
	groupRequest.Patch("/term/:id", handlers.SaveTerm)

	App.Get("/term/:id", handlers.GetTerm)

	//********************************************************************************************************
	groupRequest.Get("/manufacturers", handlers.GetManufacturersForDatatable)
	groupRequest.Post("/manufacturer", handlers.SaveManufacturer)
	groupRequest.Patch("/manufacturer/:id", handlers.SaveManufacturer)

	App.Get("/manufacturer/:id", handlers.GetManufacturer)

	//********************************************************************************************************
	// reservas de veiculos
	//********************************************************************************************************
	groupRequest.Post("/booking", handlers.SaveBooking)
	groupRequest.Patch("/booking/:id", handlers.SaveBooking)

	groupRequest.Get("/bookings/:country/:car_id/:first_day_week/:last_day_week", handlers.GetBookingsToPopulateSchedule)

	App.Get("/booking/:id", handlers.GetBooking)

	//********************************************************************************************************
	// workgroups sao pre definidos na base, nao sao editaveis, somente é possivel lista los
	//********************************************************************************************************
	groupRequest.Get("/workgroups", handlers.GetWorkgroupsForDatatable)

	App.Listen(":" + os.Getenv("PORT"))

}
