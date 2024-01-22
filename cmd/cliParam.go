package cmd

type CollectParam struct {
	InputFile string
	CollectOutput string
	AliveVerify bool

}

type CdnParam struct {
	InputFile string
	OutPut string
}

var collectParam = &CollectParam{}
var cdnParam = &CdnParam{}
