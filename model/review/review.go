package review

type ProductReview struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	Email     string `dynamodbav:"email"`
	Comment   string `dynamodbav:"comment"`
	Rating    int64  `dynamodbav:"rating"`
	DateTime  string `dynamodbav:"reviewDateTime"`
	Anonymous bool   `dynamodbav:"anonymous"`
}

type AddReviewForm struct {
	Email     string `json:"email"`
	ProductId string `json:"productId"`
	Comment   string `json:"comment"`
	Rating    int64  `json:"rating"`
	Anonymous bool   `json:"anonymous"`
}

type ReviewResponse struct {
	ReviewId  string `json:"reviewId"`
	Email     string `json:"email"`
	Comment   string `json:"comment"`
	Rating    int64  `json:"rating"`
	DateTime  string `json:"reviewDateTime"`
	Anonymous bool   `json:"anonymous"`
}

type RatingPercentage struct {
	FiveStar  int64 `json:"fiveStar"`
	FourStar  int64 `json:"fourStar"`
	ThreeStar int64 `json:"threeStar"`
	TwoStar   int64 `json:"twoStar"`
	OneStar   int64 `json:"oneStar"`
}

type ReviewDetails struct {
	AvgRating        float64          `json:"averageRating"`
	RatingPercentage RatingPercentage `json:"ratingPercentage"`
	Reviews          []ReviewResponse `json:"reviews"`
}

type UserReview struct {
	ProductId string `json:"productId"`
	ReviewId  string `json:"reviewId"`
	Comment   string `json:"comment"`
	Rating    int64  `json:"rating"`
	DateTime  string `json:"reviewDateTime"`
	Anonymous bool   `json:"anonymous"`
}
