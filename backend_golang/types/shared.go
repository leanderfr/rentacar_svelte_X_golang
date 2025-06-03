package types

type DatatableParamsRequest struct {
	Country                     string `schema:"country,required"`
	OrderBy                     string `schema:"order_by,"`
	OrderDirection              string `schema:"order_direction,"`
	SearchtTxt                  string `schema:"search_txt,"`
	OnlyActiveOrInactiveRecords string `schema:"only_active_or_inactive_records,"`
}

type FilesToCloneInRepository struct {
	OriginalFile string `json:"original_file"`
	NewFile      string `json:"new_file"`
}

type FilesToDeleteInAWS struct {
	Filename string `json:"filename"`
}

type FilesToCopyFromAdmin struct {
	AdminFilename    string
	NewGroupFilename string
}

// tabela que registra insercao/edicao de registros para avisar a quem for necessario atraves do icone 'notificacoes', no front end
type Notifications struct {
	Id                    int `gorm:"primary_key"`
	Workgroup             string
	MadeByIP              string // ip que gerou a notificacao
	DescriptionEnglish    string
	DescriptionPortuguese string
}

// retorna verificacao de nova notificacao
type NotificationsResponse struct {
	Id                    string `json:"id"`
	DescriptionEnglish    string `json:"description_english"`
	DescriptionPortuguese string `json:"description_portuguese"`
}
