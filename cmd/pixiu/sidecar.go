package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	startSideCarCmd = &cobra.Command{
		Use:   "start-sidecar",
		Short: "Run dubbo go pixiu in sidecar mode  (implement in the future)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Meet you in the Future!")
		},
	}
)
