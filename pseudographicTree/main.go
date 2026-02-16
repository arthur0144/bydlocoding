package main

import (
	"fmt"
	"strings"
)

const (
	defaultSize             = 5
	defaultBackgroundSymbol = " "
	defaultTreeSymbol       = "*"

	topLeftSymbol     = "+"
	topRightSymbol    = "+"
	bottomLeftSymbol  = "+"
	bottomRightSymbol = "+"
	verticalSymbol    = "|"
	horizontalSymbol  = "-"
)

func main() {
	var size int
	size = defaultSize

	backgroundSymbol, treeSymbol := defaultBackgroundSymbol, defaultTreeSymbol

	fmt.Printf("%s%s%s\n", topLeftSymbol, strings.Repeat(horizontalSymbol, 2*size+1), topRightSymbol)
	nBackgroundSymbols, nTreeSymbols := 0, -1
	for i := range size {
		fmt.Printf("%s", verticalSymbol)

		nBackgroundSymbols = size - i
		for range nBackgroundSymbols {
			fmt.Printf("%s", backgroundSymbol)
		}

		nTreeSymbols += 2
		for range nTreeSymbols {
			fmt.Printf("%s", treeSymbol)
		}

		for range nBackgroundSymbols {
			fmt.Printf("%s", backgroundSymbol)
		}

		fmt.Printf("%s", verticalSymbol)
		fmt.Println()
	}
	fmt.Printf("%s%s%s\n", bottomLeftSymbol, strings.Repeat(horizontalSymbol, 2*size+1), bottomRightSymbol)
}
