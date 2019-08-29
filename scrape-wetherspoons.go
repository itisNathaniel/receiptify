package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"golang.org/x/net/html"
)

func parseHTMLWetherspoon(stringHTML string) ([]receiptItem, receiptDetails) {
	bodyString := strings.NewReader(stringHTML)
	token := html.NewTokenizer(bodyString)
	content := []string{}

	receiptItems := []receiptItem{}
	receiptDetail := receiptDetails{}
	tokenCount := 0

	// Extract from HTML
	for {
		tt := token.Next()
		t := token.Token()
		td := t.Data
		tokenCount++
		// line 45 is where the receipt content starts
		if tt == html.ErrorToken {
			break
		} else if tt == html.TextToken && tokenCount > 45 {
			stringContent := strings.TrimSpace(td)
			if len(stringContent) > 0 {
				//fmt.Println(stringContent)
				content = append(content, stringContent)
			}
		}
	}

	receiptDetail.Name = content[0]
	receiptDetail.Address = content[1] + ", " + content[2]
	// some have more addresses - postcodes are 6-8 chars + space so need to offset
	var addressOffset int
	if len(content[3]) < 9 {
		addressOffset = 0
	} else {
		addressOffset = 1
	}
	receiptDetail.Postcode = content[3+addressOffset]
	receiptDetail.Phone = (content[4+addressOffset])[11:]
	receiptDetail.OrderID = content[6+addressOffset]
	receiptDetail.OrderDate = (content[7+addressOffset])[12:]
	receiptDetail.OrderTime = (content[8+addressOffset])[12:]
	receiptDetail.Table = content[10+addressOffset]

	var endOfItemsIndex int
	thisLineIndex := 14 + addressOffset

	for {
		if content[thisLineIndex] == "Payment Type" {
			endOfItemsIndex = thisLineIndex
			break
		} else {
			Quantity := content[thisLineIndex]
			Name := content[thisLineIndex+1]
			AdditionalInfoOrPrice := content[thisLineIndex+2]
			PriceOrNextItem := content[thisLineIndex+3]

			var thisItem receiptItem

			if !isInt(Quantity) {
				// handle blank quantity issue
				thisItem.Price = stringToPence(thisItem.Description)
				thisItem.Description = Quantity
				thisItem.Quantity = 0
				thisLineIndex = thisLineIndex + 2
			} else if strings.Contains(AdditionalInfoOrPrice, "£") {
				// handle prices
				thisItem.Quantity = stringToInt(Quantity)
				thisItem.Description = Name
				thisItem.Price = stringToPence(AdditionalInfoOrPrice)
				thisLineIndex = thisLineIndex + 3
			} else {
				thisItem.Quantity = stringToInt(Quantity)
				thisItem.Description = Name + " " + AdditionalInfoOrPrice
				thisItem.Price = stringToPence(PriceOrNextItem)
				thisLineIndex = thisLineIndex + 4
			}

			receiptItems = append(receiptItems, thisItem)

		}
	}

	// item for general info as Monzo does not natively support in receipts
	var thisDetail receiptItem
	thisDetail.Description = receiptDetail.Name + " - Table: " + receiptDetail.Table + ""
	receiptItems = append(receiptItems, thisDetail)

	// payment method
	receiptDetail.PayMethod = content[endOfItemsIndex+1]
	receiptDetail.OrderTotal = content[endOfItemsIndex+3]
	receiptDetail.VatNumber = stripTrailing(content[endOfItemsIndex+9])

	receiptDetail.VatTotal = stringToPence(content[endOfItemsIndex+5])
	receiptDetail.OrderWithVat = stringToPence(receiptDetail.OrderTotal)

	layout := "Monday, January 02, 2006 15:04"
	dateString := receiptDetail.OrderDate + " " + receiptDetail.OrderTime
	time, err := time.Parse(layout, dateString)

	receiptDetail.OrderDateTime = time

	if err != nil {
		fmt.Println(err)
	}

	return receiptItems, receiptDetail
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
