package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/jprobinson/eazye"
)

var procs = runtime.NumCPU()

var monzoOutput []MonzoTransaction

func main() {

	wetherspoonEmails := getMail("FROM orders@jdwetherspoon.co.uk")
	trainlineEmails := getMail(`TEXT "Your booking confirmation. Transaction Id" FROM "auto-confirm@info.thetrainline.com"`)

	//print to console how many reciepts we've got
	emailCount := len(wetherspoonEmails)
	fmt.Println(emailCount)

	// fetch monzo transactions
	monzoOutput = monzoFunc()

	// parse all trainline emails
	var WetherspoonsTransactions []Transaction
	var trainlineTransactions []Transaction

	for i := range trainlineEmails {
		if !strings.Contains(trainlineEmails[i].Subject, "Your replacement booking confirmation") {
			thisResult := parseTrainline(trainlineEmails[i])
			trainlineTransactions = append(trainlineTransactions, thisResult)
		}
	}

	//parse all spoons emails emails
	for i := range wetherspoonEmails {
		thisResult := parseWetherspoon(wetherspoonEmails[i])
		WetherspoonsTransactions = append(WetherspoonsTransactions, thisResult)
	}

	//match up transactions
	matchTransactionsMonzo(monzoOutput, trainlineTransactions, 3, "Trainline")
	matchTransactionsMonzo(monzoOutput, WetherspoonsTransactions, 1.5, "JD Wetherspoon")

}

func parseWetherspoon(mail eazye.Email) Transaction {
	// call for HTML parsing
	items, recdet := parseHTMLWetherspoon(string(mail.HTML))
	var thisTransaction Transaction
	thisTransaction.details = recdet
	thisTransaction.item = items

	return thisTransaction
}

func parseTrainline(mail eazye.Email) Transaction {

	var transaction = parseTrainlineHTML(string(mail.HTML))
	transaction.details.OrderDateTime = mail.InternalDate
	return transaction

}
