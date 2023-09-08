package main

import (
	"flag"

	"github.com/sshelll/clstelegraph/service"
)

var (
	nFlag = flag.Int("n", 10, "number of telegraphs to fetch")
	iFlag = flag.Int("i", 30, "interval of refresh screen in milliseconds")
)

func main() {
	flag.Parse()
	m := service.NewCore(*nFlag, *iFlag)
	m.Start()
}
