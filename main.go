package main

import ( 
    "fmt"
    "github.com/jprobinson/eazye"
    "runtime"
)


var procs = runtime.NumCPU()

var transactions []Transaction

func main() {

    //fetch mail
    wetherspoonEmails := getMail("orders@jdwetherspoon.co.uk")

    //print to console how many reciepts we've got
    emailCount := len(wetherspoonEmails)
    fmt.Println(emailCount)

    // fetch monzo transactions
    monzoOutput := monzoFunc()

    // parse all spooms emails emails
    mailCount := 0
    for range wetherspoonEmails {
        parseWetherspoon(wetherspoonEmails[mailCount])
        mailCount++;
    }

    // match up transactions
    matchTransactionsMonzo(monzoOutput, transactions)

}

func parseWetherspoon(mail eazye.Email){
        // call for HTML parsing
        items, recdet := parseHTMLWetherspoon(string(mail.HTML))
        var thisTransaction Transaction
        thisTransaction.details = recdet
        thisTransaction.item = items

        transactions = append(transactions, thisTransaction)
        //fmt.Println(transactions)
}
