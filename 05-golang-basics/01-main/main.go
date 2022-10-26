package main

import "fmt"

func main() {
	fmt.Println("Hello, World")
	greetText, err := greet("English", "Naka")
	if err == nil {
		fmt.Println(greetText)
	}
}

func greet(language, name string) (string, error) {
	if language == "Spanish" {
		return fmt.Sprintf("Ola, %s", name), nil
	}
	return fmt.Sprintf("Hello, %s", name), nil
}
