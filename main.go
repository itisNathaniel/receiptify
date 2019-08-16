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
    mail := getMail()

    //print to console how many reciepts we've got
    emailCount := len(mail)
    fmt.Println(emailCount)

    // fetch monzo transactions
    monzoOutput := monzoFunc()
    //fmt.Println(monzoOutput)

    // parse all emails
    mailCount := 0
    for range mail {
        parseMessages(mail[mailCount])
        mailCount++;
    }

    // match up transactions
    matchTransactionsMonzo(monzoOutput, transactions)

}

func parseMessages(mail eazye.Email){
        // call for HTML parsing
        items, recdet := parseHTML(string(mail.HTML))
        var thisTransaction Transaction
        thisTransaction.details = recdet
        thisTransaction.item = items

        transactions = append(transactions, thisTransaction)
        fmt.Println(transactions)
}


func matchTransactionsMonzo(monzoTransact []MonzoTransaction, transactions []Transaction){

    for i := range monzoTransact {
        monzoDate := monzoTransact[i].Date
        for j := range transactions {
            if(monzoTransact[i].CurrencyAmount == transactions[j].details.OrderWithVat){
                transactionDate := transactions[j].details.OrderDateTime
                // check the date 
                diff := transactionDate.Sub(monzoDate)

                // sometimes they have a backlog in sending them out, it seems
                if(diff.Hours() < 1.5){
                    fmt.Println("WOOO")

                    // PASS and make the monzo reciept

                    break;
                }

            }
        
        }
    }



}