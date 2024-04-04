package controller

import (
	. "9Kicks/model/review"
	"9Kicks/service/review"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddReview(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		return
	}
	fmt.Printf("Request Body: %s\n", string(bodyBytes))

	// Print request headers
	fmt.Printf("Request Headers: %v\n", c.Request.Header)

	// Ensure the request body can be read again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var newReview AddReviewForm
	if err := c.ShouldBindJSON(&newReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error()})
		return
	}

	success := review.CreateReview(newReview.Email, newReview.ProductId, newReview.Comment, newReview.Rating, newReview.Anonymous)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Review created successfully"})
}

func GetReviewList(c *gin.Context) {
	productId := c.Query("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Product ID is required"})
		return
	}

	reviews, success := review.GetProductReviewDetails(productId)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Reviews retrieved successfully",
		"data":    reviews})
}
