package main

import (
	"runtime"
	"strings"
	"sync"

	"github.com/jprobinson/eazye"
)

var procs = runtime.NumCPU()

func main() {

	// Initial Collection of Data
	var collect sync.WaitGroup
	var wetherspoonEmails []eazye.Email
	var trainlineEmails []eazye.Email
	var monzoTransactions []MonzoTransaction

	collect.Add(1)
	go func() {
		wetherspoonEmails = getMail("FROM orders@jdwetherspoon.co.uk")
		collect.Done()
	}()
	collect.Add(1)
	go func() {
		trainlineEmails = getMail(`TEXT "Your booking confirmation. Transaction Id" FROM "auto-confirm@info.thetrainline.com"`)
		collect.Done()
	}()
	collect.Add(1)
	go func() {
		monzoTransactions = getMonzoTransactions()
		collect.Done()
	}()

	collect.Wait()

	// Parse the data
	var parse sync.WaitGroup
	var WetherspoonsTransactions []Transaction
	var trainlineTransactions []Transaction

	parse.Add(len(trainlineEmails))
	go func() {
		for i := range trainlineEmails {
			if !strings.Contains(trainlineEmails[i].Subject, "Your replacement booking confirmation") {
				thisResult := parseTrainline(trainlineEmails[i])
				trainlineTransactions = append(trainlineTransactions, thisResult)
			}
			parse.Done()
		}
	}()

	parse.Add(len(wetherspoonEmails))
	go func() {
		for i := range wetherspoonEmails {
			thisResult := parseWetherspoon(wetherspoonEmails[i])
			WetherspoonsTransactions = append(WetherspoonsTransactions, thisResult)
			parse.Done()
		}
	}()

	parse.Wait()

	// Matching up the data
	var match sync.WaitGroup

	match.Add(1)
	go func() {
		matchTransactionsMonzo(monzoTransactions, trainlineTransactions, 3, "Trainline")
		match.Done()
	}()

	match.Add(1)
	go func() {
		matchTransactionsMonzo(monzoTransactions, WetherspoonsTransactions, 1.5, "JD Wetherspoon")
		match.Done()
	}()

	match.Wait()

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
