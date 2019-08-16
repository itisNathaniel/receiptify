package main

import ( 
    "fmt"
    "github.com/jprobinson/eazye"
    "runtime"
    "strings"
    "strconv"
    "os"
    "encoding/json"
    "github.com/dghubble/sling"
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
        //fmt.Println(transactions)
}


type Reciept_Tax struct {
	Description        string       `json:"description"`
	Amount			   int64		`json:"amount"`
	Currency		   string		`json:"currency"`
	TaxNumber		   string		`json:"tax_number"`
}

type Reciept_Merchant struct {
	Name        	   string       `json:"name"`
	Online			   bool			`json:"online"`
	StorePhone		   string		`json:"phone"`
	StoreName		   string		`json:"store_name"`
	StorePostcode	   string		`json:"store_postcode"`
}

type Reciept struct {
	TransactionID	   string			`json:"transaction_id"`
	ExternalID		   string			`json:"external_id"`
	Total			   int64			`json:"total"`
	Currency		   string			`json:"currency"`
	Items			   []MonzoRecieptItem  	`json:"items"`
	Merchant		   Reciept_Merchant `json:"merchant"`
}

type MonzoRecieptItem struct {
	Description		    string			`json:"description"`
	Quantity			int64			`json:"quantity"`
	Unit				string			`json:"unit"`
	Amount				int64			`json:"amount"`
	Currency			string			`json:"currency"`
}

type RawTransaction struct {
	Transaction
	Merchant json.RawMessage `json:"merchant"`
}

func matchTransactionsMonzo(monzoTransact []MonzoTransaction, transactions []Transaction){

    for i := range monzoTransact {
        monzoDate := monzoTransact[i].Date
        monzoTransaction := monzoTransact[i]
        for j := range transactions {
            transactionSearch := transactions[j]

            if(monzoTransaction.CurrencyAmount == transactionSearch.details.OrderWithVat){
                fmt.Println("GOT HERE")
                transactionDate := transactionSearch.details.OrderDateTime
                // check the date 
                diff := transactionDate.Sub(monzoDate)

                // generate an ID (should be replicated between runs to not create multiple reciepts for one event)
                identifier := strconv.FormatInt(monzoDate.Unix(), 10) + "-" + transactionSearch.details.Name
                identifier = strings.Replace(identifier, " ", "-", -1)


                // sometimes they have a backlog in sending them out, it seems
                if(diff.Hours() < 1.5){
                    fmt.Println("WOOO   " + " !! " + monzoTransaction.TransactionId + "!!!!! " + identifier)
                    
                    var thisTransaction = monzoTransaction
                    var items []MonzoRecieptItem

                    fmt.Println("!!!!! SEARCH RESULT")
                    fmt.Println(transactionSearch)

                    for k := range transactionSearch.item {
                        item := transactionSearch.item[k]
                        var modelItem MonzoRecieptItem

                        price, err := strconv.ParseInt(item.Price, 10, 64)
                        if(err != nil) {
                            fmt.Println(err)
                        }
                        fmt.Println(item.Price + " ---- ")
                        fmt.Println(price)

                        modelItem.Description = item.Description
                        modelItem.Unit = item.Quantity
                        modelItem.Amount = price
                        modelItem.Currency = "GBP"
                        
                        items = append(items, modelItem)
                    }
    
                    // PASS and make the monzo reciept

                    // Merchant
                    var merchant Reciept_Merchant
                    merchant.Name = "JD Wetherspoon"
                    merchant.Online = false
                    merchant.StoreName = transactionSearch.details.Name
                    merchant.StorePhone = transactionSearch.details.Phone
                    merchant.StorePostcode = transactionSearch.details.Postcode

                    // Tax
                    var tax Reciept_Tax
                    tax.Description = "VAT"
                    tax.Currency = "GBP"
                    tax.Amount = transactionSearch.details.VatTotal
                    tax.TaxNumber = transactionSearch.details.VatNumber

                    reciept := Reciept{
                        TransactionID: thisTransaction.TransactionId,
                        ExternalID: identifier,
                        Total: thisTransaction.CurrencyAmount,
                        Currency: "GBP",
                        Items: items, 
                        Merchant: merchant,
                        //Tax: tax,
                        
                    }

                    res, err := AddReciept(reciept)

                    fmt.Println(res,err)

                }

            }
        
        }
    }
}

func AddReciept(reciept Reciept) (string, error) {

    body := reciept;

    monzoBase := sling.New().Base("https://api.monzo.com/").Set("Authorization", "Bearer " + os.Getenv("accesstoken"))
    path := fmt.Sprintf("transaction-receipts")
    
    req, err := monzoBase.New().Put(path).BodyJSON(body).Request()

    fmt.Println(body, req, err)

	return "done", err
}
