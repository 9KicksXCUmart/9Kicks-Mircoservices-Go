package product

import (
	"9Kicks/dao"
	. "9Kicks/model/product"
	"errors"
	"strconv"
)

const PageSize int = 48

func FilterProducts(name, category, brand, minPriceString, maxPriceString string, pageNum int) ([]ProductInfo, error) {
	var ProductList []ProductInfo
	var filters []Filter
	var query BoolQuery

	if name != "" {
		var match Match
		match.Name = &NameField{Query: name, Operator: "and"}
		filters = append(filters, Filter{Match: &match})
	}
	if category != "" {
		var match Match
		match.Category = category
		filters = append(filters, Filter{Match: &match})
	}
	if brand != "" {
		var match Match
		match.Brand = brand
		filters = append(filters, Filter{Match: &match})
	}
	if minPriceString != "" && maxPriceString != "" {
		minPriceFloat, _ := strconv.ParseFloat(minPriceString, 64)
		maxPriceFloat, _ := strconv.ParseFloat(maxPriceString, 64)
		if minPriceFloat > maxPriceFloat {
			return nil, errors.New("minPrice should be less than maxPrice")
		}
		var price Price
		price.Gt = minPriceFloat
		price.Lt = maxPriceFloat
		var rangeFilter Range
		rangeFilter.Price = price
		filters = append(filters, Filter{Range: &rangeFilter})
	}

	// TODO: filter by size

	query.Bool.Filters = filters
	data, err := dao.Filter((pageNum-1)*PageSize, PageSize, query)
	if err != nil {
		return ProductList, err
	}
	for _, hits := range data.Hits.Hits {
		productInfo := hits.Source
		ProductList = append(ProductList, productInfo)
	}

	return ProductList, nil
}

func GetProductDetailByID(productId string) (ProductInfo, error) {
	var productInfo ProductInfo
	document, err := dao.GetProductDetailByID(productId)
	if err != nil {
		return productInfo, err
	}
	productInfo = document.Source

	return productInfo, nil
}
