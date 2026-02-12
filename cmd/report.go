package cmd

import (
	"fmt"
	"os"

	"github.com/pranavek/pomodoro/pomo"
	"github.com/spf13/cobra"
)

var (
	reportToday    bool
	reportWeek     bool
	reportMonth    bool
	reportYear     bool
	reportAll      bool
	reportDetailed bool
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report of your pomodoro statistics",
	Long: `Generate a report showing your pomodoro statistics for different time periods.

By default, shows today's statistics. Use flags to view different time ranges:
  --today    Today's sessions (default)
  --week     This week's sessions
  --month    This month's sessions
  --year     This year's sessions
  --all      All recorded sessions
  --detailed Show detailed session list`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := pomo.NewStorage()
		if err != nil {
			fmt.Printf("Error: Could not access storage: %v\n", err)
			os.Exit(1)
		}

		var records []pomo.SessionRecord
		var title string

		// Determine which time range to use
		switch {
		case reportAll:
			records, err = storage.LoadRecords()
			title = "All Time Statistics"
		case reportYear:
			records, err = storage.GetRecordsSince(pomo.GetYearStart())
			title = "This Year's Statistics"
		case reportMonth:
			records, err = storage.GetRecordsSince(pomo.GetMonthStart())
			title = "This Month's Statistics"
		case reportWeek:
			records, err = storage.GetRecordsSince(pomo.GetWeekStart())
			title = "This Week's Statistics"
		default: // today is default
			records, err = storage.GetRecordsSince(pomo.GetTodayStart())
			title = "Today's Statistics"
		}

		if err != nil {
			fmt.Printf("Error: Could not load records: %v\n", err)
			os.Exit(1)
		}

		stats := pomo.GenerateReport(records)

		if reportDetailed {
			pomo.DisplayDetailedReport(stats, title)
		} else {
			pomo.DisplayReport(stats, title)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().BoolVar(&reportToday, "today", false, "Show today's statistics (default)")
	reportCmd.Flags().BoolVar(&reportWeek, "week", false, "Show this week's statistics")
	reportCmd.Flags().BoolVar(&reportMonth, "month", false, "Show this month's statistics")
	reportCmd.Flags().BoolVar(&reportYear, "year", false, "Show this year's statistics")
	reportCmd.Flags().BoolVar(&reportAll, "all", false, "Show all time statistics")
	reportCmd.Flags().BoolVarP(&reportDetailed, "detailed", "d", false, "Show detailed session list")
}
