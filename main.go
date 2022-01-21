package main

import (
	"anylinker/core/cmd"
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "anylinker"}
	//rootCmd.AddCommand(cmd.Client())
	rootCmd.AddCommand(cmd.Server())
	//rootCmd.AddCommand(cmd.Version())
	//rootCmd.AddCommand(cmd.GeneratePemKey())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute failed", err.Error())
	}
}

