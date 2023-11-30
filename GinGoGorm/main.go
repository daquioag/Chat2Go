package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


func connectToMySQL() (*gorm.DB, error) {
    dsn := "root:new_password@tcp(localhost:3306)/nestdb?charset=utf8mb4&parseTime=True&loc=Local"
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

func createUser(db *gorm.DB, user *User) error {
    result := db.Create(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func getUserByID(db *gorm.DB, userID uint) (*User, error) {
    var user User
    result := db.First(&user, userID)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func updateUser(db *gorm.DB, user *User) error {
    result := db.Save(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func deleteUser(db *gorm.DB, user *User) error {
    result := db.Delete(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func getUsersHandler(c *gin.Context) {
	db, err := connectToMySQL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}

	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func main() {
    db, err := connectToMySQL()
    if err != nil {
        log.Fatal(err)
    }else{
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

	log.Println("Server is running on :8080")
    router.Run("localhost:8080")
	
}