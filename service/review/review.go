package review

import (
	"9Kicks/dao"
	"9Kicks/model/review"
	"time"

	"github.com/google/uuid"
)

func CreateReview(email string, productId string, comment string, rating int64, Anonymous bool) bool {

	reviewId := "REVIEW#" + uuid.New().String()
	productId = "PRODUCT#" + productId

	productReview := review.ProductReview{
		PK:        productId,
		SK:        reviewId,
		Email:     email,
		Comment:   comment,
		Rating:    rating,
		DateTime:  time.Now().String(), //TODO: Change to a better format and timezone
		Anonymous: Anonymous,
	}

	return dao.AddNewReview(productReview)
}

func GetReviewList(productId string) ([]review.ProductReview, bool) {
	productReviews, err := dao.GetReviewList(productId)
	if err != nil {
		return nil, false
	}
	return productReviews, true
}
