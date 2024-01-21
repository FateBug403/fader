package cmd

type CliParam struct {
	InputFile string
	CollectOutput string
	AliveVerify bool

}

var cliParam = &CliParam{}
