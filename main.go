package main

import ( 
    "fmt"
    "github.com/jprobinson/eazye"
    "os"
    "strconv"
    "runtime"
)


var procs = runtime.NumCPU()

//var totalSpent float64 = 0;

func main() {
    //get from env 
    mailserver := os.Getenv("mailserver");
    mailserverssl, err := strconv.ParseBool(os.Getenv("mailserverssl"));
    emailaddress := os.Getenv("emailaddress");
    password := os.Getenv("password");
    folder := os.Getenv("folder");
    
	mailbox := eazye.MailboxInfo{mailserver, mailserverssl, emailaddress, password, folder, true}
    
    mail, err := eazye.GetCommand(mailbox, "FROM orders@jdwetherspoon.co.uk", true, false)

    emailCount := len(mail)

    //Just double check 
    fmt.Println(emailCount, err)

    mailCount := 0
    for range mail {
        parseMessages(mail[mailCount])
        mailCount++;
    }

}

func parseMessages(mail eazye.Email){
        items, recdet := parseHTML(string(mail.HTML))
        fmt.Println(items, recdet)
}
