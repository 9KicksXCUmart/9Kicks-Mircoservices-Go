package review

type ProductReview struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	Email     string `dynamodbav:"email"`
	Comment   string `dynamodbav:"comment"`
	Rating    int64  `dynamodbav:"rating"`
	DateTime  string `dynamodbav:"ReviewDateTime"`
	Anonymous bool   `dynamodbav:"anonymous"`
}

type AddReviewForm struct {
	Email     string `json:"email"`
	ProductId string `json:"productId"`
	Comment   string `json:"comment"`
	Rating    int64  `json:"rating"`
	Anonymous bool   `json:"anonymous"`
}
