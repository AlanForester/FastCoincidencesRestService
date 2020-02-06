package mdl

import (
	"app/srv"
)

type ConnLog struct {
	UserID uint64 `gorm:"type:bigint" json:"user_id"`
	IpAddr string `gorm:"type:varchar(15)" json:"ip_addr"`
	Ts     int    `gorm:"default:0" json:"ts"`
}

func init() {
	srv.SQL().AutoMigrate(&ConnLog{})
}
