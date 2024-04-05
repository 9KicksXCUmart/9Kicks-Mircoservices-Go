package product

import "mime/multipart"

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Brand string `json:"brand"`
}

type ShoeSize struct {
	Size5_5  int `json:"5.5"`
	Size6_0  int `json:"6.0"`
	Size6_5  int `json:"6.5"`
	Size7_0  int `json:"7.0"`
	Size7_5  int `json:"7.5"`
	Size8_0  int `json:"8.0"`
	Size8_5  int `json:"8.5"`
	Size9_0  int `json:"9.0"`
	Size9_5  int `json:"9.5"`
	Size10_0 int `json:"10.0"`
	Size10_5 int `json:"10.5"`
	Size11_0 int `json:"11.0"`
	Size11_5 int `json:"11.5"`
	Size12_0 int `json:"12.0"`
	Size12_5 int `json:"12.5"`
	Size13_0 int `json:"13.0"`
	Size14_0 int `json:"14.0"`
}

type ProductInfo struct {
	ID             string   `json:"id"`
	Brand          string   `json:"brand"`
	Name           string   `json:"name"`
	Category       string   `json:"category"`
	Price          float64  `json:"price"`
	DiscountPrice  float64  `json:"discountPrice"`
	IsDiscount     bool     `json:"isDiscount"`
	SKU            string   `json:"sku"`
	ImageURL       string   `json:"imageUrl"`
	DetailImageURL []string `json:"detailImageUrl"`
	PublishDate    string   `json:"publishDate"`
	Size           ShoeSize `json:"size"`
	ReviewIDList   []string `json:"reviewIdList"`
	BuyCount       int      `json:"buyCount"`
}

type InnerHitsElement struct {
	Index  string      `json:"_index"`
	ID     string      `json:"_id"`
	Score  float64     `json:"_score"`
	Source ProductInfo `json:"_source"`
}

type HitsTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type OuterHits struct {
	Total    HitsTotal          `json:"total"`
	MaxScore float64            `json:"max_score"`
	Hits     []InnerHitsElement `json:"hits"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type SearchResponse struct {
	Took     int       `json:"took"`
	TimedOut bool      `json:"timed_out"`
	Share    Shards    `json:"_shards"`
	Hits     OuterHits `json:"hits"`
}

type PublishProductForm struct {
	Brand         string   `json:"brand"`
	Name          string   `json:"name"`
	Category      string   `json:"category"`
	Price         float64  `json:"price"`
	DiscountPrice float64  `json:"discountPrice"`
	IsDiscount    bool     `json:"isDiscount"`
	SKU           string   `json:"sku"`
	Size          ShoeSize `json:"size"`
}

type ProductFormData struct {
	Info string                `form:"info" binding:"required"`
	File *multipart.FileHeader `form:"image" binding:"required"`
}
