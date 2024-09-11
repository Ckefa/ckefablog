package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID            string         `json:"id"`
	Status        bool           `json:"status"`
	Intent        string         `json:"intent"`
	PurchaseUnits []PurchaseUnit `json:"purchase_units"`
	PaymentSource PaymentSource  `json:"payment_source"`
}

type PurchaseUnit struct {
	ReferenceID string   `json:"reference_id"`
	Amount      Amount   `json:"amount"`
	Shipping    Shipping `json:"shipping"`
}

type Amount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type Shipping struct {
	Address Address `json:"address"`
}

type Address struct {
	AddressLine1 string `json:"address_line_1"`
	AdminArea2   string `json:"admin_area_2"`
	AdminArea1   string `json:"admin_area_1"`
	PostalCode   string `json:"postal_code"`
	CountryCode  string `json:"country_code"`
}

type PaymentSource struct {
	PayPal PayPal `json:"paypal"`
}

type PayPal struct {
	ExperienceContext ExperienceContext `json:"experience_context"`
}

type ExperienceContext struct {
	PaymentMethodPreference string `json:"payment_method_preference"`
	BrandName               string `json:"brand_name"`
	Locale                  string `json:"locale"`
	LandingPage             string `json:"landing_page"`
	ShippingPreference      string `json:"shipping_preference"`
	UserAction              string `json:"user_action"`
	ReturnURL               string `json:"return_url"`
	CancelURL               string `json:"cancel_url"`
}

func NewOrder(amount string) *Order {
	new_id := uuid.NewString()
	return &Order{
		ID:     new_id,
		Status: false,
		Intent: "CAPTURE",
		PurchaseUnits: []PurchaseUnit{
			{
				ReferenceID: new_id,
				Amount: Amount{
					CurrencyCode: "USD",
					Value:        amount,
				},
				Shipping: Shipping{
					Address: Address{
						AddressLine1: "1234 Main St",
						AdminArea2:   "San Jose",
						AdminArea1:   "CA",
						PostalCode:   "95131",
						CountryCode:  "US",
					},
				},
			},
		},

		PaymentSource: PaymentSource{
			PayPal: PayPal{
				ExperienceContext: ExperienceContext{
					PaymentMethodPreference: "IMMEDIATE_PAYMENT_REQUIRED",
					BrandName:               "CkefaWeb Agency",
					Locale:                  "en-US",
					LandingPage:             "LOGIN",
					ShippingPreference:      "SET_PROVIDED_ADDRESS",
					UserAction:              "PAY_NOW",
					ReturnURL:               fmt.Sprintf("https://www.ckefa.com/order/confirm/%s", new_id),
					CancelURL:               fmt.Sprintf("https://www.ckefa.com/order/cancel/%s", new_id),
				},
			},
		},
	}
}
