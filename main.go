package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/jprobinson/eazye"
)

var procs = runtime.NumCPU()

var WetherspoonsTransactions []Transaction

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

	var trainlineTransactions []Transaction

	for i := range trainlineEmails {
		if !strings.Contains(trainlineEmails[i].Subject, "Your replacement booking confirmation") {
			thisResult := parseTrainline(trainlineEmails[i])
			trainlineTransactions = append(trainlineTransactions, thisResult)
		}
	}
	matchTransactionsMonzo(monzoOutput, trainlineTransactions, 3, "Trainline")

	//parse all spoons emails emails
	for i := range wetherspoonEmails {
		parseWetherspoon(wetherspoonEmails[i])
	}

	//match up transactions
	matchTransactionsMonzo(monzoOutput, WetherspoonsTransactions, 1.5, "JD Wetherspoon")

}

func parseWetherspoon(mail eazye.Email) {
	// call for HTML parsing
	items, recdet := parseHTMLWetherspoon(string(mail.HTML))
	var thisTransaction Transaction
	thisTransaction.details = recdet
	thisTransaction.item = items

	WetherspoonsTransactions = append(WetherspoonsTransactions, thisTransaction)
}

func parseTrainline(mail eazye.Email) Transaction {

	var transaction = parseTrainlineHTML(string(mail.HTML))
	transaction.details.OrderDateTime = mail.InternalDate
	return transaction

}
