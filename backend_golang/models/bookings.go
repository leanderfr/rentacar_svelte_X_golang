package models

import (
	"fmt"
	"rentacar/types"
	"rentacar/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// *************************************************************************************************************************
// obtem lista de reservas de veiculos feitas para determinado periodo
// *************************************************************************************************************************
func GetBookingsToPopulateSchedule(bookingFilter *types.BookingFilterRequest) (recordset *[]types.BookingsResponse, err error) {

	fields := " bookings.driver_name,  bookings.car_id, concat(cars.workgroup, '_car_', LPAD(cars.id, 6, '0'), '.png') as car_image, " +
		" if(bookings.country='usa', date_format(pickup_datetime, '%m/%d %h:%i - %p'), date_format(pickup_datetime, '%d/%m - %H:%i')) as pickup_formatted,   " +
		" if(bookings.country='usa', date_format(dropoff_datetime, '%m/%d %h:%i - %p'), date_format(dropoff_datetime, '%d/%m - %H:%i')) as dropoff_formatted,   " +
		" date_format(pickup_datetime, '%Y-%m-%d|%H:%i') as pickup_reference,   date_format(dropoff_datetime, '%Y-%m-%d|%H:%i') as dropoff_reference, " +
		" bookings.id as booking_id "

	join := " left join cars on bookings.car_id = cars.id "

	where := " (DATE_FORMAT(pickup_datetime,'%Y-%m-%d') between '@first' and '@last' or DATE_FORMAT(dropoff_datetime,'%Y-%m-%d') between '@first' and '@last') "

	// car_id = -1, usuario pediu para ver reservas de todos os carros
	if bookingFilter.CarId == -1 {
		where += utils.ConcatWhere(where, " bookings.country='@country' ")
	}

	// se usuario indicou qual carro listar
	if bookingFilter.CarId != -1 {
		where += utils.ConcatWhere(where, " car_id = @car_id ")
	}

	where += utils.ConcatWhere(where, fmt.Sprintf(" bookings.workgroup =  '%v' ", bookingFilter.Workgroup))

	where = strings.ReplaceAll(where, "@first", bookingFilter.FirstDayWeek.Format("2006-01-02"))
	where = strings.ReplaceAll(where, "@last", bookingFilter.LastDayWeek.Format("2006-01-02"))
	where = strings.ReplaceAll(where, "@country", bookingFilter.Country)
	where = strings.ReplaceAll(where, "@car_id", strconv.Itoa(bookingFilter.CarId))

	// busca somente os campos necessarios
	if err = Db.Model(&types.Bookings{}).Select(fields).Joins(join).Where(where).Scan(&recordset).Error; err != nil {
		return nil, err
	}
	return recordset, err
}

// *************************************************************************************************************************
// retorna registro de agendamento
// *************************************************************************************************************************
func GetBooking(bookingId int) (record *types.BookingResponse, err error) {

	fields := " bookings.car_id, driver_name, " +
		" if(bookings.country='usa', date_format(pickup_datetime, '%m/%d/%y'), date_format(pickup_datetime, '%d/%m/%y')) as pickup_date,   " +
		" if(bookings.country='usa', date_format(pickup_datetime, '%h:%i - %p'), date_format(pickup_datetime, '%H:%i')) as pickup_hour,   " +
		" if(bookings.country='usa', date_format(dropoff_datetime, '%m/%d/%y'), date_format(dropoff_datetime, '%d/%m/%y')) as dropoff_date,   " +
		" if(bookings.country='usa', date_format(dropoff_datetime, '%h:%i - %p'), date_format(dropoff_datetime, '%H:%i')) as dropoff_hour   "

	where := fmt.Sprintf(" id = %v ", bookingId)

	// busca somente os campos necessarios
	if err = Db.Model(&types.Bookings{}).Select(fields).Where(where).Scan(&record).Error; err != nil {
		return nil, err
	}
	return record, err
}

// *************************************************************************************************************************
// salva registro
// *************************************************************************************************************************
func SaveBooking(request *types.BookingRequest, c *fiber.Ctx, bookingId int) string {

	booking := types.Bookings{
		CarId:           request.CarId,
		PickupDatetime:  request.PickupDatetime,
		DropoffDatetime: request.DropoffDatetime,
		DriverName:      request.DriverName,
		Country:         request.Country,
		Workgroup:       c.Params("workgroup"),
	}

	var err error
	var result string
	var bookigId_use int

	// novo registro
	if c.Method() == "POST" {
		if err = Db.Model(&types.Bookings{}).Create(&booking).Error; err != nil {
			return err.Error()
		}

		result = "__success__|" + strconv.Itoa(booking.Id)
		bookigId_use = booking.Id
	}

	// edita registro
	if c.Method() == "PATCH" {
		if err = Db.Model(&types.Bookings{}).Where("id = ? ", bookingId).Updates(&booking).Error; err != nil {
			return err.Error()
		}

		result = "__success__"
		bookigId_use = bookingId
	}

	// registra que a operacao acima deve ser notificada a todos os usuarios do grupo, exceto o usuario que efetuou a operacao (ip cliente)
	if err = LogNotification(c, bookigId_use, request.ClientIp); err != nil {
		return err.Error()
	}

	return result
}
