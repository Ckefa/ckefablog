package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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

func CreateOrder(order *models.Order) OrderSt {

	// Define the headers
	headers := map[string]string{
		"Content-Type": "application/json",
		// "PayPal-Request-Id": "7b92603e-77ed-4896-8e78-5dea2050476a",
		"Authorization": "Bearer " + AuthToken.Token,
	}

	// Define the JSON payload
	data, err := json.MarshalIndent(order, "", " ")

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://api-m.sandbox.paypal.com/v2/checkout/orders", bytes.NewBuffer([]byte(data)))

	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		fmt.Println("Order created successfully.")
	} else {
		fmt.Printf("Failed to create order. Status code: %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	// println(string(body))

	if err := json.Unmarshal(body, &OrderStatus); err != nil {
		log.Fatal(err)
	}

	fmt.Println(OrderStatus.Status)
	return OrderStatus

	// sample := map[string]interface{}{
	// 	"id":     "65037881M3452232B",
	// 	"status": "PAYER_ACTION_REQUIRED",
	// 	"payment_source": map[string]interface{}{
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
