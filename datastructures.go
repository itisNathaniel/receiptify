package main

import (
    "time"
)

type Transaction struct {
    details         RecieptDetails
    item            []RecieptItem
}

type RecieptDetails struct {
    OrderID         string
    OrderDate       string
    OrderTime       string
    OrderDateTime   time.Time
    Table           string
    PayMethod       string
    OrderTotal      string
    VatTotal        string
    OrderWithVat    int64
    Name            string
    Address         string
    Postcode        string
    Phone           string
    VatNumber       string
}

type RecieptItem struct {
    Description     string
    Quantity        string
    Price           string
}

type MonzoTransaction struct {
    Date            time.Time
    CurrencyAmount  int64
    TransactionId   string
}