package main

import ( 
	"github.com/jprobinson/eazye"
	"os"
	"strconv"
	"fmt"
)

func getMail(email string) ([]eazye.Email){
    mailserver := os.Getenv("mailserver");
    mailserverssl, err := strconv.ParseBool(os.Getenv("mailserverssl"));
    emailaddress := os.Getenv("emailaddress");
    password := os.Getenv("password");
    folder := os.Getenv("folder");
    
	mailbox := eazye.MailboxInfo{mailserver, mailserverssl, emailaddress, password, folder, true}
    
    mail, err := eazye.GetCommand(mailbox, "FROM " + email, true, false)

    if(err != nil) {
        fmt.Println(err)
    }
    return mail
}