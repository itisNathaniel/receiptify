package main

import ( 
    "time"
    "github.com/tjvr/go-monzo"
    "os"
    "fmt"
)

func monzoFunc() ([]MonzoTransaction) {
    cl := monzo.Client{
        BaseURL: "https://api.monzo.com",
        AccessToken: os.Getenv("accesstoken"),
    }

    transactions, err := cl.Transactions("acc_00009QrR2ROd8BTPsp0vrN", true) // don't expandMerchant
    if err != nil {
        fmt.Println(err)
    }
    i := 0;

    var transactionList []MonzoTransaction

    layout := "2006-01-02T15:04:05.000Z"
    for {
        var transaction MonzoTransaction
        if(i < len(transactions)){
            if(transactions[i].Merchant.Name == "JD Wetherspoon" && transactions[i].Merchant.IsOnline == true){
                time, err := time.Parse(layout, transactions[i].Created)
                transaction.Date = time
                transaction.CurrencyAmount = -transactions[i].Amount
                transaction.TransactionId = transactions[i].ID

                transactionList = append(transactionList, transaction)
                if err != nil {
                    fmt.Println(err)
                }
            }
            } else {
                break;
            }
            i++;
        }
    
        return transactionList
    }