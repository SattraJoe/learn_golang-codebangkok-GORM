package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface
}

func (l *SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// } else {
	// 	fmt.Printf("SQL: %s, Rows affected: %d\n", sql, rowsAffected)
	// }
	fmt.Printf("%v\n==============================\n", sql)

}

func main() {
	initConfig()
	initTimeZone()
	// Load configuration from config.ymal
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=true",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.name"),
	)
	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: true,
	})
	if err != nil {
		panic(err)
	}

	// isExist := db.Migrator().HasTable(Gender{})
	// if !isExist {
	// 	db.Migrator().CreateTable(Gender{})
	// }

	// isExist = db.Migrator().HasTable(Category{})
	// if !isExist {
	// 	db.Migrator().CreateTable(Category{})
	// }

	// isExist = db.Migrator().HasTable(Person{})
	// if !isExist {
	// 	db.Migrator().CreateTable(Person{})
	// }

	// isExist = db.Migrator().HasTable(OrderDetail{})
	// if !isExist {
	// 	db.Migrator().CreateTable(OrderDetail{})
	// }

	// isExist := db.Migrator().HasTable(Test{})
	// if !isExist {
	// 	db.Migrator().CreateTable(Test{})
	// }
	// db.AutoMigrate(Test{})

	db.Migrator().CreateTable(Test{})
}

type Gender struct {
	ID   uint
	Name string
}

type Category struct {
	ID   uint
	Name string
}

type Person struct {
	ID   uint
	Name string
}

type OrderDetail struct {
	ID        uint
	Name      string
	CreatedAt string
}

type Test struct {
	ID        uint
	Name      string
	CreatedAt string
	Desc      string
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
