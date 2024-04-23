/*
Package review provides the functionality to create a review for a product, get all reviews by a user and get all reviews for a product.
*/
package review

import (
	"9Kicks/dao"
	. "9Kicks/model/review"
	"9Kicks/util"
	"strings"

	"github.com/google/uuid"
)

func CreateReview(email string, productId string, comment string, rating int64, Anonymous bool) bool {
	reviewId := "REVIEW#" + uuid.New().String()
	productId = "PRODUCT#" + productId

	productReview := ProductReview{
		PK:        productId,
		SK:        reviewId,
		Email:     email,
		Comment:   comment,
		Rating:    rating,
		DateTime:  util.GetCurrentTime(),
		Anonymous: Anonymous,
	}

	return dao.AddNewReview(productReview)
}

func GetProductReviewsByUser(email string) ([]UserReview, bool) {
	var userReviews []UserReview
	reviews, err := dao.GetUserReviews(email)
	if err != nil {
		return nil, false
	}

	for _, review := range reviews {
		reviewId := strings.Split(review.SK, "#")[1]
		userReviews = append(userReviews, UserReview{
			ProductId: strings.Split(review.PK, "#")[1],
			ReviewId:  reviewId,
			Comment:   review.Comment,
			Rating:    review.Rating,
			DateTime:  review.DateTime,
			Anonymous: review.Anonymous,
		})
	}

	return userReviews, true
}

func GetProductReviewDetails(productId string) (ReviewDetails, bool) {
	var reviewDetails ReviewDetails
	var reviewList []ReviewResponse

	productReviews, err := dao.GetReviewList(productId)
	if err != nil {
		return reviewDetails, false
	}

	avgRating := calculateAverageRating(productReviews)
	ratingPercentage := calculateRatingPercentage(productReviews)

	for _, review := range productReviews {
		reviewId := strings.Split(review.SK, "#")[1]
		reviewList = append(reviewList, ReviewResponse{
			ReviewId:  reviewId,
			Email:     review.Email,
			Comment:   review.Comment,
			Rating:    review.Rating,
			DateTime:  review.DateTime,
			Anonymous: review.Anonymous,
		})
	}

	reviewDetails.AvgRating = avgRating
	reviewDetails.RatingPercentage = ratingPercentage
	reviewDetails.Reviews = reviewList

	return reviewDetails, true
}

func calculateAverageRating(productReviews []ProductReview) float64 {
	var totalRating int64
	for _, review := range productReviews {
		totalRating += review.Rating
	}
	return float64(totalRating) / float64(len(productReviews))
}

func calculateRatingPercentage(productReviews []ProductReview) RatingPercentage {
	var ratingPercentage RatingPercentage
	for _, review := range productReviews {
		switch review.Rating {
		case 5:
			ratingPercentage.FiveStar++
		case 4:
			ratingPercentage.FourStar++
		case 3:
			ratingPercentage.ThreeStar++
		case 2:
			ratingPercentage.TwoStar++
		case 1:
			ratingPercentage.OneStar++
		}
	}
	ratingPercentage.FiveStar = int64(float64(ratingPercentage.FiveStar) / float64(len(productReviews)) * 100)
	ratingPercentage.FourStar = int64(float64(ratingPercentage.FourStar) / float64(len(productReviews)) * 100)
	ratingPercentage.ThreeStar = int64(float64(ratingPercentage.ThreeStar) / float64(len(productReviews)) * 100)
	ratingPercentage.TwoStar = int64(float64(ratingPercentage.TwoStar) / float64(len(productReviews)) * 100)
	ratingPercentage.OneStar = int64(float64(ratingPercentage.OneStar) / float64(len(productReviews)) * 100)

	return ratingPercentage
}
