package cmd

import (
	"fmt"
	"os"

	"github.com/pranavek/pomodoro/pomo"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Pomo helps to implement pomodoro in your workflow",
	Run: func(cmd *cobra.Command, args []string) {
		pomo.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
