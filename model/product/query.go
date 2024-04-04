package product

type NameField struct {
	Query    string `json:"query"`
	Operator string `json:"operator"`
}

type Match struct {
	Name     *NameField `json:"name,omitempty"`
	Category string     `json:"category,omitempty"`
	Brand    string     `json:"brand,omitempty"`
	MinPrice float64    `json:"minPrice,omitempty"`
	MaxPrice float64    `json:"maxPrice,omitempty"`
}

type Price struct {
	Gt float64 `json:"gt"`
	Lt float64 `json:"lt"`
}

type Range struct {
	Price Price `json:"price"`
}

type Filter struct {
	Match *Match `json:"match,omitempty"`
	Range *Range `json:"range,omitempty"`
}

type Filters struct {
	Filters []Filter `json:"must"`
}

type BoolQuery struct {
	Bool Filters `json:"bool"`
}

type SearchQuery struct {
	From  int       `json:"from"`
	Size  int       `json:"size"`
	Query BoolQuery `json:"query"`
}
