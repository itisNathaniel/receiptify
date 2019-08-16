package main

import ( 
    "fmt"
    //"github.com/jprobinson/eazye"
    //"os"
    //"strconv"
    "golang.org/x/net/html"
    "runtime"
    "strings"
)

type RecieptDetails struct {
    OrderID         string
    OrderDate       string
    OrderTime       string
    VatNumber       string
    BarName         string
}

type RecieptItem struct {
    Description     string
    Quantity        string
    TableNumber     int
    VatNumber       string
}


var procs = runtime.NumCPU()

func main() {
    fmt.Println("hello world")
    // get from env 
    // mailserver := os.Getenv("mailserver");
    // mailserverssl, err := strconv.ParseBool(os.Getenv("mailserverssl"));
    // emailaddress := os.Getenv("emailaddress");
    // password := os.Getenv("password");
    // folder := os.Getenv("folder");
    
	//mailbox := eazye.MailboxInfo{mailserver, mailserverssl, emailaddress, password, folder, true}
    
    //mail, err := eazye.GetCommand(mailbox, "FROM orders@jdwetherspoon.co.uk", true, false)

    //emailCount := len(mail)

    // Just double check 
    //fmt.Println(emailCount, err)

    //parseMessages(mail[0])

    items, recdet := parseHTML(horrorHTML)

    fmt.Println(items, recdet)

}

// func parseMessages(mail eazye.Email){
//     fmt.Println("Here I am", string(mail.HTML))
// }

func parseHTML(stringHTML string) ([]RecieptItem, RecieptDetails){
    bodyString := strings.NewReader(stringHTML)
    token := html.NewTokenizer(bodyString)
    content := []string{}

    recieptItems := []RecieptItem{}
    RecieptDetail := RecieptDetails{}
    tokenCount := 0

    for {
        tt := token.Next()
        t := token.Token()
        td := t.Data
        tokenCount++
        // line 45 is where the reciept content starts
        if (tt == html.TextToken && tokenCount > 45){
                stringContent :=  strings.TrimSpace(td)
                if(len(stringContent) > 0){
                        fmt.Println(stringContent)
                        content = append(content, td)
                }
        }
    }
    // Print to check the slice's content
    //fmt.Println(content)
    return recieptItems,RecieptDetail
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


var horrorHTML = `<!DOCTYPE html PUBLIC '-//W3C//DTD XHTML 1.0 Transitional//EN' 'http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd'>
<html xmlns='http://www.w3.org/1999/xhtml'>
<head>
<meta http-equiv='Content-Type' content='text/html; charset=utf-8' />
<title>&nbsp;</title>
<style type='text/css'>
.centerit {
        text-align: center;
}
body p {
        font-family: Tahoma, Geneva, sans-serif;
}
body p {
        font-size: 10px;
}
body p {
        font-size: 12px;
}
.tableclass {
        font-family: Tahoma, Geneva, sans-serif;
        font-size: 12px;
}
.headers {
        font-size: 16px;
}
</style>
</head>

<body> <table width='550' border='0' align='center' cellpadding='0' cellspacing='0'>
                <tr>
                        <td>
                                <img src='' />
                        </td>
                </tr>
                <tr>
                        <td><table width='550' border='0' cellspacing='0' cellpadding='10'>
                <tr>
                        <td><p class='headers'><strong>Thank you for your order</strong>...</p>
                        <table width='100%' border='0' cellspacing='0' cellpadding='0'>
                                <tr class='tableclass'>
                                        <td width='50%' align='left' valign='top'>
                                                <p>
                                                        The Poulton Elk<br />
                                                        22 Hardhorn Road<br />
                                                        <br />
                                                        Poulton-le-Fylde<br />
                                                        FY6 7SR<br />
                                                        <br />
                                                        Telephone: 01253 895265
                                                </p>
                                        </td>
                                        <td width='50%' align='left' valign='top'>
                                                <p>
                                                        Order Id: <strong>2123338</strong><br />
                                                        <br />
                                                        Order Date: Thursday, June 29, 2017<br />
                                                        Order Time: 21:19<br />
                                                        <br />
                                                        Table Number: <strong>21</strong>
                                                </p>
                                        </td>
                                </tr>
                        </table>
                        <br />
                        <table width='100%' border='1' cellspacing='0' cellpadding='2'>
                                <tr>
                                        <td width='12%'><span class='tableclass'>Qty</span></td>
                                        <td width='72%'><span class='tableclass'>Description</span></td>
                                        <td width='16%'><span class='tableclass'>Amount</span></td>
                                </tr>
                                <tr>
                                        <td>&nbsp;</td>
                                        <td class='tableclass'>&nbsp;</td>
                                        <td class='tableclass'>&nbsp;</td>
                                </tr>

                                                <tr>
                                                        <td class='tableclass'>1</td>
                                                        <td class='tableclass'><strong>Ruddles Best</strong>
                                                        (Pint) 
                                                        </td>
                                                        <td align='right' class='tableclass'>£1.99</td>
                                                </tr>


                                <tr>
                                        <td>&nbsp;</td>
                                        <td class='tableclass'>&nbsp;</td>
                                        <td class='tableclass'>&nbsp;</td>
                                </tr>


                                <tr>
                                        <td align='right' class='tableclass' colspan='2'><span class='tableclass'><strong>Payment Type</strong></span></td>
                                        <td align='right' class='tableclass'><strong>Apple Pay</strong></td>
                                </tr>
                                <tr>
                                        <td align='right' class='tableclass' colspan='2'><span class='tableclass'><strong>Order Total</strong></span></td>
                                        <td align='right' class='tableclass'><strong>£1.99</strong></td>
                                </tr>
                                <tr>
                                        <td align='right' class='tableclass' colspan='2'><span class='tableclass'><strong>VAT (20%)</strong></span></td>
                                        <td align='right' class='tableclass'><strong>£0.33</strong></td>
                                </tr>
                        </table>
                        <p>Thank you for your order</p>
                        <p>If you require a breakdown of the VAT, please speak to a member of the bar staff.</p>
                        <p>VAT Number - <strong><cfloop query='getVat'>396 331 433</cfloop></strong></p>
                        </td>
                </tr>
        </table>
</td>
</tr>
<tr>
        <td>
                <img src='' />
        </td>
</tr>
</table> </body>
</html>`