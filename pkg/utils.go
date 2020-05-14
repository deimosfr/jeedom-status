package pkg

import (
    "math/rand"
    "strconv"
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