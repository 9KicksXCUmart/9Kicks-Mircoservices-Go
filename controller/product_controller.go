package controller

import (
	. "9Kicks/model/product"
	"9Kicks/service/product"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FilterProducts(c *gin.Context) {
	name := c.Query("keyword")
	pageNumString := c.Query("pageNum")
	pageNum, _ := strconv.Atoi(pageNumString)
	category := c.Query("category")
	brand := c.Query("brand")
	minPriceString := c.Query("minPrice")
	maxPriceString := c.Query("maxPrice")

	products, err := product.FilterProducts(name, category, brand, minPriceString, maxPriceString, pageNum)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "Search products successfully",
		"data": gin.H{
			"amount":   len(products),
			"products": products,
		},
	})
}

func PublishProduct(c *gin.Context) {
	var formData ProductFormData
	var productInfo PublishProductForm
	if err := c.Bind(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error binding form data"})
		return
	}
	log.Println(formData.Info)
	err := json.Unmarshal([]byte(formData.Info), &productInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Wrong product information format, it must be a JSON string"})
		return
	}
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error getting the image"})
		return
	}

	success := product.CreateProduct(productInfo, *file)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product created successfully"})
}
