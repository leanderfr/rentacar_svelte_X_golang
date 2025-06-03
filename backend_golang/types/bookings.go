package types

import (
	"time"

	"gorm.io/gorm"
)

// registro na base
type Bookings struct {
	Id              int `gorm:"primary_key"`
	Workgroup       string
	CarId           int
	Cars            Cars `gorm:"foreignKey:CarId;references:Id"`
	PickupDatetime  time.Time
	DropoffDatetime time.Time
	DriverName      string
	Country         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// gravacao de registro
type BookingRequest struct {
	CarId           int       `form:"car_id" `
	PickupDatetime  time.Time `form:"pickup_datetime"`
	DropoffDatetime time.Time `form:"dropoff_datetime"`
	DriverName      string    `form:"driver_name"`
	Country         string    `form:"country"`
	ClientIp        string    `form:"client_ip"`
}

// filtro para montagem da agenda do veiculo
type BookingFilterRequest struct {
	Workgroup    string    `form:"workgroup"`
	CarId        int       `form:"car_id" `
	FirstDayWeek time.Time `form:"first_day_week"`
	LastDayWeek  time.Time `form:"last_day_week"`
	Country      string    `form:"country"`
}

// listagem de reservas
type BookingsResponse struct {
	BookingId        int    `json:"booking_id"`
	CarId            int    `json:"car_id"`
	PickupFormatted  string `json:"pickup_formatted"`
	PickupReference  string `json:"pickup_reference"`
	DropoffFormatted string `json:"dropoff_formatted"`
	DropoffReference string `json:"dropoff_reference"`
	DriverName       string `json:"driver_name"`
	CarImage         string `json:"car_image"`
}

// formulario de edicao do registro
type BookingResponse struct {
	CarId       int    `json:"car_id"`
	PickupDate  string `json:"pickup_date"`
	PickupHour  string `json:"pickup_hour"`
	DropoffDate string `json:"dropoff_date"`
	DropoffHour string `json:"dropoff_hour"`
	DriverName  string `json:"driver_name"`
}

// verifica se ha reservas em determinado horario para determinado carro
type ExistingBookings struct {
	BookingId int `json:"booking_id"`
}
