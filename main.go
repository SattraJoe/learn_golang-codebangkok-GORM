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

var db *gorm.DB

// main function initializes the configuration, sets the timezone, and connects to the database.
// It also creates the necessary tables if they do not exist.
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

	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false, // Set to true to enable dry run mode
	})
	if err != nil {
		panic(err)
	}

	// db.Migrator().CreateTable(Test{})

	// db.AutoMigrate(Gender{}, Test{})
	// CreateGender("Male")
	// CreateGender("Female")
	// GetGenders()
	// GetGender(1)
	GetGenderByName("Female")
}

func GetGenders() {
	genders := []Gender{}
	tx := db.Order("id").Find(&genders)
	if tx.Error != nil {
		fmt.Println("Error getting	", tx.Error)
		return
	}
	fmt.Println(genders)
}

func GetGenderByName(name string) {
	genders := []Gender{}
	// tx := db.Order("id").Find(&genders, "name = ?", name)
	tx := db.Where("name=?", name).Order("id").Find(&genders)
	if tx.Error != nil {
		fmt.Println("Error getting	", tx.Error)
		return
	}
	fmt.Println(genders)
}

func GetGender(id uint) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println("Error getting ", tx.Error)
		return
	}
	fmt.Println(gender)
}

func CreateGender(name string) {
	gender := Gender{Name: name}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println("Error creating ", tx.Error)
		return
	}
	fmt.Println(gender)
}

type Gender struct {
	ID   uint
	Name string `gorm:"unique;size(10)"`
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
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
	Code       uint   `gorm:"comment: 'This is Code'"`
	Name       string `gorm:"column:myname;type:varchar(50);unique;default:'Hello';not null"`
	// Price     int    `gorm:"default:1000"`
	CreatedAt string
	Desc      string
}

func (t Test) TableName() string {
	return "MyTest"
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
