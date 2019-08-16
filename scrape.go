package main

import ( 
    "golang.org/x/net/html"
    "strings"
    "unicode"

)


func parseHTML(stringHTML string) ([]RecieptItem, RecieptDetails){
    bodyString := strings.NewReader(stringHTML)
    token := html.NewTokenizer(bodyString)
    content := []string{}

    recieptItems := []RecieptItem{}
    RecieptDetail := RecieptDetails{}
    tokenCount := 0

    // Extract from HTML
    for {
        tt := token.Next()
        t := token.Token()
        td := t.Data
        tokenCount++
        // line 45 is where the reciept content starts
        if tt == html.ErrorToken {
                break
        } else if (tt == html.TextToken && tokenCount > 45){
                stringContent :=  strings.TrimSpace(td)
                if(len(stringContent) > 0){
                        //fmt.Println(stringContent)
                        content = append(content, stringContent)
                }
        }


    }

    RecieptDetail.Name = content[0]
    RecieptDetail.Address = content[1] + ", " + content[2]
    // some have more addresses - postcodes are 6-8 chars + space so need to offset
    var addressOffset int = 0;
    if(len(content[3]) < 9){
        addressOffset = 0;
    } else {
        addressOffset = 1;
    }
    RecieptDetail.Postcode = content[3 + addressOffset]
    RecieptDetail.Phone = (content[4 + addressOffset])[11:]
    RecieptDetail.OrderID = content[6 + addressOffset]
    RecieptDetail.OrderDate = (content[7 + addressOffset])[12:]
    RecieptDetail.OrderTime = (content[8 + addressOffset])[12:]
    RecieptDetail.Table = content[10 + addressOffset]

    var endOfItemsIndex int
    thisLineIndex := 14 + addressOffset

    for {
        if(content[thisLineIndex] == "Payment Type"){
                endOfItemsIndex = thisLineIndex; 
                break;
        } else {
                Quantity := content[thisLineIndex]
                Name := content[thisLineIndex + 1]
                AdditionalInfoOrPrice := content[thisLineIndex + 2]
                PriceOrNextItem := content[thisLineIndex + 3]
                
                var thisItem RecieptItem

                if(!isInt(Quantity)){
                // handle blank quantity issue
                        thisItem.Price = thisItem.Description
                        thisItem.Description = thisItem.Quantity
                        thisItem.Quantity = "0"
                        thisLineIndex = thisLineIndex + 2;
                } else if(strings.Contains(AdditionalInfoOrPrice, "£")) {
                        // handle prices
                        thisItem.Quantity = Quantity
                        thisItem.Description = Name
                        thisItem.Price = AdditionalInfoOrPrice;

                        thisLineIndex = thisLineIndex + 3;
                } else {
                        thisItem.Quantity = Quantity
                        thisItem.Description = Name + " " + AdditionalInfoOrPrice
                        thisItem.Price = PriceOrNextItem;

                        thisLineIndex = thisLineIndex + 4;
                }

                recieptItems = append(recieptItems, thisItem)

        }
    }

    RecieptDetail.PayMethod = content[endOfItemsIndex + 1]
    RecieptDetail.OrderTotal = content[endOfItemsIndex + 3]
    RecieptDetail.VatTotal = content[endOfItemsIndex + 5]
    RecieptDetail.VatNumber = content[endOfItemsIndex + 9]

// Calculates total spent
//     cost := strings.ReplaceAll(RecieptDetail.OrderTotal, "£", "")
//     vat := strings.ReplaceAll(RecieptDetail.VatTotal, "£", "")
//     totalvat, err := strconv.ParseFloat(vat, 64)
//     totalcost, err := strconv.ParseFloat(cost, 64)
//     fmt.Println(err);
//     totalSpent = totalSpent + totalcost + totalvat


    // Print to check the slice's content
    //fmt.Println(content)
    return recieptItems,RecieptDetail
}

func isInt(s string) bool {
        for _, c := range s {
            if !unicode.IsDigit(c) {
                return false
            }
        }
        return true
    }

// Concurrent stuff to return to 
//
//     // concurrency time 
//     var parsers sync.WaitGroup
// 	for i := 0; i < emailCount; i++ {
// 		parsers.Add(1)
// 		// multi goroutines so we can utilize the CPU while waiting for URLs
// 		go parseMessages(mail[i])
//     }

    
//     fmt.Println("done woo")
//     //fmt.Println(mail,err)
// }