package scripts

import (
	"app/mdl"
	"github.com/Pallinder/go-randomdata"
	"log"
	"time"
)

func LoadData(countLogs int) {
	log.Printf("Loading %d items of fake data", countLogs)
	maxUserID := 3000
	var connLogs []mdl.ConnLog
	for i := 0; i <= countLogs; i++ {
		connLogs = append(connLogs, mdl.ConnLog{
			UserID: uint64(randomdata.Number(10, maxUserID)),
			IpAddr: randomdata.IpV4Address(),
			Ts:     time.Now().Unix(),
		})
		if i%300 == 0 {
			if err := mdl.BulkCreateConnLogs(connLogs); err == nil {
				log.Printf("Last %d records\n", countLogs-i)
				connLogs = []mdl.ConnLog{}
			}
		}
	}

	log.Printf("Loaded %d records", countLogs)
}
