/*
Package product provides functions to handle product related operations
*/
package product

import (
	"9Kicks/config"
	"9Kicks/dao"
	. "9Kicks/model/product"
	"9Kicks/util"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func CreateProduct(form PublishProductForm, file multipart.FileHeader) (success bool) {
	var productInfo ProductInfo
	productInfo.ID = uuid.New().String()
	productInfo.Name = form.Name
	productInfo.Category = form.Category
	productInfo.Brand = form.Brand
	productInfo.BuyCount = 0
	productInfo.PublishDate = util.GetCurrentTime()
	productInfo.Price = form.Price
	productInfo.IsDiscount = form.IsDiscount
	productInfo.DiscountPrice = form.DiscountPrice
	productInfo.SKU = form.SKU
	productInfo.Size = form.Size

	// Upload image to S3
	publicURL, success := uploadImage(file, productInfo.ID)
	productInfo.ImageURL = publicURL

	err := dao.UpdateProduct(productInfo)
	if err != nil {
		return false
	}

	return true
}

func UpdateProduct(productInfo ProductInfo) (success bool) {
	err := dao.UpdateProduct(productInfo)
	if err != nil {
		return false
	}

	return true
}

func DeleteProduct(productId string) (success bool) {
	err := dao.DeleteProduct(productId)
	if err != nil {
		return false
	}

	return true
}

func CheckRemainingStock(productId string, size string) (remaining int, success bool) {
	remainingStock, err := dao.GetRemainingStock(productId, size)
	if err != nil {
		return 0, false
	}

	return remainingStock, true
}

func UpdateStock(productId string, size string, sold int) (success bool, message string) {
	err := dao.UpdateStock(productId, size, sold)
	if err != nil {
		return false, err.Error()
	}

	return true, "Stock updated successfully"
}

func uploadImage(file multipart.FileHeader, productId string) (publicURL string, success bool) {
	fileObject, err := file.Open()

	bucketName := "9kicks-data"
	key := "images/" + productId

	s3Client := config.GetS3Client()

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        fileObject,
		ContentType: aws.String("image/jpeg"),
	})

	if err != nil {
		return "", false
	}

	publicURL = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)

	return publicURL, true
}
