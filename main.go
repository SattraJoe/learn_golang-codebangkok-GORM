package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	// db.AutoMigrate(Gender{}, Test{}, CustomerGorm{})
	// CreateGender("Male")
	// CreateGender("Female")
	// GetGenders()
	// GetGender(1)
	// GetGenderByName("Female")
	// CreateGender("xxxx")
	// UpdateGender(4, "yyyy")
	// UpdateGenderWithModel(4, "zzzz") // Using Model to update if set value to zero, it will not update that field
	// DeleteGender(4)
	// CreateTest(0, "Test 1")
	// CreateTest(0, "Test 2")
	// CreateTest(0, "Test 3")
	// DeleteTest(3)
	// GetTests()
	// DeleteTestPermanently(3)
	// CreateCustomer("Note", 2)
	GetCustomers()
}

type CustomerGorm struct {
	ID       uint
	Name     string
	GenderID uint
	Gender   Gender
}

func GetCustomers() {
	customers := []CustomerGorm{}
	// tx := db.Preload("Gender").Find(&customers)
	tx := db.Preload(clause.Associations).Find(&customers) // This will load all associations defined in the model
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customers)
	for _, v := range customers {
		fmt.Printf("%v|%v|%v\n", v.ID, v.Name, v.Gender.Name)
	}
}

func CreateCustomer(name string, genderID uint) {
	customer := CustomerGorm{Name: name, GenderID: genderID}
	tx := db.Create(&customer)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customer)
}

func CreateTest(code uint, name string) {
	test := Test{
		Code: code,
		Name: name,
	}
	db.Create(&test)
}

func GetTests() {
	test := []Test{}
	db.Find(&test)
	for _, v := range test {
		fmt.Printf("ID: %d| Code: %d | Name: %s\n", v.ID, v.Code, v.Name)
	}
}

// DeleteTest deletes a test record by its ID (soft deleting).
func DeleteTest(id uint) {
	db.Delete(&Test{}, id)
}
func DeleteTestPermanently(id uint) {
	db.Unscoped().Delete(&Test{}, id)
}

func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println("Error deleting ", tx.Error)
		return
	}
	fmt.Println("Deleted")
	GetGender(id)
}

func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println("Error getting ", tx.Error)
		return
	}
	gender.Name = name
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println("Error updating ", tx.Error)
		return
	}
	GetGender(id)
}

func UpdateGenderWithModel(id uint, name string) {
	gender := Gender{Name: name}
	// tx := db.Model(&gender).Where("id =?", id).Updates(gender)
	tx := db.Model(&gender).Where("id =@myid", sql.Named("myid", id)).Updates(gender) // Using sql.Named to bind parameter, GORM supports named arguments
	if tx.Error != nil {
		fmt.Println("Error updating ", tx.Error)
		return
	}
	GetGender(id)
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
