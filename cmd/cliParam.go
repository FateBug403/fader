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

type TestParam struct {
	Proxy bool

}

var collectParam = &CollectParam{}
var cdnParam = &CdnParam{}

var testParam = &TestParam{}
