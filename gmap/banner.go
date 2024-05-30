package main

import (
	"fmt"

	"moul.io/banner"
)

func printBanner() {
	// Color codes
	const lines string = "--------------------------------"
	purple := "\033[35m"
	reset := "\033[0m"
	red := "\033[31m"
	blue := "\033[34m"

	// Banner string
	bannerStr := banner.Inline("gomap")

	fullBanner := fmt.Sprintf("%s%s%s\n%s%s%s\n%s%s%s\n%s%s%s\n\n%s%s%s\n\n%s%s%s\n%s%s%s\n%s%s%s\n%s%s%s\n",
		purple, lines, reset,
		purple, lines, reset,
		purple, lines, reset,
		purple, bannerStr, reset,
		red, "Made by Majir96", reset,
		blue, "Under GPL License", reset,
		purple, lines, reset,
		purple, lines, reset,
		purple, lines, reset,
	)

	fmt.Println(fullBanner)
}
