package controller

import (
	"9Kicks/service/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	pageNumString := c.Query("pageNum")
	pageNum, _ := strconv.Atoi(pageNumString)
	products, err := product.SearchProductsByKeyword(keyword, pageNum)
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
