package pomo

import (
	"fmt"
	"math"
	"time"
)

// TimeOfDayStats aggregates statistics by time of day.
type TimeOfDayStats struct {
	Morning   *ReportStats // 6am - 12pm
	Afternoon *ReportStats // 12pm - 6pm
	Evening   *ReportStats // 6pm - 12am
	Night     *ReportStats // 12am - 6am
}

// DayOfWeekStats aggregates statistics by day of week.
type DayOfWeekStats struct {
	Monday    *ReportStats
	Tuesday   *ReportStats
	Wednesday *ReportStats
	Thursday  *ReportStats
	Friday    *ReportStats
	Saturday  *ReportStats
	Sunday    *ReportStats
}

// ComparisonStats compares two time periods.
type ComparisonStats struct {
	Current        *ReportStats
	Previous       *ReportStats
	PomoChange     int     // +/- change in pomodoros
	PercentChange  float64 // percentage change
	Trend          string  // "improving", "declining", "stable"
	EfficiencyCurr float64
	EfficiencyPrev float64
}

// ProductivityInsights provides comprehensive productivity analysis.
type ProductivityInsights struct {
	TimeOfDay        *TimeOfDayStats
	DayOfWeek        *DayOfWeekStats
	BestTimeSlot     string
	BestDay          string
	AvgDailyPomos    float64
	FocusEfficiency  float64 // CompletedPomos / (CompletedPomos + SkippedSessions)
	WorkBreakRatio   float64 // WorkTime / BreakTime
	ConsistencyScore float64 // Based on session distribution
}

// getTimeOfDay returns the time slot for a given timestamp.
func getTimeOfDay(t time.Time) string {
	hour := t.Hour()

	switch {
	case hour >= 6 && hour < 12:
		return "morning"
	case hour >= 12 && hour < 18:
		return "afternoon"
	case hour >= 18 && hour < 24:
		return "evening"
	default: // 0-5
		return "night"
	}
}

// AnalyzeTimeOfDay groups sessions by time of day and returns statistics.
func AnalyzeTimeOfDay(records []SessionRecord) *TimeOfDayStats {
	morningRecords := []SessionRecord{}
	afternoonRecords := []SessionRecord{}
	eveningRecords := []SessionRecord{}
	nightRecords := []SessionRecord{}

	for _, record := range records {
		timeSlot := getTimeOfDay(record.Date)

		switch timeSlot {
		case "morning":
			morningRecords = append(morningRecords, record)
		case "afternoon":
			afternoonRecords = append(afternoonRecords, record)
		case "evening":
			eveningRecords = append(eveningRecords, record)
		case "night":
			nightRecords = append(nightRecords, record)
		}
	}

	return &TimeOfDayStats{
		Morning:   GenerateReport(morningRecords),
		Afternoon: GenerateReport(afternoonRecords),
		Evening:   GenerateReport(eveningRecords),
		Night:     GenerateReport(nightRecords),
	}
}

// AnalyzeDayOfWeek groups sessions by day of week and returns statistics.
func AnalyzeDayOfWeek(records []SessionRecord) *DayOfWeekStats {
	dayRecords := make(map[time.Weekday][]SessionRecord)

	for _, record := range records {
		weekday := record.Date.Weekday()
		dayRecords[weekday] = append(dayRecords[weekday], record)
	}

	return &DayOfWeekStats{
		Monday:    GenerateReport(dayRecords[time.Monday]),
		Tuesday:   GenerateReport(dayRecords[time.Tuesday]),
		Wednesday: GenerateReport(dayRecords[time.Wednesday]),
		Thursday:  GenerateReport(dayRecords[time.Thursday]),
		Friday:    GenerateReport(dayRecords[time.Friday]),
		Saturday:  GenerateReport(dayRecords[time.Saturday]),
		Sunday:    GenerateReport(dayRecords[time.Sunday]),
	}
}

// GetBestPerformingTime identifies the most productive time slot.
func GetBestPerformingTime(todStats *TimeOfDayStats) string {
	timeSlots := map[string]*ReportStats{
		"Morning (6am-12pm)":   todStats.Morning,
		"Afternoon (12pm-6pm)": todStats.Afternoon,
		"Evening (6pm-12am)":   todStats.Evening,
		"Night (12am-6am)":     todStats.Night,
	}

	bestSlot := ""
	bestAvg := 0.0

	for slot, stats := range timeSlots {
		if stats.TotalSessions > 0 && stats.AveragePomos > bestAvg {
			bestAvg = stats.AveragePomos
			bestSlot = slot
		}
	}

	if bestSlot == "" {
		return "No data available"
	}

	return fmt.Sprintf("%s (%.1f avg pomodoros)", bestSlot, bestAvg)
}

