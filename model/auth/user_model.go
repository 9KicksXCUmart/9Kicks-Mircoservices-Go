package auth

type UserSignUpForm struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserLoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfile struct {
	PK                string            `dynamodbav:"PK"`
	SK                string            `dynamodbav:"SK"`
	Email             string            `dynamodbav:"email"`
	Password          []byte            `dynamodbav:"password"`
	FirstName         string            `dynamodbav:"firstName"`
	LastName          string            `dynamodbav:"lastName"`
	ShippingAddress   map[string]string `dynamodbav:"shippingAddress"`
	CreditCardDetails map[string]string `dynamodbav:"creditCardDetails"`
	IsVerified        bool              `dynamodbav:"isVerified"`
}
