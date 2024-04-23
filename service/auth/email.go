/*
Package auth provides functions to send email to users for email verification and password recovery.
*/
package auth

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

func getAccessToken() (string, string) {
	baseUrl := "https://accounts.zoho.jp/oauth/v2/token?"
	url := baseUrl + config.GetTokenParams().ZOHOTokenParams
	accountNo := config.GetTokenParams().ZOHOAccountNo
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}

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
	accessToken, accountNo := getAccessToken()
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
		Subject:     "9Kicks - Email Verification",
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

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func SendResetEmailTo(email, token, name string) error {
	accessToken, accountNo := getAccessToken()
	requestUrl := "https://mail.zoho.jp/api/accounts/" + accountNo + "/messages"
	method := "POST"

	// Load the email template
	templatePath, _ := filepath.Abs("template/reset-password.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	data := struct {
		Token     string
		Email     string
		FirstName string
	}{
		Token:     token,
		Email:     email,
		FirstName: name,
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
		Subject:     "9Kicks - Password Recovery",
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

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
