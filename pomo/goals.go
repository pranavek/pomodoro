package pomo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GoalConfig stores user-defined goals.
type GoalConfig struct {
	DailyPomosGoal  int       `json:"daily_pomodoros_goal"`
	WeeklyPomosGoal int       `json:"weekly_pomodoros_goal"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// GoalProgress tracks progress toward goals.
type GoalProgress struct {
	Goal           int
	Completed      int
	Percentage     float64
	Remaining      int
	OnTrack        bool
	DaysActive     int // For streak tracking
	PercentOfDay   float64
	PercentOfWeek  float64
}

// StreakInfo tracks consistency.
type StreakInfo struct {
	CurrentStreak  int
	LongestStreak  int
	LastActiveDate time.Time
}

// getGoalsFilePath returns the path to the goals configuration file.
func getGoalsFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	pomoDir := filepath.Join(homeDir, ".pomo")
	if err := os.MkdirAll(pomoDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(pomoDir, "goals.json"), nil
}

// LoadGoalConfig reads goals from config file.
func LoadGoalConfig() (*GoalConfig, error) {
	path, err := getGoalsFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return &GoalConfig{Enabled: false}, nil
		}
		return nil, err
	}

	var config GoalConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveGoalConfig writes goals to config file.
func SaveGoalConfig(config *GoalConfig) error {
	path, err := getGoalsFilePath()
	if err != nil {
		return err
	}

	config.UpdatedAt = time.Now()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now()
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// CalculateGoalProgress computes progress toward a goal.
func CalculateGoalProgress(goal int, records []SessionRecord, startDate time.Time, isWeekly bool) *GoalProgress {
	completed := 0
	for _, record := range records {
		completed += record.CompletedPomos
	}

	percentage := 0.0
	if goal > 0 {
		percentage = (float64(completed) / float64(goal)) * 100
	}

	remaining := goal - completed
	if remaining < 0 {
		remaining = 0
	}

	// Calculate on-track status based on elapsed time
	now := time.Now()
	elapsed := now.Sub(startDate)

	var percentElapsed float64

	if isWeekly {
		percentElapsed = (elapsed.Hours() / (7 * 24)) * 100
	} else {
		// For daily: calculate from start of day to end of day
		percentElapsed = (elapsed.Hours() / 24) * 100
	}

	// On track if progress percentage >= elapsed time percentage
	onTrack := percentage >= percentElapsed || completed >= goal

	return &GoalProgress{
		Goal:          goal,
		Completed:     completed,
		Percentage:    percentage,
		Remaining:     remaining,
		OnTrack:       onTrack,
		PercentOfDay:  percentElapsed,
		PercentOfWeek: percentElapsed,
	}
}

// CheckDailyGoal evaluates today's progress.
func CheckDailyGoal(storage *Storage, goal int) (*GoalProgress, error) {
	todayStart := GetTodayStart()
	records, err := storage.GetRecordsSince(todayStart)
	if err != nil {
		return nil, err
	}

	return CalculateGoalProgress(goal, records, todayStart, false), nil
}

// CheckWeeklyGoal evaluates this week's progress.
func CheckWeeklyGoal(storage *Storage, goal int) (*GoalProgress, error) {
	weekStart := GetWeekStart()
	records, err := storage.GetRecordsSince(weekStart)
	if err != nil {
		return nil, err
	}

	return CalculateGoalProgress(goal, records, weekStart, true), nil
}

// CalculateStreak determines current and longest streaks.
func CalculateStreak(storage *Storage) (*StreakInfo, error) {
	records, err := storage.LoadRecords()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return &StreakInfo{}, nil
	}

	// Extract unique dates with completed pomodoros
	dateSet := make(map[string]bool)
	var lastActive time.Time

	for _, record := range records {
		if record.CompletedPomos > 0 {
			dateKey := record.Date.Format("2006-01-02")
			dateSet[dateKey] = true

			if record.Date.After(lastActive) {
				lastActive = record.Date
			}
		}
	}

	if len(dateSet) == 0 {
		return &StreakInfo{}, nil
	}

	// Convert to sorted slice of dates
	dates := make([]time.Time, 0, len(dateSet))
	for dateStr := range dateSet {
		date, _ := time.Parse("2006-01-02", dateStr)
		dates = append(dates, date)
	}

	// Sort dates (bubble sort for simplicity)
	for i := 0; i < len(dates); i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i].After(dates[j]) {
				dates[i], dates[j] = dates[j], dates[i]
			}
		}
	}

	// Calculate current streak (from today backwards)
	today := time.Now()
	todayStr := today.Format("2006-01-02")
	currentStreak := 0

	if dateSet[todayStr] {
		currentStreak = 1
		checkDate := today.AddDate(0, 0, -1)

		for {
			checkStr := checkDate.Format("2006-01-02")
			if !dateSet[checkStr] {
				break
			}
			currentStreak++
			checkDate = checkDate.AddDate(0, 0, -1)
		}
	}

	// Calculate longest streak
	longestStreak := 0
	tempStreak := 1

	for i := 1; i < len(dates); i++ {
		prevDate := dates[i-1]
		currDate := dates[i]

		// Check if consecutive days
		diff := currDate.Sub(prevDate)
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
			tempStreak = 1
		}
	}

	// Check final streak
	if tempStreak > longestStreak {
		longestStreak = tempStreak
	}

	// Ensure current streak doesn't exceed longest
	if currentStreak > longestStreak {
		longestStreak = currentStreak
	}

	return &StreakInfo{
		CurrentStreak:  currentStreak,
		LongestStreak:  longestStreak,
		LastActiveDate: lastActive,
	}, nil
}

// DisplayGoalProgress prints goal tracking information.
func DisplayGoalProgress(progress *GoalProgress, goalType string) {
	fmt.Printf("\n  %s Goal Progress\n", goalType)
	fmt.Println("  ════════════════════════════════════════════")

	fmt.Printf("  Goal:         %d pomodoros\n", progress.Goal)
	fmt.Printf("  Completed:    %d pomodoros (%.0f%%)\n", progress.Completed, progress.Percentage)

	if progress.Remaining > 0 {
		fmt.Printf("  Remaining:    %d pomodoros\n", progress.Remaining)
	} else {
		fmt.Println("  🎉 Goal achieved!")
	}

	statusSymbol := "✓"
	statusText := "On track"
	if !progress.OnTrack && progress.Completed < progress.Goal {
		statusSymbol = "⚠"
		statusText = "Behind schedule"
	} else if progress.Completed >= progress.Goal {
		statusSymbol = "✓"
		statusText = "Goal met"
	}

	fmt.Printf("  Status:       %s %s\n", statusSymbol, statusText)
}

// DisplayStreak prints streak information.
func DisplayStreak(streak *StreakInfo) {
	fmt.Println("\n  Consistency Streak")
	fmt.Println("  ════════════════════════════════════════════")

	if streak.CurrentStreak == 0 {
		fmt.Println("  No active streak")
		if !streak.LastActiveDate.IsZero() {
			fmt.Printf("  Last Active:     %s\n", streak.LastActiveDate.Format("Jan 02, 2006"))
		}
		return
	}

	streakEmoji := ""
	if streak.CurrentStreak >= 7 {
		streakEmoji = " 🔥"
	}

	fmt.Printf("  Current Streak:  %d days%s\n", streak.CurrentStreak, streakEmoji)
	fmt.Printf("  Longest Streak:  %d days\n", streak.LongestStreak)

	lastActiveStr := "Today"
	if !streak.LastActiveDate.IsZero() {
		today := time.Now().Format("2006-01-02")
		lastActive := streak.LastActiveDate.Format("2006-01-02")
		if today != lastActive {
			lastActiveStr = streak.LastActiveDate.Format("Jan 02, 2006")
		}
	}
	fmt.Printf("  Last Active:     %s\n", lastActiveStr)
}

// DisplayGoalConfig prints the current goal configuration.
func DisplayGoalConfig(config *GoalConfig) {
	fmt.Println("\n  Current Goals")
	fmt.Println("  ════════════════════════════════════════════")

	if !config.Enabled {
		fmt.Println("  No goals set")
		fmt.Println("\n  Set goals with: pomo goals set --daily N --weekly N")
		return
	}

	if config.DailyPomosGoal > 0 {
		fmt.Printf("  Daily Goal:   %d pomodoros\n", config.DailyPomosGoal)
	}

	if config.WeeklyPomosGoal > 0 {
		fmt.Printf("  Weekly Goal:  %d pomodoros\n", config.WeeklyPomosGoal)
	}

	if !config.CreatedAt.IsZero() {
		fmt.Printf("\n  Set on:       %s\n", config.CreatedAt.Format("Jan 02, 2006"))
	}
}
