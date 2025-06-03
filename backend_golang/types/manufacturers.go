package types

import (
	"time"

	"gorm.io/gorm"
)

// tabela de fabricantes de carros
type Manufacturers struct {
	Id         int  `gorm:"primary_key"`
	Active     bool `gorm:"default:true"`
	Name       string
	Workgroup  string
	OriginalId int // memoriza como era ID antes de ser clonado, auxilia no processo de clonagem
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type ManufacturerRequest struct {
	Name              string `json:"form"  validate:"required,min=3,max=50"`
	BypassImageUpload bool   `form:"bypass_image_upload"`
	ClientIp          string `form:"client_ip"`
}

// usada para API JSON
type ManufacturerResponse struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	ManufacturerLogo string `json:"manufacturer_logo"`
	Active           string `json:"active"`
}
