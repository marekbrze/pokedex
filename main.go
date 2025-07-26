package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	if text == "" {
		return []string{}
	}
	reg := regexp.MustCompile("[^a-z ]+")
	wordsList := strings.Fields(strings.ToLower(reg.ReplaceAllString(text, "")))
	return wordsList
}
