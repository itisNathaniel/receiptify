package main

import ( 
    "fmt"
    "github.com/jprobinson/eazye"
    "runtime"
)


var procs = runtime.NumCPU()

//var totalSpent float64 = 0;

func main() {
    //get from env 

    mail := getMail()
    emailCount := len(mail)

    //Just double check 
    fmt.Println(emailCount)

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