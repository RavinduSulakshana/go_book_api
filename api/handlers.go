package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env files", err)
	}

	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	//migrate the schema
	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Fatal("Failed to migrate schema", err)
	}

}

func CreateBook(c *gin.Context) {
	var book Book

	//bind the request body
	if err := c.ShouldBindJSON(&book); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "invalid input", nil)
		return
	}

	DB.Create(&book)
	ResponseJSON(c, http.StatusCreated, "Book created successfully", book)
}

func GetBooks(c *gin.Context) {
	var books []Book
	DB.Find(&books)
	ResponseJSON(c, http.StatusOK, "Books retrieved successfully", books)
}

func GetBook(c *gin.Context) {
	var book Book
	if err := DB.First(&book, c.Param("id")).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Book Not Found", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "book retrieved successfully", book)
}

func UpdateBook(c *gin.Context) {
	var book Book
	if err := DB.First(&book, c.Param("id")).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "invalid Input", nil)
		return
	}
	// bind the request body
	if err := c.ShouldBindJSON(&book); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid Input", nil)
		return
	}

	DB.Save(&book)
	ResponseJSON(c, http.StatusOK, "Book updated successfully", book)
}
func DeleteBook(c *gin.Context) {
	var book Book
	if err := DB.Delete(&book, c.Param("id")).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Book Not Found", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "book deleted successfully", nil)

}
