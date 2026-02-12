package pomo

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

// Config holds the configuration for a Pomodoro timer session.
type Config struct {
	WorkDuration        time.Duration
	ShortBreakDuration  time.Duration
	LongBreakDuration   time.Duration
	PomosUntilLongBreak int
	ShowCountdown       bool
	SessionTitle        string
	SessionGoal         string
}

// DefaultConfig returns the default Pomodoro timer configuration.
func DefaultConfig() Config {
	return Config{
		WorkDuration:        25 * time.Minute,
		ShortBreakDuration:  5 * time.Minute,
		LongBreakDuration:   30 * time.Minute,
		PomosUntilLongBreak: 4,
		ShowCountdown:       true,
	}
}

var shortBreakReflections = []string{
	"What did you accomplish in this session?",
	"What challenged you most?",
	"Is your current approach working?",
	"What will you focus on next?",
	"Are you working on what matters?",
	"What can you simplify?",
	"Do you need to adjust your approach?",
	"What did you learn just now?",
}

var longBreakReflections = []string{
	"What progress have you made today?",
	"Are you solving the right problem?",
	"What assumptions should you question?",
	"What would you do differently?",
	"What's the essential work remaining?",
	"How can you approach this more simply?",
	"What have you learned in this cycle?",
	"Is there a better way?",
}

// getReflection returns a random reflection prompt based on break type.
func getReflection(breakType string) string {
	if breakType == "long" {
		return longBreakReflections[rand.Intn(len(longBreakReflections))]
	}
	return shortBreakReflections[rand.Intn(len(shortBreakReflections))]
}

// SessionStats tracks statistics for the current Pomodoro session.
type SessionStats struct {
	CompletedPomos int
	SkippedSessions int
	TotalWorkTime  time.Duration
	TotalBreakTime time.Duration
	StartTime      time.Time
}

// NewSessionStats creates a new SessionStats with the start time set to now.
func NewSessionStats() *SessionStats {
	return &SessionStats{
		StartTime: time.Now(),
	}
}

// DisplaySummary prints a summary of the session statistics.
func (s *SessionStats) DisplaySummary() {
	elapsed := time.Since(s.StartTime)
	fmt.Printf("\n📊 Session Summary\n")
	fmt.Printf("  Pomodoros completed: %d\n", s.CompletedPomos)
	if s.SkippedSessions > 0 {
		fmt.Printf("  Sessions skipped: %d\n", s.SkippedSessions)
	}
	fmt.Printf("  Total work time: %s\n", formatDuration(s.TotalWorkTime))
	fmt.Printf("  Total break time: %s\n", formatDuration(s.TotalBreakTime))
	fmt.Printf("  Session duration: %s\n", formatDuration(elapsed))
}

// formatDuration formats a duration in a human-readable format.
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// alert sends a system notification or falls back to console output.
func alert(message string) error {
	err := beeep.Alert("Pomodoro", message, "assets/information.png")
	if err != nil {
		// Fall back to console if notifications fail
		fmt.Printf("\n🔔 ALERT: %s\n", message)
		return err
	}
	return nil
}

// promptContinue asks the user if they want to continue with another pomodoro.
func promptContinue() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nContinue with another pomodoro? (y/n): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			// On error, default to quit for safety
			return false
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		}
		if response == "n" || response == "no" {
			return false
		}
		// Invalid input, ask again
		fmt.Println("Please enter 'y' or 'n'")
	}
}

// countdown displays a countdown timer or sleeps silently based on config.
// Returns true if the countdown completed naturally, false if skipped.
func countdown(duration time.Duration, showCountdown bool) bool {
	if !showCountdown {
		// Original behavior: just sleep
		time.Sleep(duration)
		return true
	}

	// Buffered channel for skip signal
	skipChan := make(chan bool, 1)

	// Goroutine to listen for skip command
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if strings.ToLower(strings.TrimSpace(input)) == "s" {
				// Try to send skip signal
				select {
				case skipChan <- true:
					return
				default:
					// Countdown already over
					return
				}
			}
			// If not 's', continue listening
		}
	}()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	remaining := duration
	totalMinutes := int(duration.Minutes())

	fmt.Printf("  Time remaining: %d minutes (Press 's' + Enter to skip)\n", totalMinutes)

	for remaining > 0 {
		select {
		case <-skipChan:
			fmt.Printf("\r  ⏭️  Session skipped!                                           \n")
			return false
		case <-ticker.C:
			remaining -= 1 * time.Minute
			if remaining > 0 {
				mins := int(remaining.Minutes())
				fmt.Printf("\r  Time remaining: %d minutes (Press 's' + Enter to skip)   ", mins)
			}
		}
	}

	fmt.Printf("\r  ✅ Time's up!                                                \n")
	return true
}

