package main

import "strings"

func removePrefixAndCamelCase(input, clear string) string {
	if strings.HasPrefix(input, "wl_") {
		input = input[3:]
	}
	if strings.HasPrefix(input, clear+"_") {
		input = input[len(clear)+1:]
	}

	input = strings.Replace(input, "_"+clear+"_", "_", -1)

	// Split the string into words based on underscore
	words := strings.Split(input, "_")

	// Capitalize the first letter of each word (including the first word)
	for i := 0; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}

	// Join the words to form the camel case string
	result := strings.Join(words, "")

	return result
}

func before_wl(input string) string {
	var output = strings.Split(input, "_")[0]
	if output == "wayland" {
		return "wl"
	}
	return output
}

func sanitizeSingleLineComment(input string) string {
	input = strings.Replace(input, "\n", "", -1)
	input = strings.Replace(input, "\r", "", -1)
	input = strings.Replace(input, "\t", " ", -1)
	input = strings.Replace(input, "  ", " ", -1)
	return input
}
