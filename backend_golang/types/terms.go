package types

import (
	"time"

	"gorm.io/gorm"
)

// registro na base
type Terms struct {
	Id          int  `gorm:"primary_key" `
	Active      bool `gorm:"default:true"`
	Workgroup   string
	Item        string
	English     string
	Portuguese  string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// internamente a tabela chama-se 'languages'
func (Terms) TableName() string {
	return "terms"
}

// gravacao de registro
type TermRequest struct {
	Item        string `form:"item" validate:"required,min=5,max=200"`
	Portuguese  string `form:"portuguese" validate:"required,min=2,max=1200"`
	English     string `form:"english" validate:"required,min=2,max=1200"`
	Description string `form:"description" validate:"required,min=10,max=200"`
	ClientIp    string `form:"client_ip"`
}

// todos os termos, expressoes do idioma atualmente selecionado no front end
type TermsForPopulateFronEndResponse struct {
	Item       string `json:"item"`
	Expression string `json:"expression"`
}

// listagem de termos/expressoes
type TermsResponse struct {
	Id          string `json:"id"`
	Item        string `json:"item"`
	Portuguese  string `json:"portuguese"`
	English     string `json:"english"`
	Active      string `json:"active"`
	Description string `json:"description"`
}

// formulario de termo/expressao
type TermResponse struct {
	Item        string `json:"item"`
	Portuguese  string `json:"portuguese"`
	English     string `json:"english"`
	Description string `json:"description"`
}
