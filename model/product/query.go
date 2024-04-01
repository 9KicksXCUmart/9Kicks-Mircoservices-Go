package product

type NameField struct {
	Query    string `json:"query"`
	Operator string `json:"operator"`
}

type MatchName struct {
	Name NameField `json:"name"`
}

type QueryName struct {
	Match MatchName `json:"match"`
}

type SearchNameQuery struct {
	From  int       `json:"from"`
	Size  int       `json:"size"`
	Query QueryName `json:"query"`
}
