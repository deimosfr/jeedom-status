package pkg

import (
    "fmt"
    "strings"
)

func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}

func CheckArgContent(selected string, allowedArgs []string) bool {
    _, found := Find(allowedArgs, selected)
    if !found {
        fmt.Printf(
            "Value %s is not a valid bar type, allowed values are: %s\n",
            selected,
            strings.Join(allowedArgs, " "),
        )
        return false
    }
    return true
}