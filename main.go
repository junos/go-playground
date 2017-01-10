package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
)

type Users struct {
	Id        int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Firstname string `gorm:"not null" form:"firstname" json:"firstname"`
	Lastname  string `gorm:"not null" form:"lastname" json:"lastname"`
}

type Prices struct {
	Id        int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Field1 string `gorm:"not null" form:"field1" json:"field1"`
	Field2 string `gorm:"not null" form:"field2" json:"field2"`
	Field3 string `gorm:"not null" form:"field3" json:"field3"`
	Field4 string `gorm:"not null" form:"field4" json:"field4"`
	Field5 string `gorm:"not null" form:"field5" json:"field5"`
	Field6 string `gorm:"not null" form:"field6" json:"field6"`
	Field7 string `gorm:"not null" form:"field7" json:"field7"`
	Field8 string `gorm:"not null" form:"field8" json:"field8"`
	Field9 string `gorm:"not null" form:"field9" json:"field9"`
	Field10 string `gorm:"not null" form:"field10" json:"field10"`
	Field11 string `gorm:"not null" form:"field11" json:"field11"`
	Field12 string `gorm:"not null" form:"field12" json:"field12"`
	Field13 string `gorm:"not null" form:"field13" json:"field13"`
	Field14 string `gorm:"not null" form:"field14" json:"field14"`
	Field15 string `gorm:"not null" form:"field15" json:"field15"`
}

/*
const (
	DB_USER = ""
	DB_PASSWORD = ""
	DB_NAME = ""
)
*/

func InitDb() *gorm.DB {
	// Openning file
	db, err := gorm.Open("sqlite3", "./data.db")
	// Display SQL queries
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}
	// Creating the table
	if !db.HasTable(&Users{}) {
		db.CreateTable(&Users{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Users{})
	}

	return db
}

func InitPgDb() *gorm.DB {
	// Openning PG database
	db, err := gorm.Open("postgres", "host=localhost user=postgres password=mpostgres dbname=api_phoenix_dev sslmode=disable")

	// Display SQL queries
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}

	return db
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Use(Cors())

	v1 := r.Group("api/v1")
	{
		v1.POST("/users", PostUser)
		v1.GET("/users", GetUsers)
		v1.GET("/users/:id", GetUser)
		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
		v1.GET("/prices", GetPrices)
	}

	r.Run(":8080")
}

func PostUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var user Users
	c.Bind(&user)

	if user.Firstname != "" && user.Lastname != "" {
		// INSERT INTO "users" (name) VALUES (user.Name);
		db.Create(&user)
		// Display error
		c.JSON(201, gin.H{"success": user})
	} else {
		// Display error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func GetUsers(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	var users []Users
	// SELECT * FROM users
	db.Find(&users)

	// Display JSON result
	c.JSON(200, users)

	// curl -i http://localhost:8080/api/v1/users
}

func GetPrices(c *gin.Context) {
	// Connection to the database
	db := InitPgDb()
	// Close connection database
	defer db.Close()

	var prices []Prices
	db.Limit(1000000).Find(&prices)
	//db.Where("field2 = ?", "BDBD").Limit(10).Find(&prices)

	// Display JSON result
	c.JSON(200, prices)

	// curl -i http://localhost:8080/api/v1/prices
}

func GetUser(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	id := c.Params.ByName("id")
	var user Users
	// SELECT * FROM users WHERE id = 1;
	db.First(&user, id)

	if user.Id != 0 {
		// Display JSON result
		c.JSON(200, user)
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "User not found"})
	}

	// curl -i http://localhost:8080/api/v1/users/1
}

func UpdateUser(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	// Get id user
	id := c.Params.ByName("id")
	var user Users
	// SELECT * FROM users WHERE id = 1;
	db.First(&user, id)

	if user.Firstname != "" && user.Lastname != "" {

		if user.Id != 0 {
			var newUser Users
			c.Bind(&newUser)

			result := Users{
				Id:        user.Id,
				Firstname: newUser.Firstname,
				Lastname:  newUser.Lastname,
			}

			// UPDATE users SET firstname='newUser.Firstname', lastname='newUser.Lastname' WHERE id = user.Id;
			db.Save(&result)
			// Display modified data in JSON message "success"
			c.JSON(200, gin.H{"success": result})
		} else {
			// Display JSON error
			c.JSON(404, gin.H{"error": "User not found"})
		}

	} else {
		// Display JSON error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}

func DeleteUser(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	// Get id user
	id := c.Params.ByName("id")
	var user Users
	// SELECT * FROM users WHERE id = 1;
	db.First(&user, id)

	if user.Id != 0 {
		// DELETE FROM users WHERE id = user.Id
		db.Delete(&user)
		// Display JSON result
		c.JSON(200, gin.H{"success": "User #" + id + " deleted"})
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "User not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/users/1
}

func OptionsUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}
