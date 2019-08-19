package main

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	details receiptDetails
	item    []receiptItem
}

type receiptDetails struct {
	OrderID       string
	OrderDate     string
	OrderTime     string
	OrderDateTime time.Time
	Table         string
	PayMethod     string
	OrderTotal    string
	VatTotal      int64
	OrderWithVat  int64
	Name          string
	Address       string
	Postcode      string
	Phone         string
	VatNumber     string
}

type receiptItem struct {
	Description string
	Quantity    string
	Price       string
}

type MonzoTransaction struct {
	Date           time.Time
	CurrencyAmount int64
	TransactionId  string
}

type receipt_Tax struct {
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	TaxNumber   string `json:"tax_number"`
}

type receipt_Merchant struct {
	Name          string `json:"name"`
	Online        bool   `json:"online"`
	StorePhone    string `json:"phone"`
	StoreName     string `json:"store_name"`
	StorePostcode string `json:"store_postcode"`
}

type receipt struct {
	TransactionID string             `json:"transaction_id"`
	ExternalID    string             `json:"external_id"`
	Total         int64              `json:"total"`
	Currency      string             `json:"currency"`
	Items         []MonzoreceiptItem `json:"items"`
	Merchant      receipt_Merchant   `json:"merchant"`
	receiptTaxes  []receipt_Tax      `json:"taxes"`
}

type MonzoreceiptItem struct {
	Description string `json:"description"`
	Quantity    int64  `json:"quantity"`
	Unit        string `json:"unit"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
}

type RawTransaction struct {
	Transaction
	Merchant json.RawMessage `json:"merchant"`
}
