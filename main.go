package main

import ( 
    "fmt"
    "github.com/jprobinson/eazye"
    "runtime"
    "strings"
    "strconv"
    "bytes"
    "encoding/json"
    "net/http"
    "os"
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
    RecieptTaxes       []Reciept_Tax    `json:"taxes"`
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
                transactionDate := transactionSearch.details.OrderDateTime
                // check the date 
                diff := transactionDate.Sub(monzoDate)

                // generate an ID (should be replicated between runs to not create multiple reciepts for one event)
                identifier := strconv.FormatInt(monzoDate.Unix(), 10) + "-" + transactionSearch.details.Name
                identifier = strings.Replace(identifier, " ", "-", -1)


                // sometimes they have a backlog in sending them out, it seems
                if(diff.Hours() < 1.5 && diff.Hours() > 0){
                    
                    var thisTransaction = monzoTransaction
                    var items []MonzoRecieptItem

                    fmt.Println("!!!!! SEARCH RESULT", diff.Hours())
                    fmt.Println(transactionSearch)

                    for k := range transactionSearch.item {
                        item := transactionSearch.item[k]
                        var modelItem MonzoRecieptItem

                        price, err := strconv.ParseInt(item.Price, 10, 64)
                        if(err != nil) {
                            fmt.Println(err)
                        }

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
                    var taxes []Reciept_Tax

                    var tax Reciept_Tax
                    tax.Description = "VAT"
                    tax.Currency = "GBP"
                    tax.Amount = transactionSearch.details.VatTotal
                    tax.TaxNumber = transactionSearch.details.VatNumber
                    taxes= append(taxes, tax)

                    reciept := Reciept{
                        TransactionID: thisTransaction.TransactionId,
                        ExternalID: identifier,
                        Total: thisTransaction.CurrencyAmount,
                        Currency: "GBP",
                        Items: items, 
                        Merchant: merchant,
                        RecieptTaxes: taxes,
                        
                    }

                    res, err := AddReciept(reciept)

                    fmt.Println(res,err)

                }

            }
        
        }
    }
}

func AddReciept(reciept Reciept) (string, error) {

    json, err := json.Marshal(reciept)
    if err != nil {
        panic(err)
    }

    fmt.Println(json)

    client := &http.Client{}
    url := "https://api.monzo.com/transaction-receipts"
    req, err := http.NewRequest("PUT", url,  bytes.NewBuffer(json))
    req.Header.Set("Content-Type", "application/json; charset=utf-8")
    req.Header.Set("Authorization", "Bearer " + os.Getenv("accesstoken"))

    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    fmt.Println(resp.StatusCode)

	return "done", err
}
