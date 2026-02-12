package pomo

import (
	"fmt"
	"time"
)

// ReportStats aggregates statistics from multiple session records.
type ReportStats struct {
	TotalSessions   int
	TotalPomos      int
	TotalSkipped    int
	TotalWorkTime   time.Duration
	TotalBreakTime  time.Duration
	TotalDuration   time.Duration
	AveragePomos    float64
	Records         []SessionRecord
}

// GenerateReport creates a report from the given records.
func GenerateReport(records []SessionRecord) *ReportStats {
	if len(records) == 0 {
		return &ReportStats{}
	}

	stats := &ReportStats{
		TotalSessions: len(records),
		Records:       records,
	}

	for _, record := range records {
		stats.TotalPomos += record.CompletedPomos
		stats.TotalSkipped += record.SkippedSessions
		stats.TotalWorkTime += record.WorkTime
		stats.TotalBreakTime += record.BreakTime
		stats.TotalDuration += record.Duration
	}

	if stats.TotalSessions > 0 {
		stats.AveragePomos = float64(stats.TotalPomos) / float64(stats.TotalSessions)
	}

	return stats
}

// DisplayReport prints a formatted report to the console.
func DisplayReport(stats *ReportStats, title string) {
	if stats.TotalSessions == 0 {
		fmt.Printf("\n📊 %s\n\n", title)
		fmt.Println("  No sessions recorded yet. Start a pomodoro to begin tracking!")
		return
	}

	fmt.Printf("\n📊 %s\n", title)
	fmt.Println("  ════════════════════════════════════════════")
	fmt.Printf("  Total sessions: %d\n", stats.TotalSessions)
	fmt.Printf("  Total pomodoros: %d\n", stats.TotalPomos)
	if stats.TotalSkipped > 0 {
		fmt.Printf("  Sessions skipped: %d\n", stats.TotalSkipped)
	}
	fmt.Printf("  Average pomodoros per session: %.1f\n", stats.AveragePomos)
	fmt.Println()
	fmt.Printf("  Total work time: %s\n", formatDuration(stats.TotalWorkTime))
	fmt.Printf("  Total break time: %s\n", formatDuration(stats.TotalBreakTime))
	fmt.Printf("  Total time: %s\n", formatDuration(stats.TotalDuration))
	fmt.Println("  ════════════════════════════════════════════")
}

// DisplayDetailedReport shows a report with individual session details.
func DisplayDetailedReport(stats *ReportStats, title string) {
	DisplayReport(stats, title)

	if len(stats.Records) == 0 {
		return
	}

	fmt.Println("\n  Recent Sessions:")
	fmt.Println("  ────────────────────────────────────────────")

	// Show last 10 sessions
	start := len(stats.Records) - 10
	if start < 0 {
		start = 0
	}

	for i := len(stats.Records) - 1; i >= start; i-- {
		record := stats.Records[i]
		date := record.Date.Format("Jan 02, 2006 15:04")

		if record.Title != "" {
			fmt.Printf("  %s - %s\n", date, record.Title)
			fmt.Printf("    %d 🍅 (%s work)\n",
				record.CompletedPomos,
				formatDuration(record.WorkTime))
		} else {
			fmt.Printf("  %s - %d 🍅 (%s work)\n",
				date,
				record.CompletedPomos,
				formatDuration(record.WorkTime))
		}
	}
	fmt.Println()
}

// GetTodayStart returns the start of today (midnight).
func GetTodayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// GetWeekStart returns the start of the current week (Monday).
func GetWeekStart() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	daysBack := weekday - 1 // Days back to Monday
	monday := now.AddDate(0, 0, -daysBack)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location())
}

// GetMonthStart returns the start of the current month.
func GetMonthStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

// GetYearStart returns the start of the current year.
func GetYearStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
}

// GetPreviousWeekStart returns the start of the previous week (Monday).
func GetPreviousWeekStart() time.Time {
	currentWeekStart := GetWeekStart()
	return currentWeekStart.AddDate(0, 0, -7)
}

// GetPreviousMonthStart returns the start of the previous month.
func GetPreviousMonthStart() time.Time {
	now := time.Now()
	// Go back one month
	previousMonth := now.AddDate(0, -1, 0)
	return time.Date(previousMonth.Year(), previousMonth.Month(), 1, 0, 0, 0, 0, now.Location())
}
