package source

import (
	"strings"
	"fmt"
)

const field1MaxWidth = 30
const field2MaxWidth = 100
const fieldPadding = 2


func (s StrMatrix) PrintList () {
	for _, datum := range s {
		fmt.Println(datum[0], datum[1])
	}
}

func (s StrMatrix) PrintTable () {
	maxwidthF1, maxwidthF2 := 1, 1

	for _, datum := range s {
		if len(datum[0]) > maxwidthF1 {
			maxwidthF1 = len(datum[0])
		}
		if len(datum[1]) > maxwidthF2 {
			maxwidthF2 = len(datum[1])
		}
	}

	if maxwidthF1 > field1MaxWidth {
		maxwidthF1 = field1MaxWidth
	}
	if maxwidthF2 > field2MaxWidth {
		maxwidthF2 = field2MaxWidth
	}

	maxwidthF1 = maxwidthF1 + fieldPadding
	maxwidthF2 = maxwidthF2 + fieldPadding

	strFormat := fmt.Sprintf("│%%%ds │ %%-%ds│\n", maxwidthF1, maxwidthF2)

	fmt.Println("┌" + strings.Repeat("─", maxwidthF1+1) + "┬" + strings.Repeat("─", maxwidthF2+1) + "┐")

	for _, datum := range s {
		f1 := breakString(datum[0])[0] // truncate field 1
		f2 := breakString(datum[1])
		for _, ff2 := range f2 {
			fmt.Printf(strFormat, f1, ff2)
			f1 = ""
		}
	}

	fmt.Println("└" + strings.Repeat("─", maxwidthF1+1) + "┴" + strings.Repeat("─", maxwidthF2+1) + "┘")
}

func breakString(bigString string) []string {
	if len(bigString) < field2MaxWidth {
		return []string{bigString}
	} else {
		return append([]string{bigString[0:field2MaxWidth-1]}, breakString(bigString[field2MaxWidth:])...)
	}
}
