package cmd

import (
	"fmt"
	"os"

	"github.com/pranavek/pomodoro/pomo"
	"github.com/spf13/cobra"
)

var (
	analyzeWeek  bool
	analyzeMonth bool
	analyzeAll   bool
	compareWeeks bool
	compareMonths bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze your productivity patterns and trends",
	Long: `Provides comprehensive productivity analysis including best times, day patterns,
comparisons, and insights to help you evaluate and improve your productivity.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help when called without subcommand
		cmd.Help()
	},
}

var insightsCmd = &cobra.Command{
	Use:   "insights",
	Short: "Show comprehensive productivity insights",
	Long:  `Displays an overview of your productivity including efficiency metrics, best times and days, and overall patterns.`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := pomo.NewStorage()
		if err != nil {
			fmt.Printf("Error: Could not access storage: %v\n", err)
			os.Exit(1)
		}
		defer storage.Close()

		var records []pomo.SessionRecord
		var title string

		// Determine time range
		switch {
		case analyzeAll:
			records, err = storage.LoadRecords()
			title = "Productivity Insights - All Time"
		case analyzeMonth:
			records, err = storage.GetRecordsSince(pomo.GetMonthStart())
			title = "Productivity Insights - This Month"
		default: // week is default
			records, err = storage.GetRecordsSince(pomo.GetWeekStart())
			title = "Productivity Insights - This Week"
		}

		if err != nil {
			fmt.Printf("Error: Could not load records: %v\n", err)
			os.Exit(1)
		}

		// Generate insights
		insights := pomo.GenerateProductivityInsights(records, storage)
		stats := pomo.GenerateReport(records)

		// Display
		pomo.DisplayProductivityInsights(insights, title, stats)

		// Check for streak
		streak, err := pomo.CalculateStreak(storage)
		if err == nil && streak.CurrentStreak > 0 {
			fmt.Println()
			fmt.Printf("  Streak:          %d days", streak.CurrentStreak)
			if streak.CurrentStreak >= 7 {
				fmt.Print(" 🔥")
			}
			fmt.Println()
		}

		// Check for active goals and show progress
		goalConfig, err := pomo.LoadGoalConfig()
		if err == nil && goalConfig.Enabled {
			if goalConfig.WeeklyPomosGoal > 0 && !analyzeAll && !analyzeMonth {
				fmt.Println()
				progress, err := pomo.CheckWeeklyGoal(storage, goalConfig.WeeklyPomosGoal)
				if err == nil {
					fmt.Printf("  Weekly Goal:     %d/%d (%.0f%%)\n",
						progress.Completed,
						progress.Goal,
						progress.Percentage)
				}
			}
		}

		fmt.Println()
	},
}

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Analyze productivity by time of day",
	Long:  `Shows when you're most productive by breaking down your sessions into morning, afternoon, evening, and night.`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := pomo.NewStorage()
		if err != nil {
			fmt.Printf("Error: Could not access storage: %v\n", err)
			os.Exit(1)
		}
		defer storage.Close()

		var records []pomo.SessionRecord
		var title string

		// Determine time range
		switch {
		case analyzeAll:
			records, err = storage.LoadRecords()
			title = "Time of Day Analysis - All Time"
		case analyzeMonth:
			records, err = storage.GetRecordsSince(pomo.GetMonthStart())
			title = "Time of Day Analysis - This Month"
		default: // week is default
			records, err = storage.GetRecordsSince(pomo.GetWeekStart())
			title = "Time of Day Analysis - This Week"
		}

		if err != nil {
			fmt.Printf("Error: Could not load records: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n📊 %s\n", title)
		todStats := pomo.AnalyzeTimeOfDay(records)
		pomo.DisplayTimeOfDayAnalysis(todStats)
		fmt.Println()
	},
}

