package types

import (
	"time"

	"gorm.io/gorm"
)

type Workgroups struct {
	Id                    int `gorm:"primary_key" `
	Active                bool
	Name                  string
	InUse                 int
	ClientIp              string
	ClientCountry         string
	ClientCity            string
	ClientLoc             string
	ClientHostname        string
	ClientOrg             string
	ClientPostal          string
	ClientRegion          string
	ClientTimezone        string
	DatabaseChangesAmount int
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

type WorkgroupResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type WorkgroupReport struct {
	DatabaseChangesAmount int `json:"database_changes_amount"`
}

type CloningWokgroupStatus struct {
	Status          string `json:"status"`
	PercentReady    string `json:"percent_ready"`
	ChosenWorkgroup string `json:"chosen_workgroup"`
}

type WorkgroupsResponse struct {
	Id            int    `json:"id"`
	Active        string `json:"active"`
	Name          string `json:"name"`
	InUse         string `json:"in_use"`
	ClientIp      string `json:"client_ip"`
	ClientCountry string `json:"client_country"`
	ClientCity    string `json:"client_city"`
	ClientLoc     string `json:"client_loc"`
	UpdatedAt     string `json:"updated_at"`
	DeletedAt     string `json:"deleted_at"`
}

type ClientInfo struct {
	Ip       string `json:"ip"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Hostname string `json:"hostname"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Region   string `json:"region"`
	Timezone string `json:"timezone"`
}
