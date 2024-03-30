package auth_service

import (
	"9Kicks/config"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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

func SendEmailTo(email, token string) error {
	accessToken, accountNo := GetAccessToken()
	requestUrl := "https://mail.zoho.jp/api/accounts/" + accountNo + "/messages"
	method := "POST"

	// Load the email template
	templatePath, _ := filepath.Abs("template/verify-email.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	data := struct {
		Token string
		Email string
	}{
		Token: token,
		Email: email,
	}

	// Execute the template and get the rendered HTML
	var emailBody bytes.Buffer
	err = tmpl.Execute(&emailBody, data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	type EmailContent struct {
		FromAddress string `json:"fromAddress"`
		ToAddress   string `json:"toAddress"`
		CcAddress   string `json:"ccAddress"`
		Subject     string `json:"subject"`
		MailFormat  string `json:"mailFormat"`
		Content     string `json:"content"`
	}

	emailContent := EmailContent{
		FromAddress: "noreply@9kicks.shop",
		ToAddress:   email,
		CcAddress:   "",
		Subject:     "Email Verification",
		MailFormat:  "html",
		Content:     emailBody.String(),
	}

	contentBytes, err := json.Marshal(emailContent)
	contentString := strings.Replace(string(contentBytes), `\u003c`, `<`, -1)
	contentString = strings.Replace(contentString, `\u003e`, `>`, -1)
	contentString = strings.Replace(contentString, `\u0026`, `&`, -1)
	payload := strings.NewReader(contentString)

	client := &http.Client{}
	req, err := http.NewRequest(method, requestUrl, payload)

	if err != nil {
		log.Fatal(err)
		return err
	}
	req.Header.Add("Authorization", "Zoho-oauthtoken "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(body)

	return nil
}
