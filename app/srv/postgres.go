package srv

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

var sql *gorm.DB

type Postgres struct{}

func (s *Postgres) Connect() *gorm.DB {
	dsn := "user=docker password=docker sslmode=disable host=db"
	instance, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Panicf("Postgres Error: %+v", err)
	}
	return instance
}

func initDB() *gorm.DB {
	conn := &Postgres{}
	ormInstance := conn.Connect()
	ormInstance.BlockGlobalUpdate(true)
	ormInstance.Debug().LogMode(true)

	sqlInstance := ormInstance.DB()
	sqlInstance.SetConnMaxLifetime(0)
	log.Printf("Postgres is ready")
	return ormInstance
}

func SQL() *gorm.DB {
	if sql == nil {
		sql = initDB()
	}
	return sql
}
