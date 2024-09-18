package paypal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type OrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func CheckOrderStatus(payID string) OrderResponse {
	paypalUrl := os.Getenv("paypalUrl")
	if paypalUrl != "" {
		log.Println("<< Env varibale PaypalUrl Failed to load")
	}

	// Set up the URL for the PayPal order
	url := fmt.Sprintf("%s/v2/checkout/orders/%s", paypalUrl, payID)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+AuthToken.Token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// println(string(body))

	// Unmarshal JSON response into struct
	var order OrderResponse
	if err := json.Unmarshal(body, &order); err != nil {
		log.Fatal("Error unmarshalling response:", err)
	}

	// Output the order status
	fmt.Printf("Order ID: %s, Status: %s\n", order.ID, order.Status)
	return order
}
