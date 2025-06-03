package handlers

import (
	"rentacar/models"
	"rentacar/types"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// *************************************************************************************************************************
// obtem reservas feitas de veiculos para determinada semana
// *************************************************************************************************************************
func GetBookingsToPopulateSchedule(c *fiber.Ctx) error {
	var err error

	//ip := utils.ClientIp(c)

	var firstDayWeek time.Time
	var lastDayWeek time.Time
	var carId int
	country := c.Params("country")
	workgroup := c.Params("workgroup")

	if firstDayWeek, err = time.Parse("2006-01-02", c.Params("first_day_week")); err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if lastDayWeek, err = time.Parse("2006-01-02", c.Params("last_day_week")); err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if carId, err = strconv.Atoi(c.Params("car_id")); err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if country != "usa" && country != "brazil" {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString("Missing country info")
	}

	if strings.TrimSpace(workgroup) == "" {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}

	var data *[]types.BookingsResponse
	var bookingFilter = types.BookingFilterRequest{
		CarId:        carId,
		FirstDayWeek: firstDayWeek,
		LastDayWeek:  lastDayWeek,
		Country:      country,
		Workgroup:    workgroup,
	}
	if data, err = models.GetBookingsToPopulateSchedule(&bookingFilter); err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if data == nil {
		emptyJSON := make([]types.BookingsResponse, 0)
		c.JSON(emptyJSON)
	} else {
		c.JSON(*data)
	}
	return nil
}

// *************************************************************************************************************************
// grava reserva de veiculo
// *************************************************************************************************************************
func SaveBooking(c *fiber.Ctx) error {

	// verifica se todos campos/seus tipos vieram ok
	var request types.BookingRequest
	var err error

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if strings.TrimSpace(c.Params("workgroup")) == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("Missing workgroup")
	}

	if err = c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// verifica se campos (text/textarea) obrigatorios vieram, tamanhos e quantidades min/maxima, caso alguem tente injetar campos "gigantes", estando logado
	validate := validator.New()

	if err = validate.Struct(&request); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusInternalServerError).SendString(errors.Error())
	}

	// c.Method()  ==>  POST (novo reg) ou PATCH (editar reg)
	bookingId := 0 // POST
	if c.Method() == "PATCH" {
		if bookingId, err = strconv.Atoi(c.Params("id")); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	// verifica se o veiculo já esta reservado no(s) dia(s)/horario(s)

	pickupDate := request.PickupDatetime.Format("2006-01-02 15:04")
	dropoffDate := request.DropoffDatetime.Format("2006-01-02 15:04")

	//  - INTERVAL 1 MINUTE / + INTERVAL 1 MINUTE   para que reservas possam ser limitrofes, por exemplo:
	// carro foi reservado entre 08:00 e 15:00  por detemrinado usuario e reservado entre 15:00 e 16:00 por outro usuario, por isso essa tolerancia de 1 minuto
	where :=
		" ( '@pickupDate' between date_format(pickup_datetime + INTERVAL 1 MINUTE, '%Y-%m-%d %H:%i') and date_format(dropoff_datetime - INTERVAL 1 MINUTE, '%Y-%m-%d %H:%i') or " +
			"   '@dropoffDate' between date_format(pickup_datetime + INTERVAL 1 MINUTE, '%Y-%m-%d %H:%i') and date_format(dropoff_datetime - INTERVAL 1 MINUTE, '%Y-%m-%d %H:%i') or  " +
			"   date_format(pickup_datetime + INTERVAL 1 MINUTE, '%Y-%m-%d %H:%i') between '@pickupDate' and '@dropoffDate' or " +
			"   date_format(dropoff_datetime + INTERVAL 1 MINUTE, '%Y-%m-%d %H:%i') between '@pickupDate' and '@dropoffDate' ) " +
			" and " +
			"  bookings.car_id = @carId and " +
			"  bookings.id <> @bookingId  and " +
			"  bookings.workgroup = '@workgroup' "

	where = strings.ReplaceAll(where, "@pickupDate", pickupDate)
	where = strings.ReplaceAll(where, "@dropoffDate", dropoffDate)
	where = strings.ReplaceAll(where, "@bookingId", strconv.Itoa(bookingId))
	where = strings.ReplaceAll(where, "@carId", strconv.Itoa(request.CarId))
	where = strings.ReplaceAll(where, "@workgroup", c.Params("workgroup"))

	var existingBookings []types.ExistingBookings

	if err = models.Db.Model(&types.Bookings{}).Select("bookings.id as booking_id").Where(where).Scan(&existingBookings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	contBook := 0
	for range existingBookings {
		contBook++
	}

	if contBook > 0 {
		var _existingBookings_ string
		if request.Country == "usa" {
			_existingBookings_ = "There is already a booking for this car, in this period of time."
		}
		if request.Country == "brazil" {
			_existingBookings_ = "Já há uma reserva para este carro, neste período de tempo."
		}

		return c.Status(fiber.StatusInternalServerError).SendString(_existingBookings_)
	}

	try := models.SaveBooking(&request, c, bookingId)

	if !strings.Contains(try, "__success__") {
		return c.Status(fiber.StatusInternalServerError).SendString(try) // erro
	} else {
		return c.Status(fiber.StatusOK).SendString(try) // gravacao ok
	}
}

// *************************************************************************************************************************
// lê registro de determinada reserva
// *************************************************************************************************************************
func GetBooking(c *fiber.Ctx) error {

	var bookingId int
	var err error
	var data *types.BookingResponse

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	if bookingId, err = strconv.Atoi(c.Params("id")); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if data, err = models.GetBooking(bookingId); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// se nenhum registro retornado
	if data == nil {
		emptyJSON := make([]types.BookingResponse, 0)
		c.JSON(emptyJSON)

	} else {
		c.JSON(*data)
	}
	return nil

}
