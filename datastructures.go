package main

type RecieptDetails struct {
    OrderID         string
    OrderDate       string
    OrderTime       string
    Table           string
    PayMethod       string
    OrderTotal      string
    VatTotal        string
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
