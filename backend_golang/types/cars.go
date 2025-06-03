package types

import (
	"time"

	"gorm.io/gorm"
)

// registro na base
type Cars struct {
	Id                 int `gorm:"primary_key"`
	Active             bool
	Country            string
	Workgroup          string
	OriginalId         int // memoriza como era ID antes de ser clonado, auxilia no processo de clonagem
	Year               string
	ManufacturerId     int
	Manufacturers      Manufacturers `gorm:"foreignKey:ManufacturerId;references:Id"`
	Name               string
	RentalPrice        string
	Odometer           string
	Mpg                string
	Cylinders          string
	Hp                 string
	TransmissionManual *bool // true= manual   false= automatica
	Doors              string
	Cc                 string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

// gravacao de registro
type CarRequest struct {
	Country            string `form:"country" validate:"required"`
	Year               string `form:"year" validate:"required,min=4,max=4"`
	Name               string `form:"name" validate:"required,min=2,max=50"`
	ManufacturerId     int    `form:"manufacturer_id" validate:"required" `
	RentalPrice        string `form:"rental_price" validate:"required"`
	Odometer           string `form:"odometer" validate:"required,min=2,max=7"`
	FileCarImage       string `form:"file_car_image"`
	BypassImageUpload  bool   `form:"bypass_image_upload"`
	Mpg                string `form:"mpg" validate:"required,min=1,max=5"`
	Cylinders          string `form:"cylinders" validate:"required,min=1,max=1"`
	Hp                 string `form:"hp" validate:"required,min=2,max=3"`
	TransmissionManual bool   `form:"transmission_manual" binding:"required,boolean"`
	Doors              string `form:"doors" validate:"required,min=1,max=1"`
	Cc                 string `form:"cc" validate:"required,min=2,max=5"`
	ClientIp           string `form:"client_ip"  validate:"required"`
}

// formulario de veiculo
type CarResponse struct {
	Country            string `json:"country"`
	Year               string `json:"year"`
	Name               string `json:"name"`
	ManufacturerName   string `json:"manufacturer_name"`
	ManufacturerLogo   string `json:"manufacturer_logo"`
	RentalPrice        string `json:"rental_price"`
	Odometer           string `json:"odometer"`
	CarImage           string `json:"car_image"`
	Mpg                string `json:"mpg"`
	Cylinders          string `json:"cylinders"`
	Hp                 string `json:"hp"`
	TransmissionManual string `json:"transmission_manual"`
	Doors              string `json:"doors"`
	Cc                 string `json:"cc"`
	DaysIsReserved     string `json:"days_reserved"` // retorna os dias do mês/ano em que o carro esta reservado no formato dia 1, dia 2, dia 3, ....
}

// listagem de veiculos
type CarsResponse struct {
	Id               int    `json:"id"`
	Active           int    `json:"active"`
	Name             string `json:"name"`
	ManufacturerLogo string `json:"manufacturer_logo"` // so o nome do arquivo, a imagem é gravada no repositorio via FTP
	CarImage         string `json:"car_image"`
}
