package target

import (
    "fmt"
)

type Target struct {
    name        string
    address     string
}

func MakeNewTarget(name string, address string) Target {
    toReturn := Target{}
    toReturn.name = name
    toReturn.address = address
    return toReturn
}

func (t Target) Print () {
    fmt.Printf("Name: %s Address: %s\n", t.name, t.address)
}

func Print () {
    fmt.Println("useless")
}
