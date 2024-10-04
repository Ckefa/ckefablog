package paypal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
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
	Url          string
	Mode         string
	authToken    AccessToken
	clientid     string
	clientsecret string
	tokenExpiry  time.Time
	mu           sync.Mutex
	authFile     = ".auth"
)

const tokenExpiryBuffer = 30 * time.Second // Buffer to handle potential clock skew

// InitPayment initializes the PayPal token by loading it or generating a new one.
func InitPayment() error {
	Mode = os.Getenv("Mode")

	if Mode == "sandbox" {
		Url = os.Getenv("PaypalSandbox")
		clientid = os.Getenv("SandboxClientID")
		clientsecret = os.Getenv("SandboxClientSecret")
	} else {
		Url = os.Getenv("PaypalLive")
		clientid = os.Getenv("LiveClientID")
		clientsecret = os.Getenv("LiveClientSecret")
	}

	// Check required environment variables
	if Url == "" || clientid == "" || clientsecret == "" {
		return errors.New("PayPal URL, Client ID or Client Secret is not configured")
	}

	if err := LoadAuthToken(); err != nil {
		log.Println("Error loading auth token:", err)
		return GenerateToken() // Generate token on load failure
	}

	// Refresh token if it's expired or close to expiring
	if authToken.Token == "" || time.Now().Add(tokenExpiryBuffer).After(tokenExpiry) {
		return GenerateToken()
	}

	return nil
}

// GetAuthToken returns a valid access token.
func GetAuthToken() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if authToken.Token == "" || time.Now().Add(tokenExpiryBuffer).After(tokenExpiry) {
		if err := InitPayment(); err != nil {
			return "", err // Return error if token generation fails
		}
	}

	return authToken.Token, nil
}

// LoadAuthToken reads the token and expiry from the .auth file.
func LoadAuthToken() error {
	file, err := os.Open(authFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No error if the file does not exist yet
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

	authToken.Token = tokenData.Token
	tokenExpiry = time.Unix(tokenData.ExpiresAt, 0)

	log.Println("Loaded existing access token:", authToken.Token)
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
		Token:     authToken.Token,
		ExpiresAt: tokenExpiry.Unix(),
	}

	return json.NewEncoder(file).Encode(tokenData)
}

// GenerateToken fetches a new PayPal access token and retries on failure with exponential backoff.
func GenerateToken() error {
	log.Println("<< Generating new AccessToken >>")

	mu.Lock()
	defer mu.Unlock()

	var attempt int
	maxRetries := 3
	delay := 1 * time.Second

	for attempt = 0; attempt < maxRetries; attempt++ {
		err := generateNewToken()
		if err == nil {
			return nil // Token generated successfully
		}

		log.Printf("Error generating token, retrying in %s: %v", delay, err)
		time.Sleep(delay)
		delay *= 2 // Exponential backoff
	}

	return errors.New("failed to generate access token after retries")
}

// generateNewToken sends the request to fetch a new access token.
func generateNewToken() error {
	auth := base64.StdEncoding.EncodeToString([]byte(clientid + ":" + clientsecret))

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", Url+"/v1/oauth2/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to get token: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &authToken); err != nil {
		return err
	}

	tokenExpiry = time.Now().Add(time.Duration(authToken.ExpiresIn) * time.Second)

	if err := SaveAuthToken(); err != nil {
		log.Println("Error saving auth token:", err)
	}

	log.Printf("Access token obtained: %s", authToken.Token)
	log.Printf("Token expires at: %v", tokenExpiry)

	return nil
}
