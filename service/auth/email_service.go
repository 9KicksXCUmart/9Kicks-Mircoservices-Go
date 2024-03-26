package auth

import (
	"9Kicks/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetAccessToken() (string, string) {
	baseUrl := "https://accounts.zoho.jp/oauth/v2/token?"
	url := baseUrl + config.GetTokenParams().ZOHOTokenParams
	accountNo := config.GetTokenParams().ZOHOAccountNo
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Cookie", "066c82b352=63009310f5ee34788612bc676d7e34fb; _zcsr_tmp=8362b2b3-6c46-4a09-adb3-c5af743cb93a; iamcsr=8362b2b3-6c46-4a09-adb3-c5af743cb93a")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	token := data["access_token"]

	return fmt.Sprintf("%s", token), accountNo
}

func SendEmailTo(email string) {
	accessToken, accountNo := GetAccessToken()
	url := "https://mail.zoho.jp/api/accounts/" + accountNo + "/messages"
	method := "POST"

	//payload := strings.NewReader(`{
	//  "fromAddress": "noreply@9kicks.shop",
	//  "toAddress": "nickwkt2001@gmail.com",
	//  "subject": "Email Verification",
	//  "mailFormat": "html",
	//  "content": "9Kicks shop"
	//}`)
	payload := strings.NewReader(fmt.Sprintf(`{ "fromAddress": "noreply@9kicks.shop", "toAddress": "%s", "ccAddress": "", "subject": "Email Verification", "mailFormat": "html", "content": "9Kicks shop" }`, email))
	//fmt.Println(aod)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Add("Authorization", "Zoho-oauthtoken "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(string(body))
}
