package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/bendersilver/pgeasy"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := pgeasy.InitConf()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pgeasy.Start()
	if err != nil {
		fmt.Println(err)
	}
}