var daysCmd = &cobra.Command{
	Use:   "days",
	Short: "Analyze productivity by day of week",
	Long:  `Shows which days you're most productive by breaking down your sessions across Monday through Sunday.`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := pomo.NewStorage()
		if err != nil {
			fmt.Printf("Error: Could not access storage: %v\n", err)
			os.Exit(1)
		}
		defer storage.Close()

		var records []pomo.SessionRecord
		var title string

		// Determine time range
		switch {
		case analyzeAll:
			records, err = storage.LoadRecords()
			title = "Day of Week Analysis - All Time"
		case analyzeMonth:
			records, err = storage.GetRecordsSince(pomo.GetMonthStart())
			title = "Day of Week Analysis - This Month"
		default: // week is default
			records, err = storage.GetRecordsSince(pomo.GetWeekStart())
			title = "Day of Week Analysis - This Week"
		}

		if err != nil {
			fmt.Printf("Error: Could not load records: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n📊 %s\n", title)
		dowStats := pomo.AnalyzeDayOfWeek(records)
		pomo.DisplayDayOfWeekAnalysis(dowStats)
		fmt.Println()
	},
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare productivity across time periods",
	Long:  `Compares your current productivity to previous periods (week-over-week or month-over-month) to show trends.`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := pomo.NewStorage()
		if err != nil {
			fmt.Printf("Error: Could not access storage: %v\n", err)
			os.Exit(1)
		}
		defer storage.Close()

		var comp *pomo.ComparisonStats
		var title string

		if compareMonths {
			comp, err = pomo.CompareMonths(storage)
			title = "Month-over-Month Comparison"
		} else {
			// Default to weeks
			comp, err = pomo.CompareWeeks(storage)
			title = "Week-over-Week Comparison"
		}

		if err != nil {
			fmt.Printf("Error: Could not compare periods: %v\n", err)
			os.Exit(1)
		}

		pomo.DisplayComparison(comp, title)
		fmt.Println()
	},
}

var streakCmd = &cobra.Command{
	Use:   "streak",
	Short: "Show your consistency streak",
	Long:  `Displays your current and longest activity streaks, showing how many consecutive days you've completed pomodoros.`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := pomo.NewStorage()
		if err != nil {
			fmt.Printf("Error: Could not access storage: %v\n", err)
			os.Exit(1)
		}
		defer storage.Close()

		streak, err := pomo.CalculateStreak(storage)
		if err != nil {
			fmt.Printf("Error: Could not calculate streak: %v\n", err)
			os.Exit(1)
		}

		pomo.DisplayStreak(streak)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// Add subcommands
	analyzeCmd.AddCommand(insightsCmd)
	analyzeCmd.AddCommand(timeCmd)
	analyzeCmd.AddCommand(daysCmd)
	analyzeCmd.AddCommand(compareCmd)
	analyzeCmd.AddCommand(streakCmd)

	// Time range flags for insights, time, and days commands
	insightsCmd.Flags().BoolVar(&analyzeWeek, "week", false, "Analyze this week (default)")
	insightsCmd.Flags().BoolVar(&analyzeMonth, "month", false, "Analyze this month")
	insightsCmd.Flags().BoolVar(&analyzeAll, "all", false, "Analyze all time")

	timeCmd.Flags().BoolVar(&analyzeWeek, "week", false, "Analyze this week (default)")
	timeCmd.Flags().BoolVar(&analyzeMonth, "month", false, "Analyze this month")
	timeCmd.Flags().BoolVar(&analyzeAll, "all", false, "Analyze all time")

	daysCmd.Flags().BoolVar(&analyzeWeek, "week", false, "Analyze this week (default)")
	daysCmd.Flags().BoolVar(&analyzeMonth, "month", false, "Analyze this month")
	daysCmd.Flags().BoolVar(&analyzeAll, "all", false, "Analyze all time")

	// Comparison type flags
	compareCmd.Flags().BoolVar(&compareWeeks, "weeks", false, "Compare weeks (default)")
	compareCmd.Flags().BoolVar(&compareMonths, "months", false, "Compare months")
}
