package main

import ( 
    "fmt"
    "github.com/jprobinson/eazye"
    "os"
    "strconv"
    "runtime"
)

var procs = runtime.NumCPU()

func main() {
    fmt.Println("hello world")
    // get from env 
    mailserver := os.Getenv("mailserver");
    mailserverssl, err := strconv.ParseBool(os.Getenv("mailserverssl"));
    emailaddress := os.Getenv("emailaddress");
    password := os.Getenv("password");
    folder := os.Getenv("folder");
    
	mailbox := eazye.MailboxInfo{mailserver, mailserverssl, emailaddress, password, folder, true}
    
    mail, err := eazye.GetCommand(mailbox, "FROM orders@jdwetherspoon.co.uk", true, false)

    emailCount := len(mail)

    // Just double check 
    fmt.Println(emailCount, err)

    parseMessages(mail[0])

}

func parseMessages(mail eazye.Email){
    fmt.Println("Here I am", string(mail.HTML))
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
