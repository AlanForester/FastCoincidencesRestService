package mdl

import (
	"app/srv"
	"fmt"
	"strings"
)

type ConnLog struct {
	UserID uint64 `gorm:"type:bigint" json:"user_id"`
	IpAddr string `gorm:"type:varchar(15)" json:"ip_addr"`
	Ts     int64  `gorm:"default:0" json:"ts"`
}

func BulkCreateConnLogs(rs []ConnLog) error {
	valueStrings := []string{}
	valueArgs := []interface{}{}

	for _, f := range rs {
		valueStrings = append(valueStrings, "(?, ?, ?)")

		valueArgs = append(valueArgs, f.UserID)
		valueArgs = append(valueArgs, f.IpAddr)
		valueArgs = append(valueArgs, f.Ts)
	}

	smt := `INSERT INTO conn_logs(user_id, ip_addr, ts) VALUES %s `

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	tx := srv.SQL().Begin()
	if err := tx.Exec(smt, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func init() {
	srv.SQL().AutoMigrate(&ConnLog{})
}
