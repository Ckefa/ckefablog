package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Ckefa/ckefablog/models"
)

// Define struct for the response
type OrderSt struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	PaymentSource struct {
		PayPal map[string]interface{} `json:"paypal"`
	} `json:"payment_source"`
	Links []Link `json:"links"`
}

// Define struct for the Links
type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

var OrderStatus OrderSt

func CreateOrder(order *models.Order) error {
	paypalUrl := os.Getenv("PaypalUrl")
	if paypalUrl != "" {
		log.Println("<< Env varibale PaypalUrl Failed to load")
	}
	// Define the headers
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + AuthToken.Token,
	}

	// Define the JSON payload
	payload := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"reference_id": order.ID,
				"amount": map[string]interface{}{
					"currency_code": "USD",
					"value":         order.Amount,
				},
			},
		},
		"payment_source": map[string]interface{}{
			"paypal": map[string]interface{}{
				"experience_context": map[string]interface{}{
					"payment_method_preference": "IMMEDIATE_PAYMENT_REQUIRED",
					"brand_name":                "CkefaWeb Agency",
					"locale":                    "en-US",
					"landing_page":              "LOGIN",
					"shipping_preference":       "NO_SHIPPING",
					"user_action":               "PAY_NOW",
					"return_url":                fmt.Sprintf("https://www.ckefa.com/order/confirm/%s", order.ID),
					"cancel_url":                fmt.Sprintf("https://example.com/order/cancel/%s", order.ID),
				},
			},
		},
	}

	// Marshal the payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling payload:", err)
		return err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", paypalUrl+"/v2/checkout/orders", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Printf("Failed to create order. Status code: %d\n", resp.StatusCode)
		return fmt.Errorf("failed to create order, status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return err
	}

	// Unmarshal response body into OrderStatus
	if err := json.Unmarshal(body, &OrderStatus); err != nil {
		log.Println("Error unmarshalling response:", err)
		log.Println("Response body:", string(body)) // Log the raw body for debugging
		return err
	}

	// Ensure Links array has at least 2 elements
	if len(OrderStatus.Links) < 2 {
		log.Println("Insufficient links in response")
		return fmt.Errorf("unexpected response format")
	}

	// Update the order with the response details
	order.PayId = OrderStatus.ID
	order.StatusMsg = OrderStatus.Status
	order.Link1 = OrderStatus.Links[0].Href
	order.Link2 = OrderStatus.Links[1].Href

	return nil

	// sample := map[string]interface{}{
	// 	"id":     "65037881M3452232B",
	// 	"status": "PAYER_ACTION_REQUIRED",
	// 	"payment_source": [string]interface{}{
	// 		"paypal": map[string]interface{}{},
	// 	},
	// 	"links": []map[string]interface{}{
	// 		{
	// 			"href":   "https://api.sandbox.paypal.com/v2/checkout/orders/65037881M3452232B",
	// 			"rel":    "self",
	// 			"method": "GET",
	// 		},
	// 		{
	// 			"href":   "https://www.sandbox.paypal.com/checkoutnow?token=65037881M3452232B",
	// 			"rel":    "payer-action",
	// 			"method": "GET",
	// 		},
	// 	},
	// }
}
