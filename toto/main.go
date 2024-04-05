package main

import (
	"github.com/spf13/toto/cmd"
	"github.com/common-nighthawk/go-figure"
)


func main() {
	figure.NewFigure("Welcome!", "basic", true).Scroll(3000, 3000, "right")
	fmt.Printf("\n\n")
	cmd.Execute()
}
