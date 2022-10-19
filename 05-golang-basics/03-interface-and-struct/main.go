package main

import "fmt"

type Object interface {
	GetName() string
}

type Parent struct {
	Id   int
	Name string
}

func (p Parent) GetName() string {
	return p.Name
}

func (p *Parent) SetName(name string) {
	p.Name = name
}

type Child struct {
	Parent
	OtherField string
}

func main() {
    p := Parent{Id: 1, Name: "John"}
    printObjectName(&p)
    ch := Child{Parent: p}
    printObjectName(&ch)

    ch = Child{Parent{Id: 1, Name: "child"}, "Other"}
    printObjectName(&ch)
}

func printObjectName(obj Object) {
	fmt.Println(obj.GetName())
}
