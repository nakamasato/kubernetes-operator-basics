package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v47/github"
	"github.com/nakamasato/kubernetes-operator-basics/05-golang-basics/02-module/mypackage"
)

func main() {
	name := mypackage.GetName()
	fmt.Printf("Hello, %s\n", name)

	client := github.NewClient(nil)
	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List(context.Background(), "willnorris", nil)
	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
	for i, org := range orgs {
		fmt.Println(i, *org.Login)
	}
}
