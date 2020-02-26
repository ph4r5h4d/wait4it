package main

import (
	"wait4it/cmd"
	"wait4it/inputParser/envParser"
	"wait4it/inputParser/flagParser"
	"wait4it/model"
)

func main() {
	c := &model.CheckContext{}
	c = envParser.Parse(c)
	c = flagParser.Parse(c)
	cmd.RunCheck(*c)
}
