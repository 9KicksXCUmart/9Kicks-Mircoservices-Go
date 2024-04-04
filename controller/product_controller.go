package controller

import (
	"9Kicks/service/product"
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
