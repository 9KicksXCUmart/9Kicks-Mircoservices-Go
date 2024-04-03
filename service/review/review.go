package review

import (
	"9Kicks/dao"
	"9Kicks/model/review"
	"github.com/google/uuid"
	"time"
)

func CreateReview(email string, productId string, comment string, rating int64, Anonymous bool) bool {

	reviewId := "REVIEW#" + uuid.New().String()

	productReview := review.ProductReview{
		PK:        productId,
		SK:        reviewId,
		Email:     email,
		Comment:   comment,
		Rating:    rating,
		DateTime:  time.Now().String(),
		Anonymous: Anonymous,
	}
	return dao.AddNewReview(productReview)

}

func GetReviewList(productId string) []review.ProductReview {
	return dao.GetReviewList(productId)
}
