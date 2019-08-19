package main

import (
	"runtime"
	"strings"

	"github.com/jprobinson/eazye"
)

var procs = runtime.NumCPU()

var monzoOutput []MonzoTransaction

func main() {

	// fetch emails
	wetherspoonEmails := getMail("FROM orders@jdwetherspoon.co.uk")
	trainlineEmails := getMail(`TEXT "Your booking confirmation. Transaction Id" FROM "auto-confirm@info.thetrainline.com"`)

	// fetch monzo transactions
	monzoOutput = getMonzoTransactions()

	// structs for transactions
	var WetherspoonsTransactions []Transaction
	var trainlineTransactions []Transaction

	// parse all trainline emails
	for i := range trainlineEmails {
		if !strings.Contains(trainlineEmails[i].Subject, "Your replacement booking confirmation") {
			thisResult := parseTrainline(trainlineEmails[i])
			trainlineTransactions = append(trainlineTransactions, thisResult)
		}
	}

	// parse all spoons emails emails
	for i := range wetherspoonEmails {
		thisResult := parseWetherspoon(wetherspoonEmails[i])
		WetherspoonsTransactions = append(WetherspoonsTransactions, thisResult)
	}

	// match up transactions
	matchTransactionsMonzo(monzoOutput, trainlineTransactions, 3, "Trainline")
	matchTransactionsMonzo(monzoOutput, WetherspoonsTransactions, 1.5, "JD Wetherspoon")

}

// Build up transaction from HTML and metadata
func parseWetherspoon(mail eazye.Email) Transaction {
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