// GetBestPerformingDay identifies the most productive day of week.
func GetBestPerformingDay(dowStats *DayOfWeekStats) string {
	days := map[string]*ReportStats{
		"Monday":    dowStats.Monday,
		"Tuesday":   dowStats.Tuesday,
		"Wednesday": dowStats.Wednesday,
		"Thursday":  dowStats.Thursday,
		"Friday":    dowStats.Friday,
		"Saturday":  dowStats.Saturday,
		"Sunday":    dowStats.Sunday,
	}

	bestDay := ""
	bestAvg := 0.0

	for day, stats := range days {
		if stats.TotalSessions > 0 && stats.AveragePomos > bestAvg {
			bestAvg = stats.AveragePomos
			bestDay = day
		}
	}

	if bestDay == "" {
		return "No data available"
	}

	return fmt.Sprintf("%s (%.1f avg pomodoros)", bestDay, bestAvg)
}

// CompareWeeks compares current week to previous week.
func CompareWeeks(storage *Storage) (*ComparisonStats, error) {
	currentWeekStart := GetWeekStart()
	previousWeekStart := GetPreviousWeekStart()

	// Get current week records
	currentRecords, err := storage.GetRecordsSince(currentWeekStart)
	if err != nil {
		return nil, err
	}

	// Get previous week records
	previousRecords, err := storage.GetRecordsInRange(
		previousWeekStart,
		currentWeekStart.Add(-1*time.Second),
	)
	if err != nil {
		return nil, err
	}

	current := GenerateReport(currentRecords)
	previous := GenerateReport(previousRecords)

	return buildComparisonStats(current, previous), nil
}

// CompareMonths compares current month to previous month.
func CompareMonths(storage *Storage) (*ComparisonStats, error) {
	currentMonthStart := GetMonthStart()
	previousMonthStart := GetPreviousMonthStart()

	// Get current month records
	currentRecords, err := storage.GetRecordsSince(currentMonthStart)
	if err != nil {
		return nil, err
	}

	// Get previous month records
	previousRecords, err := storage.GetRecordsInRange(
		previousMonthStart,
		currentMonthStart.Add(-1*time.Second),
	)
	if err != nil {
		return nil, err
	}

	current := GenerateReport(currentRecords)
	previous := GenerateReport(previousRecords)

	return buildComparisonStats(current, previous), nil
}

// buildComparisonStats creates comparison statistics from two report periods.
func buildComparisonStats(current, previous *ReportStats) *ComparisonStats {
	pomoChange := current.TotalPomos - previous.TotalPomos
	var percentChange float64

	if previous.TotalPomos > 0 {
		percentChange = (float64(pomoChange) / float64(previous.TotalPomos)) * 100
	} else if current.TotalPomos > 0 {
		percentChange = 100.0 // From 0 to something is 100% increase
	}

	trend := "stable"
	if percentChange > 10 {
		trend = "improving"
	} else if percentChange < -10 {
		trend = "declining"
	}

	return &ComparisonStats{
		Current:        current,
		Previous:       previous,
		PomoChange:     pomoChange,
		PercentChange:  percentChange,
		Trend:          trend,
		EfficiencyCurr: CalculateFocusEfficiency(current),
		EfficiencyPrev: CalculateFocusEfficiency(previous),
	}
}

// CalculateFocusEfficiency computes completion rate percentage.
func CalculateFocusEfficiency(stats *ReportStats) float64 {
	totalActivities := stats.TotalPomos + stats.TotalSkipped
	if totalActivities == 0 {
		return 0
	}
	return (float64(stats.TotalPomos) / float64(totalActivities)) * 100
}

// CalculateWorkBreakRatio computes work time to break time ratio.
func CalculateWorkBreakRatio(stats *ReportStats) float64 {
	if stats.TotalBreakTime == 0 {
		return 0
	}
	return float64(stats.TotalWorkTime) / float64(stats.TotalBreakTime)
}

