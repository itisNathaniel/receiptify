package main

import ( 
    "fmt"
    "github.com/jprobinson/eazye"
    "runtime"
)


var procs = runtime.NumCPU()

func main() {

    // fetch mail
    mail := getMail()

    // print to console how many reciepts we've got
    emailCount := len(mail)
    fmt.Println(emailCount)

    // parse all emails
    mailCount := 0
    for range mail {
        parseMessages(mail[mailCount])
        mailCount++;
    }

    // look up reciepts in Monzo
}

func parseMessages(mail eazye.Email){
        // call for HTML parsing
        items, recdet := parseHTML(string(mail.HTML))
        fmt.Println(items, recdet)
}