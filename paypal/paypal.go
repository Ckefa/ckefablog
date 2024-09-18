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
	"sync"
	"time"
)

type AccessToken struct {
	Scope     string `json:"scope"`
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	AppID     string `json:"app_id"`
	ExpiresIn int    `json:"expires_in"`
	Nonce     string `json:"nonce"`
}

var (
	AuthToken   AccessToken
	tokenExpiry time.Time
	mu          sync.Mutex
	authFile    = ".auth"
)

// LoadAuthToken reads the token and expiry from the .auth file if it exists.
func LoadAuthToken() error {
	file, err := os.Open(authFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, nothing to load
			return nil
		}
		return err
	}
	defer file.Close()

	var tokenData struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}

	if err := json.NewDecoder(file).Decode(&tokenData); err != nil {
		return err
	}

	AuthToken.Token = tokenData.Token
	tokenExpiry = time.Unix(tokenData.ExpiresAt, 0)

	return nil
}

// SaveAuthToken saves the token and expiry to the .auth file.
func SaveAuthToken() error {
	file, err := os.Create(authFile)
	if err != nil {
		return err
	}
	defer file.Close()

	tokenData := struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{
		Token:     AuthToken.Token,
		ExpiresAt: tokenExpiry.Unix(),
	}
	return json.NewEncoder(file).Encode(tokenData)
}

func InitPayment() {
	paypalUrl := os.Getenv("PaypalUrl")
	if paypalUrl != "" {
		log.Println("<< Env varibale PaypalUrl Failed to load")
	}
	mu.Lock()
	defer mu.Unlock()

	// Load existing token from file
	if err := LoadAuthToken(); err != nil {
		log.Println("Error loading auth token:", err)
	}

	// Check if the current token is valid and not expired
	if AuthToken.Token != "" && time.Now().Before(tokenExpiry) {
		log.Println("Using existing access token:", AuthToken.Token)
		return
	}

	clientID := os.Getenv("ClientID")
	clientSecret := os.Getenv("ClientSecret")

	// Create basic authentication header
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Prepare the form data
	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	// Create the HTTP request
	req, err := http.NewRequest("POST", paypalUrl+"/v1/oauth2/token", bytes.NewBufferString(form.Encode()))
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

	// Unmarshal the response into the AuthToken struct
	if err := json.Unmarshal(body, &AuthToken); err != nil {
		log.Fatal(err)
	}

	// Set the token expiry time (current time + expires_in)
	tokenExpiry = time.Now().Add(time.Duration(AuthToken.ExpiresIn) * time.Second)

	// Save the token to the file
	if err := SaveAuthToken(); err != nil {
		log.Println("Error saving auth token:", err)
	}

	log.Println("Access token obtained:", AuthToken.Token)
	log.Println("Token expires at:", tokenExpiry)
}
