package main

import (
	"regexp"
	"strconv"
	"strings"
)

func StripTrailing(start string) string {
	nospace := strings.TrimSpace(start)

	return nospace
}

func stringToInt(money string) int64 {
	nospace := strings.Replace(money, " ", "", -1)
	val, err := strconv.Atoi(nospace)

	if err != nil {
		panic(err)
	}

	return int64(val)
}

func stringToPence(money string) int64 {

	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		panic(err)
	}
	processedString := reg.ReplaceAllString(money, "")

	currencyVal, err := strconv.Atoi(processedString)

	if err != nil {
		panic(err)
	}

	return int64(currencyVal)
}

func createTextTransaction(transaction Transaction) string {
	var stringOut string
	stringOut = (transaction.details.Name + "\n" + transaction.details.Address + ", " + transaction.details.Postcode + "\n\n")
	for i := range transaction.item {
		stringOut = stringOut + (StripTrailing(transaction.item[i].Quantity) + " " + StripTrailing(transaction.item[i].Price) + " " + StripTrailing(transaction.item[i].Description) + "\n")
	}
	stringOut = stringOut + ("\nVAT (" + transaction.details.VatNumber + ") " + "\n\n")
	return stringOut
}
