package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Model untuk database
type Item struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// Inisialisasi database
var db *gorm.DB
var err error

func main() {
	// Koneksi ke SQLite database
	db, err = gorm.Open(sqlite.Open("items.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Migrasi model Item ke dalam database
	db.AutoMigrate(&Item{})

	// Inisialisasi Echo
	e := echo.New()

	// Endpoint CRUD
	e.POST("/items", createItem)    // Create
	e.GET("/items", getItems)       // Read All
	e.GET("/items/:id", getItem)    // Read One
	e.PUT("/items/:id", updateItem) // Update
	e.DELETE("/items/:id", deleteItem) // Delete

	// Jalankan server di localhost:8080
	e.Logger.Fatal(e.Start(":8080"))
}

// Create item
func createItem(c echo.Context) error {
	item := new(Item)
	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db.Create(&item)
	return c.JSON(http.StatusCreated, item)
}

// Read all items
func getItems(c echo.Context) error {
	var items []Item
	db.Find(&items)
	return c.JSON(http.StatusOK, items)
}

// Read one item by ID
func getItem(c echo.Context) error {
	id := c.Param("id")
	var item Item
	if result := db.First(&item, id); result.Error != nil {
		return c.JSON(http.StatusNotFound, "Item not found")
	}
	return c.JSON(http.StatusOK, item)
}

// Update item
func updateItem(c echo.Context) error {
	id := c.Param("id")
	var item Item
	if result := db.First(&item, id); result.Error != nil {
		return c.JSON(http.StatusNotFound, "Item not found")
	}

	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db.Save(&item)
	return c.JSON(http.StatusOK, item)
}

// Delete item
func deleteItem(c echo.Context) error {
	id := c.Param("id")
	var item Item
	if result := db.First(&item, id); result.Error != nil {
		return c.JSON(http.StatusNotFound, "Item not found")
	}
	db.Delete(&item)
	return c.JSON(http.StatusOK, "Item deleted")
}
