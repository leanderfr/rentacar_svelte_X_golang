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

// *************************************************************************************************************************
// retorna lista dos carros baseado no país escolhido no front end
// *************************************************************************************************************************
func GetCarsForDatatable(request *types.DatatableParamsRequest, c *fiber.Ctx) (recordset *[]types.CarsResponse, err error) {

	// nao foi passado campo order by inicialmente
	if request.OrderBy == "" {
		request.OrderBy = "name"
		request.OrderDirection = "asc"
	}

	// arq imagem do carro:   workgroup_car_999999.png

	fields :=
		` cars.id , cars.name,  concat(cars.workgroup, '_car_', LPAD(cars.id, 6, '0'), '.png') as car_image,  ` +
			` ifnull(cars.active, false) as active, concat(cars.workgroup, '_manufacturer_', LPAD(manufacturers.id, 6, '0'), '.png') as manufacturer_logo  `

	join := " left join manufacturers on manufacturers.id = cars.manufacturer_id "

	where := fmt.Sprintf(" country = '%v' and cars.workgroup= '%v' ", request.Country, c.Params("workgroup"))

	if request.SearchtTxt != "" {
		// usa todas as colunas visiveis no datatable do front end
		request.SearchtTxt = "%" + request.SearchtTxt + "%"
		where += utils.ConcatWhere(where, fmt.Sprintf("  (cars.name like '%v' or manufacturers.name like '%v') ", request.SearchtTxt, request.SearchtTxt))
	}
	if request.OnlyActiveOrInactiveRecords == "active" {
		// so ativos
		where += utils.ConcatWhere(where, " ifnull(cars.active, false) = true ")
	}
	if request.OnlyActiveOrInactiveRecords == "inactive" {
		// so inativos
		where += utils.ConcatWhere(where, " ifnull(cars.active, false) = false ")
	}

	request.OrderBy = CarsOrderByTranslated(request.OrderBy) + " " + request.OrderDirection

	// busca somente os campos necessarios
	if err = Db.Model(&types.Cars{}).Order(request.OrderBy).Select(fields).Joins(join).Where(where).Scan(&recordset).Error; err != nil {
		return nil, err
	}

	return recordset, nil
}

