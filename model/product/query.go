package product

type SearchQuery struct {
	From  int `json:"from"`
	Size  int `json:"size"`
	Query struct {
		Match struct {
			Name string `json:"name"`
		} `json:"match"`
	} `json:"query"`
}
