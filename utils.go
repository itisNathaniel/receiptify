package main

import ( 
	   "strings"
	   "strconv"
   )

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