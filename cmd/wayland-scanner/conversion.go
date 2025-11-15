package main

import "strings"

func removePrefixAndCamelCase(input, clear string) string {
	input = strings.TrimPrefix(input, "wl_")
	input = strings.TrimPrefix(input, clear+"_")

	input = strings.Replace(input, "_"+clear+"_", "_", -1)

	// Split the string into words based on underscore
	words := strings.Split(input, "_")

	// Capitalize the first letter of each word (including the first word)
	// Filter out empty strings from consecutive underscores
	var result strings.Builder
	for i := 0; i < len(words); i++ {
		if len(words[i]) > 0 {
			result.WriteString(strings.ToUpper(words[i][:1]))
			result.WriteString(words[i][1:])
		}
	}

	return result.String()
}

func before_wl(input string) string {
	var output = strings.Split(input, "_")[0]
	if output == "wayland" {
		return "wl"
	}
	return output
}

func sanitizeSingleLineComment(input string) string {
	replacer := strings.NewReplacer(
		"\n", "",
		"\r", "",
		"\t", " ",
	)
	input = replacer.Replace(input)
	
	// Normalize multiple spaces to single space
	for strings.Contains(input, "  ") {
		input = strings.Replace(input, "  ", " ", -1)
	}
	
	return input
}
