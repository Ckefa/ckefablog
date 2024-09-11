package paypal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type AccessToken struct {
	Token string `json:"access_token"`
	AppId string `json:"app_id"`
}

var AuthToken AccessToken

func InitPayment() {
	clientID := os.Getenv("ClientID")
	clientSecret := os.Getenv("ClientSecret")

	// Create basic authentication header
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Prepare the form data
	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://api-m.sandbox.paypal.com/v1/oauth2/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the headers
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	if err := json.Unmarshal(body, &AuthToken); err != nil {
		log.Fatal(err)
	}

	// println(AuthToken.AppId, AuthToken.Token)
	log.Println("Access Authorized")
	// sample := map[string]interface{}{
	// 	"scope":        "",
	// 	"access_token": "A21AAKaEMcjp9opdL-S_5I2jtG1qbFug6qHsEYlNvmi4FZ3ShCU7SkMM_nEZCTf-Riz6PjBwel3syoEomHq2yIDmpa0m851Gw",
	// 	"token_type":   "Bearer",
	// 	"app_id":       "APP-80W284485P519543T",
	// 	"expires_in":   32400,
	// 	"nonce":        "2024-09-11T07:56:48ZnGouJNDBbUnDn5GJ79FL-3YwdklB-OuvJqbRLFospMk"}

}
