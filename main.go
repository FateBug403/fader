/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/FateBug403/fader/cmd"
	"github.com/FateBug403/fader/initialize"
)

func main() {
	err := initialize.RunInitialize()
	if err != nil {
		return
	}
	cmd.Execute()
}
