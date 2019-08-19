package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tjvr/go-monzo"
)

func monzoFunc() []MonzoTransaction {
	cl := monzo.Client{
		BaseURL:     "https://api.monzo.com",
		AccessToken: os.Getenv("accesstoken"),
	}

	accounts, err := cl.Accounts("uk_retail")
	if err != nil {
		panic(err)
	}
	if len(accounts) == 0 {
		panic("no retail account")
	}
	acc := accounts[0]

	transactions, err := cl.Transactions(acc.ID, true) // don't expandMerchant
	if err != nil {
		fmt.Println(err)
	}
	i := 0

	var transactionList []MonzoTransaction

	layout := "2006-01-02T15:04:05.000Z"
	layoutAlt := "2006-01-02T15:04:05.00Z"
	for {
		var transaction MonzoTransaction
		if i < len(transactions) {
			if (transactions[i].Merchant.Name == "JD Wetherspoon" || transactions[i].Merchant.Name == "Trainline") && transactions[i].Merchant.IsOnline == true {
				thisTime, err := time.Parse(layout, transactions[i].Created)
				if err != nil {
					err = nil
					thisTime, err = time.Parse(layoutAlt, transactions[i].Created)
				}

				transaction.Date = thisTime
				transaction.CurrencyAmount = -transactions[i].Amount
				transaction.TransactionId = transactions[i].ID

				transactionList = append(transactionList, transaction)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			break
		}
		i++
	}

	return transactionList
}

func matchTransactionsMonzo(monzoTransact []MonzoTransaction, transactions []Transaction, hourTolerance float64, MerchantName string) {

	for i := range monzoTransact {
		monzoDate := monzoTransact[i].Date
		monzoTransaction := monzoTransact[i]
		for j := range transactions {
			transactionSearch := transactions[j]

			if monzoTransaction.CurrencyAmount == transactionSearch.details.OrderWithVat {
				transactionDate := transactionSearch.details.OrderDateTime
				// check the date
				diff := transactionDate.Sub(monzoDate)

				// generate an ID (should be replicated between runs to not create multiple reciepts for one event)
				identifier := strconv.FormatInt(monzoDate.Unix(), 10) + "-" + transactionSearch.details.Name
				identifier = strings.Replace(identifier, " ", "-", -1)

				// time tolerance depends on merchant
				if diff.Hours() < hourTolerance && diff.Hours() > 0 {

					var thisTransaction = monzoTransaction
					var items []MonzoRecieptItem

					for k := range transactionSearch.item {
						item := transactionSearch.item[k]
						var modelItem MonzoRecieptItem

						nospace := strings.Replace(item.Price, " ", "", -1)
						stringVal := strings.ReplaceAll(nospace, ".", "")
						stringVal = strings.ReplaceAll(stringVal, "£", "")

						price, err := strconv.ParseInt(stringVal, 10, 64)
						if err != nil {
							fmt.Println(err)
						}

						modelItem.Description = item.Description
						modelItem.Quantity = stringToInt(item.Quantity)
						modelItem.Amount = price
						modelItem.Currency = "GBP"

						items = append(items, modelItem)
					}

					// PASS and make the monzo reciept

					// Merchant
					var merchant Reciept_Merchant
					merchant.Name = MerchantName
					merchant.Online = false
					merchant.StoreName = transactionSearch.details.Name
					merchant.StorePhone = transactionSearch.details.Phone
					merchant.StorePostcode = transactionSearch.details.Postcode

					// Tax
					var taxes []Reciept_Tax

					var tax Reciept_Tax
					tax.Description = "VAT"
					tax.Currency = "GBP"
					tax.Amount = transactionSearch.details.VatTotal
					tax.TaxNumber = transactionSearch.details.VatNumber
					taxes = append(taxes, tax)

					reciept := Reciept{
						TransactionID: thisTransaction.TransactionId,
						ExternalID:    identifier,
						Total:         thisTransaction.CurrencyAmount,
						Currency:      "GBP",
						Items:         items,
						Merchant:      merchant,
						RecieptTaxes:  taxes,
					}

					fmt.Println("Match Found ✅")

					res, err := addReciept(reciept)

					fmt.Println(res, err)

				}

			}

		}
	}
}

func addReciept(reciept Reciept) (string, error) {

	json, err := json.Marshal(reciept)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	url := "https://api.monzo.com/transaction-receipts"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("accesstoken"))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)

	return "done", err
}
