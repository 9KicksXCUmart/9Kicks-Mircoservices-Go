package product

import (
	"9Kicks/dao"
	. "9Kicks/model/product"
)

const PageSize int = 48

func SearchProductsByKeyword(keyword string, pageNum int) ([]ProductInfo, error) {
	var ProductList []ProductInfo
	// search products by keyword with pagination
	data, err := dao.Search((pageNum-1)*PageSize, PageSize, keyword)
	if err != nil {
		return ProductList, err
	}
	for _, hits := range data.Hits.Hits {
		productInfo := hits.Source
		ProductList = append(ProductList, productInfo)
	}

	return ProductList, nil
}
