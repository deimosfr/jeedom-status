package pkg

import (
    "fmt"
    "math/rand"
    "strconv"
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

func RandomNumberAsString() string {
    return strconv.Itoa(rand.Intn(8)+1)
}

func GetSampleJeedomGlobalStatus() map[string]string {
    return map[string]string{
        "alarm": "1",
        "door": RandomNumberAsString(),
        "humidity": RandomNumberAsString(),
        "light": RandomNumberAsString(),
        "luminosity": RandomNumberAsString(),
        "motion": RandomNumberAsString(),
        "outlet": RandomNumberAsString(),
        "power": RandomNumberAsString(),
        "security": RandomNumberAsString(),
        "shutter": RandomNumberAsString(),
        "temperature": RandomNumberAsString(),
        "windows": RandomNumberAsString(),
    }
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