package scripts

import (
	"app/helpers"
	"app/mdl"
	"github.com/Pallinder/go-randomdata"
	"log"
	"time"
)

func LoadData(countLogs int) {
	log.Printf("Loading %d items of fake data", countLogs)
	maxUserID := 3000
	maxIPs := 10000
	var connLogs []mdl.ConnLog
	for i := 0; i < countLogs; i++ {
		randIpInt := randomdata.Number(1, maxIPs)
		ip := helpers.IntToIP(uint32(randIpInt))
		connLogs = append(connLogs, mdl.ConnLog{
			UserID: uint64(randomdata.Number(10, maxUserID)),
			IpAddr: ip,
			Ts:     time.Now().Unix(),
		})
		if i%300 == 0 || countLogs < 300 {
			if err := mdl.BulkCreateConnLogs(connLogs); err == nil {
				log.Printf("Last %d records\n", countLogs-i)
				connLogs = []mdl.ConnLog{}
			}
		}
	}

	log.Printf("Loaded %d records", countLogs)
}
