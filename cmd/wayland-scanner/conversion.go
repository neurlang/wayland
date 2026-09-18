package main

import "strings"

func removePrefixAndCamelCase(input, clear string) string {
	input = strings.TrimPrefix(input, "wl_")
	input = strings.TrimPrefix(input, clear+"_")

	input = strings.ReplaceAll(input, "_"+clear+"_", "_")

	// Split the string into words based on underscore
	words := strings.Split(input, "_")

	// Capitalize the first letter of each word (including the first word)
	// Filter out empty strings from consecutive underscores
	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(word[:1]))
			result.WriteString(word[1:])
		}
	}

	return result.String()
}

func beforeWl(input string) string {
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
		input = strings.ReplaceAll(input, "  ", " ")
	}

	return input
}
