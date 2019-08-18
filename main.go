package main

import ( 
 //   "fmt"
    "github.com/jprobinson/eazye"
    "runtime"
    "golang.org/x/net/html"
    "strings"
    "strconv"
)


var procs = runtime.NumCPU()

var WetherspoonsTransactions []Transaction

var monzoOutput []MonzoTransaction

func main() {

    //wetherspoonEmails := getMail("FROM orders@jdwetherspoon.co.uk")
    trainlineEmails := getMail(`TEXT "Your booking confirmation. Transaction Id" FROM "auto-confirm@info.thetrainline.com"`)


    //print to console how many reciepts we've got
    //emailCount := len(wetherspoonEmails)
    //fmt.Println(emailCount)

    // fetch monzo transactions
    monzoOutput = monzoFunc()

    // parse all trainline emails

    var trainlineTransactions []Transaction
    
    for i := range trainlineEmails {
        if (!strings.Contains(trainlineEmails[i].Subject, "Your replacement booking confirmation")){
            trainlineTransactions = append(trainlineTransactions, parseTrainline(trainlineEmails[i]))
        }
    }
    //fmt.Println(trainlineTransactions)
    matchTransactionsMonzo(monzoOutput, trainlineTransactions, 3, "Trainline")

    // //parse all spoons emails emails
    // for i := range wetherspoonEmails {
    //     parseWetherspoon(wetherspoonEmails[i])
    // }

    // //match up transactions
    // matchTransactionsMonzo(monzoOutput, WetherspoonsTransactions, 1.5, "JD Wetherspoon")

}


func StripTrailing(start string) (string){
    nospace := strings.TrimSpace(start)

    return nospace
}

func stringToInt(money string) (int64){
    nospace := strings.Replace(money, " ", "", -1)
    val, err := strconv.Atoi(nospace)
    
    if(err != nil) {
        panic(err)
    }

    return int64(val)
}

func stringToPence(money string) (int64){
    nospace := strings.Replace(money, " ", "", -1)
    stringVal := strings.ReplaceAll(nospace, ".", "")
    stringVal = strings.ReplaceAll(stringVal, "£", "")
    currencyVal, err := strconv.Atoi(stringVal)
    
    if(err != nil) {
        panic(err)
    }

    return int64(currencyVal)
}

func parseWetherspoon(mail eazye.Email){
        // call for HTML parsing
        items, recdet := parseHTMLWetherspoon(string(mail.HTML))
        var thisTransaction Transaction
        thisTransaction.details = recdet
        thisTransaction.item = items

        WetherspoonsTransactions = append(WetherspoonsTransactions, thisTransaction)
        //fmt.Println(transactions)
}


func parseTrainline(mail eazye.Email) (Transaction) {
     
    var transaction = parseTrainlineHTML(string(mail.HTML))
    transaction.details.OrderDateTime = mail.InternalDate
    return transaction
    
}

func parseTrainlineHTML(htmlStr string) (Transaction) {
        //fmt.Println(htmlStr)
        bodyString := strings.NewReader(htmlStr)
        token := html.NewTokenizer(bodyString)
        tokenCount := 0
        content := []string{}
        tripDetails := []string{}
        relevantInfoPayment := false
        relevantInfoTrip := false
    
        for {
            tt := token.Next()
            t := token.Token()
            td := t.Data
            tokenCount++

            if(tt == html.TextToken){
                //fmt.Println(td)
            }

            if tt == html.ErrorToken {
                    break
            
            } else if(strings.Contains(td, "Payment information")){
                    relevantInfoTrip = false
                    relevantInfoPayment = true
 
            } else if (tt == html.TextToken && relevantInfoTrip){
                stringContent :=  strings.TrimSpace(td)
                if(len(stringContent) > 0){
                            tripDetails = append(tripDetails, stringContent)
                }
            } else if (tt == html.TextToken && relevantInfoPayment){
                    stringContent :=  strings.TrimSpace(td)
                    if(len(stringContent) > 0){
                            content = append(content, stringContent)
                    }
            }
        }

    
        var transaction Transaction
        var details RecieptDetails

        var items []RecieptItem
        recieptItems := false;


        // For the actual content
        k := 0;
        for {      
            if(k < len(content)){

            if(content[k] == "Transaction Date:") {
                //transactionDate = content[k+1]
                //fmt.Println(transactionDate)
            } else if(content[k] == "Fare details") {
                recieptItems = true;
    
            } else if (strings.Contains(content[k], "Other costs*")){ // end of items
                recieptItems = false;
            
            } else if(recieptItems){
               
                if(strings.HasPrefix(content[k], "Booking")){
                    k++;
                }
                if(content[k] == "Outbound"){
                    k++;
                }
                if(content[k] == "Return"){
                    k++;
                }

                var thisItem RecieptItem
                thisItem.Quantity = strings.Split(content[k], "x")[0]
                thisItem.Description = StripTrailing(strings.Split(content[k], "x")[1])
                thisItem.Price = strings.Split(content[k + 1], "x")[1]
                k++;
                
                items = append(items, thisItem)
                
            }else if(content[k] == "Booking Fee:") { 
                var thisItem RecieptItem
                thisItem.Quantity = "1"
                thisItem.Description = "Booking Fee"
                
                if(content[k + 1] == "No booking fee"){
                    thisItem.Price = "0"
                }else{
                    thisItem.Price = content[k + 1]
                }
    
                items = append(items, thisItem)
            } else if(content[k] == "Total amount:"){
                details.OrderWithVat = stringToPence(content[k + 1])
    
            } else if(strings.HasPrefix(content[k], "VAT number")){
                details.VatNumber = strings.Split(content[k], "VAT number ")[1]
                details.VatTotal = 0
    
                break;
            }
        }

        k++;
    
        }
        transaction.details = details
        transaction.item = items
        return transaction;
}
