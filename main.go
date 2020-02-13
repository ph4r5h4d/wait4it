package main

import (
	"wait4it/cmd"
	"wait4it/inputParser"
)

func main() {
	c := inputParser.GetInput()
	cmd.RunCheck(c)
}
