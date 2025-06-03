package models

import (
	"fmt"
	"log"
	"os"
	"rentacar/types"
	"rentacar/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB
var DBIP, DBPORT, DBNAME, DBPASS, DBUSER string

var Ctx *fiber.Ctx

// ***************************************************************************************************************************************************************
// struct usada para receber registros de determinada tabela e retorna los como itens de autocomplete - destinados a um: input type=textbox
// ***************************************************************************************************************************************************************
type ItensForAutoComplete struct {
	ItemId   string
	ItemName string
}

// ***************************************************************************************************************************************************************
// struct para retornar metadados de qq tipo de regisro
// ***************************************************************************************************************************************************************

type RecordMetadataResponse struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Workgroup string `json:"workgroup"`
}

//********************************************************************************************************
// este arquivo contem structs que sao usadas por outros arquivos do mesmo package 'models'
// ****************************************************************************************************************************************

type OrderByTranslated struct {
	FieldChosenInDatatable string
	FieldInQuery           string
}

// ***************************************************************************************************************************************************************
// instancia banco de dados mysql
// ***************************************************************************************************************************************************************
func DbSetup() error {

	var err error

	if err = godotenv.Load(); err != nil {
		utils.HighlightError(err)
		return err
	}

	DBIP = os.Getenv("DBIP")
	DBNAME = os.Getenv("DBNAME")
	DBPASS = os.Getenv("DBPASS")
	DBPORT = os.Getenv("DBPORT")
	DBUSER = os.Getenv("DBUSER")

	// instancia logger do gorm, para exibir comandos SQL e detectar eventual erro
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 500 * time.Millisecond, // Slow SQL threshold
			LogLevel:      logger.Info,            // Log level
			Colorful:      true,                   // Disable color
		},
	)

	// conecta à base de dados
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True",
		DBUSER, DBPASS, DBIP, DBPORT, DBNAME)

	var dbTMP *gorm.DB
	dbTMP, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		utils.HighlightError(err)
		return err

	}

	Db = dbTMP

	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")
	fmt.Println("CONECTOU")

	return nil

}

// *************************************************************************************************************************
// muda status (ativo/inativo) de determinado registro
// *************************************************************************************************************************
func ChangeRecordStatus(tableName, recordId string) string {

	result := "__success__"

	// GORM updates() nao suporta 'case when', necessario usar 'Exec'
	try := Db.Exec(fmt.Sprintf("update %v set active = CASE WHEN ifnull(active, false)=false THEN true else false END, updated_at=now() where id=%v", tableName, recordId))

	if try.Error != nil {
		result = try.Error.Error()
	}

	return result
}

// *************************************************************************************************************************
// exclui registro de qq tabela
// *************************************************************************************************************************
func DeleteRecord(tableName, recordId string) string {

	result := "__success__"

	try := Db.Exec(fmt.Sprintf("update %v set deleted_at = now() where id=%v", tableName, recordId))

	if try.Error != nil {
		result = try.Error.Error()
	}
	return result
}

// *************************************************************************************************************************
// exclui registros de qq tabela
// *************************************************************************************************************************
func DeleteRecords(tableName, recordsIds string) string {

	result := "__success__"

	try := Db.Exec(fmt.Sprintf("update %v set deleted_at = now() where id in (%v)", tableName, recordsIds))

	if try.Error != nil {
		result = try.Error.Error()
	}

	return result
}

