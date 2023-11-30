package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectToMySQL() (*gorm.DB, error) {
	dsn := "<user>:<pass>@tcp(localhost:3306)/<db>?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Admin    bool   `json:"admin"`
	Password string `json:"-"` // "-" means this field will be ignored in JSON serialization
	ApiCalls uint   `json:"apiCalls"`
}

func main() {
	db, err := connectToMySQL()
	if err != nil {
		log.Fatal(err)
	} else {
		println("SDFsdf")
	}

	// Perform database migration
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	// Serve static files (including HTML)
	router.StaticFS("/static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/users", getUsersHandler)
	router.GET("/users/:id", getUserByIDHandler)
	router.DELETE("/users/:id", deleteUserByIDHandler)
	router.POST("/create", createUserHandler)

	log.Println("Server is running on :8080")
	router.Run("localhost:8080")

}
