package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getUserByIDHandler(c *gin.Context) {
	db, err := connectToMySQL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := getUserByID(db, uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func getUserByID(db *gorm.DB, userID uint) (*User, error) {
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
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

func deleteUserByIDHandler(c *gin.Context) {
	db, err := connectToMySQL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := getUserByID(db, uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = deleteUser(db, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func deleteUser(db *gorm.DB, user *User) error {
	result := db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func createUserHandler(c *gin.Context) {
	db, err := connectToMySQL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}

	var request struct {
		Email    string `json:"email" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON request to struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create the user
	err = createUser(db, request.Email, request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func createUser(db *gorm.DB, username, email, password string) error {
	// Create a new User instance with default values
	user := &User{
		Username: username,
		Email:    email,
		Password: password,
		Admin:    false, // Set default value for Admin
		ApiCalls: 0,     // Set default value for ApiCalls
	}

	// Create the user in the database
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func updateUserHandler(c *gin.Context) {
	db, err := connectToMySQL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}

	// Extract user ID from the URL parameter
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Extract updated user information from the request body
	var request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind JSON request to struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create an updatedUser instance with the provided fields
	updatedUser := &User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	// Update the user in the database
	err = updateUserByID(db, uint(userID), updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func updateUserByID(db *gorm.DB, userID uint, updatedUser *User) error {
	// Check if the user exists
	existingUser, err := getUserByID(db, userID)
	if err != nil {
		return err
	}

	// Update the fields of the existing user
	existingUser.Username = updatedUser.Username
	existingUser.Email = updatedUser.Email
	existingUser.Password = updatedUser.Password

	// Save the updated user to the database
	result := db.Save(existingUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