// *************************************************************************************************************************
// retorna informacoes minimas sobre o carro para montagem do seu card
// *************************************************************************************************************************
func GetCarsForCards(country string, workgroup string) (car_cards []*types.CarsResponse, err error) {

	// lê demais campos que serao retornados ao frontend
	fields := " cars.id as Id, cars.name as Name, rental_price as RentalPrice, concat(cars.workgroup, '_car_', LPAD(cars.id, 6, '0'), '.png') as car_image  "
	where := fmt.Sprintf("country = '%v' and ifnull(active, false) = true and workgroup='%v' ", country, workgroup)

	// busca somente os campos necessarios, organiza randomicamente, para sortear carros a serem vistos por primeiro
	if err = Db.Model(&types.Cars{}).Order("rand()").Select(fields).Where(where).Scan(&car_cards).Error; err != nil {
		return nil, err
	}

	for i := range car_cards {
		// concatena numero aleatorio ao SRC da imagem do carro para evitar cache do browser
		car_cards[i].CarImage += "?" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	return car_cards, nil
}

// *************************************************************************************************************************
// retorna campos de determinado registro de fabricante
// *************************************************************************************************************************
func GetCar(carId int, c *fiber.Ctx) (record *types.CarResponse, err error) {

	// logo_filename ficara parecido com: 	manufacturer_000002.png, manufacturer_000041.png, etc
	// car_image ficara parecido com: 	car_000002.png, car_000041.png, etc
	_select := " cars.name, manufacturer_id, country, manufacturers.name as manufacturer_name, rental_price, cars.id, " +
		" concat(cars.workgroup, '_car_', LPAD(cars.id, 6, '0'), '.png') as car_image, year, doors, odometer, mpg, cylinders, hp, cc,  " +
		" concat(cars.workgroup, '_manufacturer_', LPAD(cars.manufacturer_id, 6, '0'), '.png') as manufacturer_logo, " +
		" ifnull(transmission_manual, false) as transmission_manual "

	join := " left join manufacturers on manufacturers.id = cars.manufacturer_id "

	// retorna so os campos necessarios
	if err = Db.Model(&types.Cars{}).Joins(join).Where("cars.id = ?", carId).Select(_select).Scan(&record).Error; err != nil {
		return nil, err
	}

	// obtem os dias em que o carro esta reservado, mas somente se na rota (url) foi passado 'first_day', 'last_day'
	if c.Params("first_day") != "" && c.Params("last_day") != "" {

		type DaysCarIsReserved struct {
			Pickup  int `json:"pickup"`
			Dropoff int `json:"dropoff"`
		}

		var days []DaysCarIsReserved

		where := " ( date_format(pickup_datetime, '%Y-%m-%d') between '@firstday' and '@lastday' or " +
			" date_format(dropoff_datetime, '%Y-%m-%d') between '@firstday' and '@lastday' )  and car_id = @carid " +
			" AND deleted_at IS null "

			// obtem todos os dias entre a data retirada e data devolucao do carro
		sql := " SELECT DAY(pickup_datetime) as pickup, DAY(dropoff_datetime) as dropoff " +
			" from bookings " +
			" where " + where

		sql = strings.ReplaceAll(sql, "@firstday", c.Params("first_day"))
		sql = strings.ReplaceAll(sql, "@lastday", c.Params("last_day"))
		sql = strings.ReplaceAll(sql, "@carid", strconv.Itoa(carId))

		Db.Raw(sql).Scan(&days)

		// se o dia encontrado de reserva do veiculo
		for _, dates := range days {

			for day := dates.Pickup; day <= dates.Dropoff; day++ {
				if record.DaysIsReserved != "" {
					record.DaysIsReserved += ","
				}
				record.DaysIsReserved += strconv.Itoa(day)
			}
		}
	}

	// concatena numero aleatorio ao SRC da imagem do carro para evitar cache do browser
	record.CarImage += "?" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	return record, nil
}

// *************************************************************************************************************************
// salva registro
// *************************************************************************************************************************
func SaveCar(request *types.CarRequest, c *fiber.Ctx, carId int) string {

	car := types.Cars{
		Name:               request.Name,
		Country:            request.Country,
		Year:               request.Year,
		ManufacturerId:     request.ManufacturerId,
		RentalPrice:        request.RentalPrice,
		Odometer:           request.Odometer,
		Mpg:                request.Mpg,
		Cylinders:          request.Cylinders,
		Hp:                 request.Hp,
		TransmissionManual: &request.TransmissionManual,
		Doors:              request.Doors,
		Cc:                 request.Cc,
		Active:             true,
		Workgroup:          c.Params("workgroup"),
	}

	var err error
	var result string
	var carId_use int

	// se der erro no upload da image, cancela alteracao/criacao do registro
	trans := Db.Begin()

	// novo registro
	if c.Method() == "POST" {
		if err = trans.Model(&types.Cars{}).Create(&car).Error; err != nil {
			trans.Rollback()
			return err.Error()
		}

		carId_use = car.Id
		result = "__success__|" + strconv.Itoa(car.Id)
	}

	// edita registro
	if c.Method() == "PATCH" {
		if err = trans.Model(&types.Cars{}).Where("id = ? ", carId).Updates(&car).Error; err != nil {
			trans.Rollback()
			return err.Error()
		}

		carId_use = carId
		result = "__success__"
	}

	// quando registro esta sendo editado, ja ha imagem gravada, a verificacao/upload
	// de nova imagem só é necessaria caso o usuario escolheu outra imagem
	// 'BypassImageUpload' é decidida pelo front end
	if !request.BypassImageUpload {
		// se der qq erro com o upload ou integridade da imagem enviada, cancela a insercao/update feito acima
		if err = utils.ValidateIAndUploadImageFile(c, carId_use); err != nil {

			fmt.Println("deu erro")
			fmt.Println(err)

			trans.Rollback()
			return err.Error()
		}
	}

	// necessario efetivar a transacao aqui, caso contrario nao sera possivel registrar notificacao
	trans.Commit()

	// registra que a operacao acima deve ser notificada a todos os usuarios do grupo, exceto o usuario que efetuou a operacao (ip cliente)
	if err = LogNotification(c, carId_use, request.ClientIp); err != nil {
		return err.Error()
	}

	return result
}

// *********************************************************************************************************************************************
// prepara adaptacao do campo 'order by' vindo do front end, requisitado pelo usuario, para que corresponda ao campo correto na base
// ** EXEMPLO: ***
//
//	usuario pediu no datatable organizacao por 'manufacturers_name',  a query dever executar 'order by manufacturers.name'
//
// *********************************************************************************************************************************************
func CarsOrderByTranslated(FieldChosenInDatatable string) string {

	var translate []OrderByTranslated
	translate = append(translate, OrderByTranslated{FieldChosenInDatatable: "manufacturer_name", FieldInQuery: "manufacturers.name"})

	corrected_field := ""
	for _, field := range translate {
		if FieldChosenInDatatable == field.FieldChosenInDatatable {
			corrected_field = field.FieldInQuery
		}
	}

	// se nao houve alteracao, o campo é da tabela 'cars', mesmo assim, concatena a string "cars." para que nao ocorra eventualmente o erro "ambiguous column name in order by"
	if corrected_field == "" {
		corrected_field = "cars." + FieldChosenInDatatable
	}
	return corrected_field
}
