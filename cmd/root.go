package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/pranavek/pomodoro/pomo"
	"github.com/spf13/cobra"
)

var (
	workDuration    int
	shortBreak      int
	longBreak       int
	pomosBeforeLong int
	showCountdown   bool
)

var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Pomo helps to implement pomodoro in your workflow",
	Long: `Pomo is a Pomodoro timer CLI tool that helps you stay focused and productive.
It implements the Pomodoro Technique with configurable work sessions, short breaks, and long breaks.

During breaks, it presents reflection prompts to help you think about your work and approach.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate work duration
		if workDuration < 1 || workDuration > 120 {
			return fmt.Errorf("work duration must be between 1 and 120 minutes, got %d", workDuration)
		}

		// Validate short break duration
		if shortBreak < 1 || shortBreak > 60 {
			return fmt.Errorf("short break duration must be between 1 and 60 minutes, got %d", shortBreak)
		}

		// Validate long break duration
		if longBreak < 1 || longBreak > 120 {
			return fmt.Errorf("long break duration must be between 1 and 120 minutes, got %d", longBreak)
		}

		// Validate pomodoros before long break
		if pomosBeforeLong < 1 || pomosBeforeLong > 10 {
			return fmt.Errorf("pomodoros before long break must be between 1 and 10, got %d", pomosBeforeLong)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Construct config from flags
		config := pomo.Config{
			WorkDuration:        time.Duration(workDuration) * time.Minute,
			ShortBreakDuration:  time.Duration(shortBreak) * time.Minute,
			LongBreakDuration:   time.Duration(longBreak) * time.Minute,
			PomosUntilLongBreak: pomosBeforeLong,
			ShowCountdown:       showCountdown,
		}

		// Run the Pomodoro timer
		pomo.Run(config)
	},
}

func init() {
	// Work session duration flag
	rootCmd.Flags().IntVarP(&workDuration, "work", "w", 25,
		"Work session duration in minutes (1-120)")

	// Short break duration flag
	rootCmd.Flags().IntVarP(&shortBreak, "short-break", "s", 5,
		"Short break duration in minutes (1-60)")

	// Long break duration flag
	rootCmd.Flags().IntVarP(&longBreak, "long-break", "l", 30,
		"Long break duration in minutes (1-120)")

	// Pomodoros before long break flag
	rootCmd.Flags().IntVarP(&pomosBeforeLong, "count", "c", 4,
		"Number of pomodoros before a long break (1-10)")

	// Countdown display flag
	rootCmd.Flags().BoolVarP(&showCountdown, "countdown", "d", true,
		"Show real-time countdown during sessions")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
