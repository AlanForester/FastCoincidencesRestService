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

func Intersection(a, b []int64) (c []int64) {
	m := make(map[int64]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func IntersectionSQL(a, b int64) (c []string) {
	var firstUserIps []string
	srv.SQL().Model(ConnLog{}).Select("ip_addr").Where("user_id = ?", a).Pluck("ip_addr", &firstUserIps)
	var secondUserIps []string
	srv.SQL().Model(ConnLog{}).Select("ip_addr").Where("user_id = ?", b).Pluck("ip_addr", &secondUserIps)

	m := make(map[string]bool)

	for _, item := range firstUserIps {
		m[item] = true
	}

	for _, item := range secondUserIps {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func init() {
	srv.SQL().AutoMigrate(&ConnLog{}).AddIndex("conn_logs_user_id_index", "user_id")
}
