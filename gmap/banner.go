package main

import (
	"fmt"
	"gmap/utils"

	"moul.io/banner"
)

func printBanner() {
	// Banner string
	bannerStr := banner.Inline("gomap")

	fullBanner := fmt.Sprintf("%s%s%s\n%s%s%s\n%s%s%s\n%s%s%s\n\n%s%s%s\n\n%s%s%s\n%s%s%s\n%s%s%s\n%s%s%s\n",
		utils.Purple, utils.Lines, utils.Reset,
		utils.Purple, utils.Lines, utils.Reset,
		utils.Purple, utils.Lines, utils.Reset,
		utils.Purple, bannerStr, utils.Reset,
		utils.Red, "Made by Majir96", utils.Reset,
		utils.Blue, "Under GPL License", utils.Reset,
		utils.Purple, utils.Lines, utils.Reset,
		utils.Purple, utils.Lines, utils.Reset,
		utils.Purple, utils.Lines, utils.Reset,
	)

	fmt.Println(fullBanner)
}
