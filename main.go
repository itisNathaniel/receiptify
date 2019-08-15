package main

import ( 
    "fmt"
    //"net/mail"
    "time"
    "github.com/jprobinson/eazye"
    "os"
    "strconv"
)


func main() {
    fmt.Println("hello world")
    // get from env 
    mailserver := os.Getenv("mailserver");
    mailserverssl, err := strconv.ParseBool(os.Getenv("mailserverssl"));
    emailaddress := os.Getenv("emailaddress");
    password := os.Getenv("password");
    folder := os.Getenv("folder");
    
	mailbox := eazye.MailboxInfo{mailserver, mailserverssl, emailaddress, password, folder, true}
    
    t := time.Now()
    rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

    fmt.Println(rounded)
    
    mail, err := eazye.GetUnread(mailbox, false, false)
    fmt.Println("done woo")
    fmt.Println(mail,err)
}