// *************************************************************************************************************************
// insere na tabela de notificacoes operacao que foi realizada e precisa ser informada a usuarios de determinado grupo
// *************************************************************************************************************************
func LogNotification(c *fiber.Ctx, recordId int, clientIp string) error {

	method := strings.ToLower(c.Method())
	url := string(c.Request().URI().FullURI())

	var expressions types.TermResponse
	var err error

	/****************************************************
	  registro de carro foi inserido ou editado
	 ****************************************************/

	// obtem as frases 'registro de carro inserido/registro de carro alterado', preenche com a imagem do respectivo carro e retorna
	if strings.Contains(url, "/car") {
		if method == "post" {
			expressions, err = GetTermByItem("notification_created_car")
		}
		if method == "patch" {
			expressions, err = GetTermByItem("notification_updated_car")
		}
		if err != nil {
			return err
		}

		// as expressoes notification_updated_car e notification_created_car exigem a imagem <img> do  carro que foi alterado

		// monta nome da imagem do carro baseado em seu ID
		suffixWithID := fmt.Sprintf("%06d", recordId)
		uniqueFileName := fmt.Sprintf("%v_%v_%v.png", c.Params("workgroup"), "car", suffixWithID) // workgroup_car_000005.png, workgroup_car_000053.png, etc
		url := os.Getenv("FTP_URL") + uniqueFileName

		var carRecord *types.CarResponse
		carRecord, err = GetCar(recordId, c)
		if err != nil {
			return err
		}

		expressions.English = strings.ReplaceAll(expressions.English, "@img", url) // substitui a string '@img' pelo link real da imagem
		expressions.Portuguese = strings.ReplaceAll(expressions.Portuguese, "@img", url)

		expressions.English = strings.ReplaceAll(expressions.English, "@name", carRecord.Name) // substitui a string '@name' pelo nome real do carro
		expressions.Portuguese = strings.ReplaceAll(expressions.Portuguese, "@name", carRecord.Name)
	}

	/****************************************************
	  registro de fabricante foi inserido ou editado
	 ****************************************************/

	// obtem as frases 'registro de fabricante inserido/registro de fabricante alterado', preenche com a imagem do respectivo fabricante e retorna
	if strings.Contains(url, "/manufacturer") {
		if method == "post" {
			expressions, err = GetTermByItem("notification_created_manufacturer")
		}
		if method == "patch" {
			expressions, err = GetTermByItem("notification_updated_manufacturer")
		}
		if err != nil {
			return err
		}

		// as expressoes notification_updated_manufacturer e notification_created_manufacturer exigem o logotipo <img> do fabricante que foi alterado

		// monta nome do logotipo do fabricante baseado em seu ID
		suffixWithID := fmt.Sprintf("%06d", recordId)
		uniqueFileName := fmt.Sprintf("%v_%v_%v.png", c.Params("workgroup"), "manufacturer", suffixWithID) // workgroup_manufacturer_000005.png, workgroup_manufacturer_000053.png, etc
		url := os.Getenv("FTP_URL") + uniqueFileName

		var manufacturerRecord *types.ManufacturerResponse
		manufacturerRecord, err = GetManufacturer(recordId)
		if err != nil {
			return err
		}
		expressions.English = strings.ReplaceAll(expressions.English, "@img", url) // substitui a string '@img' pelo link real do logotipo
		expressions.Portuguese = strings.ReplaceAll(expressions.Portuguese, "@img", url)

		expressions.English = strings.ReplaceAll(expressions.English, "@name", manufacturerRecord.Name) // substitui a string '@name' pelo nome real do fabricante
		expressions.Portuguese = strings.ReplaceAll(expressions.Portuguese, "@name", manufacturerRecord.Name)
	}

	/****************************************************************************************
	  registro de expressao ingles/portugues foi inserido ou editado
	****************************************************************************************/

	// obtem as frases 'registro de expressão inserida/registro de expressão alterada', e preenche com uma imagem ilustrativa
	if strings.Contains(url, "/term") {
		if method == "post" {
			expressions, err = GetTermByItem("notification_created_expression")
		}
		if method == "patch" {
			expressions, err = GetTermByItem("notification_updated_expression")
		}
		if err != nil {
			return err
		}

		// as expressoes notification_updated_expression e notification_created_expression exigem uma imagem ilustrativa <img> que esta no repositorio FTP
		url := os.Getenv("FTP_URL_STATIC") + "language_notification.png"

		var termRecord *types.TermResponse
		termRecord, err = GetTerm(recordId)
		if err != nil {
			return err
		}
		expressions.English = strings.ReplaceAll(expressions.English, "@img", url) // substitui a string '@img' pela imagem ilustrativa
		expressions.Portuguese = strings.ReplaceAll(expressions.Portuguese, "@img", url)

		expressions.English = strings.ReplaceAll(expressions.English, "@item", termRecord.Item) // substitui a string '@item' pelo item indicador da expressao
		expressions.Portuguese = strings.ReplaceAll(expressions.Portuguese, "@item", termRecord.Item)
	}

	/****************************************************************************************
	  registro de reserva de veiculo foi inserido ou editado
	****************************************************************************************/

	// obtem as frases 'registro de reserva inserida/registro de reserva alterada', e preenche com uma imagem ilustrativa
	if strings.Contains(url, "/booking") {
		if method == "post" {
			expressions, err = GetTermByItem("notification_created_booking")
		}
		if method == "patch" {
			expressions, err = GetTermByItem("notification_updated_booking")
		}
		if err != nil {
			return err
		}

		// as expressoes notification_updated_booking e notification_created_booking exigem a imagem <img> do carro que foi reservado

		var bookingRecord *types.BookingResponse
		bookingRecord, err = GetBooking(recordId)
		if err != nil {
			return err
		}

		// monta nome dda imagem do carrro baseado em seu ID
		suffixWithID := fmt.Sprintf("%06d", bookingRecord.CarId)
		uniqueFileName := fmt.Sprintf("%v_%v_%v.png", c.Params("workgroup"), "car", suffixWithID) // workgroup_car_000005.png, workgroup_car_000053.png, etc
		url := os.Getenv("FTP_URL") + uniqueFileName

		expressions.English = strings.ReplaceAll(expressions.English, "@img", url) // substitui a string '@img' pela imagem do carro
		expressions.Portuguese = strings.ReplaceAll(expressions.Portuguese, "@img", url)

		expressions.English = strings.ReplaceAll(expressions.English, "@date", bookingRecord.PickupDate) // substitui a string '@date' pela data de retirada do carro
		expressions.English = strings.ReplaceAll(expressions.English, "@hour", bookingRecord.PickupHour) // @hour = hora da retirada do carro

		expressions.Portuguese = strings.ReplaceAll(expressions.Portuguese, "@date", bookingRecord.PickupDate) // @hour = hora da retirada do carro
		expressions.Portuguese = strings.ReplaceAll(expressions.Portuguese, "@hour", bookingRecord.PickupHour) // @hour = hora da retirada do carro
	}

	//****************************************************************************************
	// grava a notificacao preparada acima
	//****************************************************************************************
	notification := types.Notifications{
		Workgroup:             c.Params("workgroup"),
		MadeByIP:              clientIp,
		DescriptionEnglish:    expressions.English,
		DescriptionPortuguese: expressions.Portuguese,
	}

	if err = Db.Model(&types.Notifications{}).Create(&notification).Error; err != nil {
		return err
	}

	// contabiliza 1 operacao (inclusao, edicao) a mais feita pelo grupo
	WorkgroupAddLog(c.Params("workgroup"))

	return nil
}

//*************************************************************************************************************
// verifica se ha notificacoes a serem enviadas ao front end, dentro do grupo atual
// notificacao = alteracao registro de carro, nova reserva, etc
//*************************************************************************************************************

func CheckNotifications(c *fiber.Ctx) (newNotifications *[]types.NotificationsResponse, err error) {

	fields := " id, description_english, description_portuguese "
	where := fmt.Sprintf(" trim(workgroup) = trim('%v') and made_by_ip<>'%v' ", c.Params("workgroup"), c.Params("client_ip"))

	if err = Db.Model(&types.Notifications{}).Order("id desc").Select(fields).Where(where).Scan(&newNotifications).Error; err != nil {
		return nil, err
	}

	return newNotifications, nil

}
