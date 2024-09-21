package paypal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

// InitPayment initializes the PayPal token by loading it or generating a new one.
func InitPayment() {
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

	if Url == "" {
		log.Println("Paypal URL not configured")
		return
	}

	if err := LoadAuthToken(); err != nil {
		log.Println("Error loading auth token:", err)
		GenerateToken()
	} else if authToken.Token == "" || time.Now().After(tokenExpiry) {
		GenerateToken()
	}
}

// GetAuthToken returns a valid access token.
func GetAuthToken() string {
	mu.Lock()
	defer mu.Unlock()

	if authToken.Token == "" || time.Now().After(tokenExpiry) {
		GenerateToken() // Refresh the token
	} else {
		log.Println("Using existing access token:", authToken.Token)
	}

	return authToken.Token
}

// LoadAuthToken reads the token and expiry from the .auth file.
func LoadAuthToken() error {
	file, err := os.Open(authFile)
	if err != nil {
		if os.IsNotExist(err) {
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

	authToken.Token = tokenData.Token
	tokenExpiry = time.Unix(tokenData.ExpiresAt, 0)

	log.Println("Loading existing access token:", authToken.Token)

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

// GenerateToken fetches a new PayPal access token.
func GenerateToken() {
	log.Println("<< Generating new AccessToken >>")

	mu.Lock()
	defer mu.Unlock()

	// Create basic authentication header
	auth := base64.StdEncoding.EncodeToString([]byte(clientid + ":" + clientsecret))

	// Prepare the form data
	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	// Create the HTTP request
	req, err := http.NewRequest("POST", Url+"/v1/oauth2/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return
	}

	// Set headers
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read and unmarshal response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
		return
	}

	// Unmarshal the response into the AuthToken struct
	if err := json.Unmarshal(body, &authToken); err != nil {
		log.Fatalf("Error unmarshaling response: %v", err)
		return
	}

	// Set token expiry time
	tokenExpiry = time.Now().Add(time.Duration(authToken.ExpiresIn) * time.Second)

	// Save the token to the file
	if err := SaveAuthToken(); err != nil {
		log.Println("Error saving auth token:", err)
	}

	log.Printf("Access token obtained: %s", authToken.Token)
	log.Printf("Token expires at: %v", tokenExpiry)
}
