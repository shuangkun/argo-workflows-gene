package main

import (
	"fmt"
	"os"

	"github.com/shuangkun/argo-workflows-gene/cmd"
)

func main() {

	if err := cmd.ArgoGeneCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
