package mdl

import (
	"app/helpers"
	"app/srv"
	"log"
	"time"
)

var userIps map[int64][]uint32

func LoadDuplicates() {
	log.Printf("Load duplicates in memory...")

	var countHandled = 0
	var countRecords = 0
	var maxFlows = 3

	srv.SQL().Find(&ConnLog{}).Count(&countRecords)
	if countRecords > 0 {
		var workingFlows = 0
		var recordsPack = 3000
		iters := int(countRecords / recordsPack)
		for i := 0; i <= iters; i++ {
			workingFlows++
			for workingFlows == maxFlows {
				time.Sleep(5 * time.Second)
			}

			go func(flows int) {
				if rows, err := srv.SQL().Model(ConnLog{}).Select("user_id, ip_addr").Limit(countRecords).Offset(countRecords * i).Rows(); err == nil {
					defer rows.Close()

					for rows.Next() {
						var userId int64
						var ipStr string
						_ = rows.Scan(&userId, &ipStr)

						if ip, err := helpers.Ip2long(ipStr); err == nil {
							if _, ok := userIps[userId]; !ok {
								userIps[userId] = append(userIps[userId], ip)
								log.Printf("Invalid IP: %s\n", ipStr)
							} else {
								log.Printf("[%d] Loaded IP: %s\n", countHandled, ipStr)
								countHandled++
							}
						}
					}
					flows--
				} else {
					log.Panicf("Error load duplicates in memory!")
				}
			}(workingFlows)
		}

	} else {
		log.Printf("Records not found")
	}
}

func init() {
	userIps = make(map[int64][]uint32)
}
