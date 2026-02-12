package cmd

import (
	"fmt"
	"os"

	"github.com/pranavek/pomodoro/pomo"
	"github.com/spf13/cobra"
)

var (
	goalDaily  int
	goalWeekly int
	showDaily  bool
	showWeekly bool
)

var goalsCmd = &cobra.Command{
	Use:   "goals",
	Short: "Manage and track your pomodoro goals",
	Long:  `Set daily or weekly pomodoro goals, track your progress, and stay motivated to maintain consistent productivity.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help when called without subcommand
		cmd.Help()
	},
}

var setGoalsCmd = &cobra.Command{
	Use:   "set",
	Short: "Set your daily or weekly pomodoro goals",
	Long:  `Configure your target number of pomodoros per day or per week to help track your productivity.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load existing config or create new one
		config, err := pomo.LoadGoalConfig()
		if err != nil {
			fmt.Printf("Error: Could not load goal config: %v\n", err)
			os.Exit(1)
		}

		// Check if any goals were specified
		if goalDaily == 0 && goalWeekly == 0 {
			fmt.Println("Error: Please specify at least one goal")
			fmt.Println("\nUsage:")
			fmt.Println("  pomo goals set --daily N     Set daily goal")
			fmt.Println("  pomo goals set --weekly N    Set weekly goal")
			fmt.Println("  pomo goals set --daily N --weekly N    Set both")
			os.Exit(1)
		}

		// Update config
		if goalDaily > 0 {
			config.DailyPomosGoal = goalDaily
		}
		if goalWeekly > 0 {
			config.WeeklyPomosGoal = goalWeekly
		}
		config.Enabled = true

		// Save config
		if err := pomo.SaveGoalConfig(config); err != nil {
			fmt.Printf("Error: Could not save goal config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✓ Goals updated successfully!")
		pomo.DisplayGoalConfig(config)
		fmt.Println()
	},
}

var showGoalsCmd = &cobra.Command{
	Use:   "show",
	Short: "Display your current goals",
	Long:  `Shows your configured daily and weekly pomodoro goals.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := pomo.LoadGoalConfig()
		if err != nil {
			fmt.Printf("Error: Could not load goal config: %v\n", err)
			os.Exit(1)
		}

		pomo.DisplayGoalConfig(config)
		fmt.Println()
	},
}

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Check your progress toward goals",
	Long:  `Shows how many pomodoros you've completed toward your daily or weekly goals.`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := pomo.NewStorage()
		if err != nil {
			fmt.Printf("Error: Could not access storage: %v\n", err)
			os.Exit(1)
		}
		defer storage.Close()

		config, err := pomo.LoadGoalConfig()
		if err != nil {
			fmt.Printf("Error: Could not load goal config: %v\n", err)
			os.Exit(1)
		}

		if !config.Enabled {
			fmt.Println("\nNo goals configured yet.")
			fmt.Println("Set goals with: pomo goals set --daily N --weekly N\n")
			return
		}

		// If no flags specified, show both (if configured)
		if !showDaily && !showWeekly {
			showDaily = config.DailyPomosGoal > 0
			showWeekly = config.WeeklyPomosGoal > 0
		}

		// Show daily progress
		if showDaily && config.DailyPomosGoal > 0 {
			progress, err := pomo.CheckDailyGoal(storage, config.DailyPomosGoal)
			if err != nil {
				fmt.Printf("Error: Could not check daily goal: %v\n", err)
				os.Exit(1)
			}
			pomo.DisplayGoalProgress(progress, "Daily")
		}

		// Show weekly progress
		if showWeekly && config.WeeklyPomosGoal > 0 {
			progress, err := pomo.CheckWeeklyGoal(storage, config.WeeklyPomosGoal)
			if err != nil {
				fmt.Printf("Error: Could not check weekly goal: %v\n", err)
				os.Exit(1)
			}

			if showDaily && config.DailyPomosGoal > 0 {
				fmt.Println()
			}

			pomo.DisplayGoalProgress(progress, "Weekly")
		}

		// Show streak
		streak, err := pomo.CalculateStreak(storage)
		if err == nil && streak.CurrentStreak > 0 {
			fmt.Println()
			pomo.DisplayStreak(streak)
		}

		fmt.Println()
	},
}

var clearGoalsCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all goals",
	Long:  `Clears your configured goals. Your historical data is preserved.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := &pomo.GoalConfig{
			Enabled: false,
		}

		if err := pomo.SaveGoalConfig(config); err != nil {
			fmt.Printf("Error: Could not clear goals: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✓ Goals cleared successfully\n")
	},
}

func init() {
	rootCmd.AddCommand(goalsCmd)

	// Add subcommands
	goalsCmd.AddCommand(setGoalsCmd)
	goalsCmd.AddCommand(showGoalsCmd)
	goalsCmd.AddCommand(progressCmd)
	goalsCmd.AddCommand(clearGoalsCmd)

	// Flags for set command
	setGoalsCmd.Flags().IntVar(&goalDaily, "daily", 0, "Daily pomodoro goal")
	setGoalsCmd.Flags().IntVar(&goalWeekly, "weekly", 0, "Weekly pomodoro goal")

	// Flags for progress command
	progressCmd.Flags().BoolVar(&showDaily, "daily", false, "Show daily progress")
	progressCmd.Flags().BoolVar(&showWeekly, "weekly", false, "Show weekly progress")
}