// displayProgress shows a visual representation of pomodoro progress.
func displayProgress(current, total int) {
	fmt.Printf("\nProgress: ")
	for i := 1; i <= total; i++ {
		if i <= current {
			fmt.Print("✓ ")
		} else {
			fmt.Print("○ ")
		}
	}
	fmt.Printf("(%d/%d)\n", current, total)
}

// runWorkSession executes a work session with countdown and tracking.
func runWorkSession(config Config, pomoNumber int, stats *SessionStats) bool {
	mins := int(config.WorkDuration.Minutes())
	fmt.Printf("\n🎯 Starting pomodoro #%d (%d minutes)\n", pomoNumber, mins)
	alert("It's time to get into the flow")

	completed := countdown(config.WorkDuration, config.ShowCountdown)

	if completed {
		fmt.Println("  ✓ Work session completed!")
		stats.TotalWorkTime += config.WorkDuration
		return true
	} else {
		stats.SkippedSessions++
		return false
	}
}

// runBreak executes a break session with countdown and tracking.
func runBreak(breakType string, duration time.Duration, showCountdown bool, stats *SessionStats) {
	mins := int(duration.Minutes())
	reflection := getReflection(breakType)

	if breakType == "long" {
		fmt.Printf("\n☕ Take a long break (%d minutes)\n", mins)
		alert(fmt.Sprintf("Take a long break - %d minutes", mins))
	} else {
		fmt.Printf("\n☕ Take a short break (%d minutes)\n", mins)
		alert(fmt.Sprintf("Take a short break - %d minutes", mins))
	}

	fmt.Printf("\n💭 %s\n\n", reflection)

	completed := countdown(duration, showCountdown)

	if completed {
		alert(fmt.Sprintf("%d minute break is over", mins))
		fmt.Println("  ✓ Break completed!")
		stats.TotalBreakTime += duration
	} else {
		fmt.Println("  Break skipped!")
		stats.SkippedSessions++
	}
}

// Run starts the Pomodoro timer with the given configuration.
func Run(config Config) {
	rand.Seed(time.Now().UnixNano())
	stats := NewSessionStats()
	pomoCount := 0
	carryOn := true

	fmt.Println("\n🍅 Pomodoro Timer Started!")
	if config.SessionGoal != "" {
		fmt.Printf("Goal: %s\n", config.SessionGoal)
	}
	if config.SessionTitle != "" {
		fmt.Printf("Session: %s\n", config.SessionTitle)
	}
	fmt.Printf("Configuration: %dm work / %dm short break / %dm long break\n",
		int(config.WorkDuration.Minutes()),
		int(config.ShortBreakDuration.Minutes()),
		int(config.LongBreakDuration.Minutes()))

	for carryOn {
		// Run work session
		if runWorkSession(config, pomoCount+1, stats) {
			pomoCount++
			stats.CompletedPomos++

			fmt.Printf("\n✓ Pomodoro #%d completed!\n", pomoCount)
			displayProgress(pomoCount, config.PomosUntilLongBreak)
		}

		// Determine break type
		if pomoCount >= config.PomosUntilLongBreak {
			runBreak("long", config.LongBreakDuration, config.ShowCountdown, stats)
			pomoCount = 0
			fmt.Println("\n🔄 Starting a new pomodoro cycle!")
		} else {
			runBreak("short", config.ShortBreakDuration, config.ShowCountdown, stats)
		}

		// Ask user to continue
		carryOn = promptContinue()
	}

	stats.DisplaySummary()

	// Save session statistics
	if stats.CompletedPomos > 0 {
		storage, err := NewStorage()
		if err == nil {
			record := SessionRecord{
				Date:            stats.StartTime,
				Title:           config.SessionTitle,
				Goal:            config.SessionGoal,
				CompletedPomos:  stats.CompletedPomos,
				SkippedSessions: stats.SkippedSessions,
				WorkTime:        stats.TotalWorkTime,
				BreakTime:       stats.TotalBreakTime,
				Duration:        time.Since(stats.StartTime),
			}
			if err := storage.SaveRecord(record); err != nil {
				fmt.Printf("\n⚠️  Warning: Could not save session data: %v\n", err)
			} else {
				fmt.Println("✓ Session saved!")
			}
		}
	}

	fmt.Println("\n👋 Good bye!")
}