// CalculateConsistencyScore evaluates session regularity (0-100).
func CalculateConsistencyScore(records []SessionRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	// Group pomodoros by date
	pomosByDate := make(map[string]int)
	for _, record := range records {
		dateKey := record.Date.Format("2006-01-02")
		pomosByDate[dateKey] += record.CompletedPomos
	}

	// Calculate standard deviation of daily pomodoro counts
	values := make([]float64, 0, len(pomosByDate))
	var sum float64
	for _, count := range pomosByDate {
		val := float64(count)
		values = append(values, val)
		sum += val
	}

	if len(values) == 0 {
		return 0
	}

	mean := sum / float64(len(values))

	var variance float64
	for _, val := range values {
		diff := val - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	stdDev := math.Sqrt(variance)

	// Score based on coefficient of variation (lower is better/more consistent)
	// CV = stdDev / mean
	if mean == 0 {
		return 0
	}

	cv := stdDev / mean

	// Convert to 0-100 score (lower CV = higher score)
	// CV of 0 = 100 points, CV of 1 = 50 points, CV > 2 = 0 points
	score := 100 - (cv * 50)
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// GenerateProductivityInsights creates comprehensive analytics.
func GenerateProductivityInsights(records []SessionRecord, storage *Storage) *ProductivityInsights {
	if len(records) == 0 {
		return &ProductivityInsights{
			TimeOfDay: &TimeOfDayStats{
				Morning:   GenerateReport([]SessionRecord{}),
				Afternoon: GenerateReport([]SessionRecord{}),
				Evening:   GenerateReport([]SessionRecord{}),
				Night:     GenerateReport([]SessionRecord{}),
			},
			DayOfWeek: &DayOfWeekStats{
				Monday:    GenerateReport([]SessionRecord{}),
				Tuesday:   GenerateReport([]SessionRecord{}),
				Wednesday: GenerateReport([]SessionRecord{}),
				Thursday:  GenerateReport([]SessionRecord{}),
				Friday:    GenerateReport([]SessionRecord{}),
				Saturday:  GenerateReport([]SessionRecord{}),
				Sunday:    GenerateReport([]SessionRecord{}),
			},
		}
	}

	todStats := AnalyzeTimeOfDay(records)
	dowStats := AnalyzeDayOfWeek(records)
	stats := GenerateReport(records)

	// Calculate average daily pomodoros
	dateSet := make(map[string]bool)
	for _, record := range records {
		dateKey := record.Date.Format("2006-01-02")
		dateSet[dateKey] = true
	}
	uniqueDays := len(dateSet)
	avgDaily := 0.0
	if uniqueDays > 0 {
		avgDaily = float64(stats.TotalPomos) / float64(uniqueDays)
	}

	return &ProductivityInsights{
		TimeOfDay:        todStats,
		DayOfWeek:        dowStats,
		BestTimeSlot:     GetBestPerformingTime(todStats),
		BestDay:          GetBestPerformingDay(dowStats),
		AvgDailyPomos:    avgDaily,
		FocusEfficiency:  CalculateFocusEfficiency(stats),
		WorkBreakRatio:   CalculateWorkBreakRatio(stats),
		ConsistencyScore: CalculateConsistencyScore(records),
	}
}

// DisplayTimeOfDayAnalysis prints time-of-day breakdown.
func DisplayTimeOfDayAnalysis(todStats *TimeOfDayStats) {
	fmt.Println("\n  Time of Day Analysis")
	fmt.Println("  ────────────────────────────────────────────")

	displayTimeSlot("Morning (6am-12pm)", todStats.Morning)
	displayTimeSlot("Afternoon (12pm-6pm)", todStats.Afternoon)
	displayTimeSlot("Evening (6pm-12am)", todStats.Evening)
	displayTimeSlot("Night (12am-6am)", todStats.Night)

	fmt.Println()
	fmt.Printf("  Best performing time: %s\n", GetBestPerformingTime(todStats))
}

// displayTimeSlot prints statistics for a single time slot.
func displayTimeSlot(name string, stats *ReportStats) {
	if stats.TotalSessions == 0 {
		fmt.Printf("  %s:\n    No sessions\n", name)
		return
	}

	fmt.Printf("  %s:\n", name)
	fmt.Printf("    Sessions: %d | Pomodoros: %d | Avg: %.1f\n",
		stats.TotalSessions,
		stats.TotalPomos,
		stats.AveragePomos)
}

// DisplayDayOfWeekAnalysis prints day-of-week breakdown.
func DisplayDayOfWeekAnalysis(dowStats *DayOfWeekStats) {
	fmt.Println("\n  Day of Week Analysis")
	fmt.Println("  ────────────────────────────────────────────")

	displayDaySlot("Monday", dowStats.Monday)
	displayDaySlot("Tuesday", dowStats.Tuesday)
	displayDaySlot("Wednesday", dowStats.Wednesday)
	displayDaySlot("Thursday", dowStats.Thursday)
	displayDaySlot("Friday", dowStats.Friday)
	displayDaySlot("Saturday", dowStats.Saturday)
	displayDaySlot("Sunday", dowStats.Sunday)

	fmt.Println()
	fmt.Printf("  Most productive day: %s\n", GetBestPerformingDay(dowStats))
}

// displayDaySlot prints statistics for a single day.
func displayDaySlot(name string, stats *ReportStats) {
	if stats.TotalSessions == 0 {
		fmt.Printf("  %-10s No sessions\n", name+":")
		return
	}

	fmt.Printf("  %-10s Sessions: %-3d | Pomodoros: %-3d | Avg: %.1f\n",
		name+":",
		stats.TotalSessions,
		stats.TotalPomos,
		stats.AveragePomos)
}

// DisplayComparison prints period comparison.
func DisplayComparison(comp *ComparisonStats, title string) {
	fmt.Printf("\n  %s\n", title)
	fmt.Println("  ════════════════════════════════════════════")

	if comp.Current.TotalSessions == 0 && comp.Previous.TotalSessions == 0 {
		fmt.Println("  No data available for comparison")
		return
	}

	fmt.Printf("  This Period:  %d pomodoros | %s work time\n",
		comp.Current.TotalPomos,
		formatDuration(comp.Current.TotalWorkTime))

	fmt.Printf("  Last Period:  %d pomodoros | %s work time\n",
		comp.Previous.TotalPomos,
		formatDuration(comp.Previous.TotalWorkTime))

	fmt.Println()

	// Show change
	changeSign := "+"
	if comp.PomoChange < 0 {
		changeSign = ""
	}
	fmt.Printf("  Change:       %s%d pomodoros (%+.1f%%)\n",
		changeSign,
		comp.PomoChange,
		comp.PercentChange)

	// Show trend
	trendSymbol := "→"
	if comp.Trend == "improving" {
		trendSymbol = "✓"
	} else if comp.Trend == "declining" {
		trendSymbol = "↓"
	}
	fmt.Printf("  Trend:        %s %s\n", trendSymbol, comp.Trend)

	// Show focus efficiency change
	effChange := comp.EfficiencyCurr - comp.EfficiencyPrev
	if comp.Previous.TotalSessions > 0 {
		fmt.Printf("  Focus Rate:   %.0f%% (", comp.EfficiencyCurr)
		if effChange >= 0 {
			fmt.Printf("up %.0f%% from %.0f%%)\n", effChange, comp.EfficiencyPrev)
		} else {
			fmt.Printf("down %.0f%% from %.0f%%)\n", -effChange, comp.EfficiencyPrev)
		}
	}
}

// DisplayProductivityInsights prints formatted insights.
func DisplayProductivityInsights(insights *ProductivityInsights, title string, stats *ReportStats) {
	fmt.Printf("\n📊 %s\n", title)
	fmt.Println("  ════════════════════════════════════════════")

	if stats.TotalSessions == 0 {
		fmt.Println("  No sessions recorded yet")
		return
	}

	fmt.Printf("  Sessions:        %d sessions\n", stats.TotalSessions)
	fmt.Printf("  Total Pomos:     %d pomodoros\n", stats.TotalPomos)

	// Focus efficiency rating
	effRating := "needs improvement"
	if insights.FocusEfficiency >= 90 {
		effRating = "excellent"
	} else if insights.FocusEfficiency >= 75 {
		effRating = "good"
	}
	fmt.Printf("  Focus Rate:      %.0f%% (%s)\n", insights.FocusEfficiency, effRating)

	// Work/break ratio
	if insights.WorkBreakRatio > 0 {
		ratioRating := "optimal"
		if insights.WorkBreakRatio < 4 || insights.WorkBreakRatio > 6 {
			ratioRating = "off-target"
		}
		fmt.Printf("  Work/Break:      %.1f:1 (%s)\n", insights.WorkBreakRatio, ratioRating)
	}

	// Consistency score
	consistencyRating := "needs work"
	if insights.ConsistencyScore >= 80 {
		consistencyRating = "excellent"
	} else if insights.ConsistencyScore >= 60 {
		consistencyRating = "good"
	}
	fmt.Printf("  Consistency:     %.0f/100 (%s)\n", insights.ConsistencyScore, consistencyRating)

	fmt.Println()
	fmt.Printf("  Best Time:       %s\n", insights.BestTimeSlot)
	fmt.Printf("  Best Day:        %s\n", insights.BestDay)
	fmt.Printf("  Avg Daily:       %.1f pomodoros\n", insights.AvgDailyPomos)
}
