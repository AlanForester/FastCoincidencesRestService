package srv

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
)

var sql *gorm.DB

type Postgres struct{}

func (s *Postgres) Connect() *gorm.DB {
	dsn := "user=docker password=docker sslmode=disable host=db"
	instance, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Printf("Postgres Error: %+v", err)
		if err.Error() == "pq: the database system is starting up" {
			log.Printf("Sleep and reconnect to DB")
			time.Sleep(5 * time.Second)
			s.Connect()
		}
	}
	return instance
}

func initDB() *gorm.DB {
	conn := &Postgres{}
	ormInstance := conn.Connect()
	ormInstance.BlockGlobalUpdate(false)

	sqlInstance := ormInstance.DB()
	sqlInstance.SetConnMaxLifetime(0)
	log.Printf("Postgres is ready")
	return ormInstance.Debug()
}

func SQL() *gorm.DB {
	if sql == nil {
		sql = initDB()
	}
	return sql
}
