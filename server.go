package main

import (
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

type Product struct {
	gorm.Model
	Name  string
	Price int
}

var db *gorm.DB

func initDB() {
    var err error
    db, err = gorm.Open(sqlite.Open("products.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    db.AutoMigrate(&Product{})

    products := []Product{
        {Name: "Laptop", Price: 1500},
        {Name: "Smartphone", Price: 800},
        {Name: "Klawiatura", Price: 100},
    }

    for _, p := range products {
        db.Create(&p)
    }
}

func getProducts(c echo.Context) error {
	var products []Product
	db.Find(&products)
	return c.JSON(http.StatusOK, products)
}

func getProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Product not found")
	}
	return c.JSON(http.StatusOK, product)
}

func createProduct(c echo.Context) error {
	var product Product
	if err := c.Bind(&product); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	db.Create(&product)
	return c.JSON(http.StatusCreated, product)
}

func updateProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Product not found")
	}
	if err := c.Bind(&product); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	db.Save(&product)
	return c.JSON(http.StatusOK, product)
}

func deleteProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Product not found")
	}
	db.Delete(&product)
	return c.NoContent(http.StatusNoContent)
}



func main() {
	initDB()

	e := echo.New()

	e.GET("/products", getProducts)
	e.GET("/products/:id", getProduct)
	e.POST("/products", createProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)

	fmt.Println("Server is running...")
	e.Logger.Fatal(e.Start(":1323"))
}
