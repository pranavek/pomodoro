package pomo

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// SessionRecord represents a completed pomodoro session with all its statistics.
type SessionRecord struct {
	ID              int           `json:"id"`
	Date            time.Time     `json:"date"`
	Title           string        `json:"title,omitempty"`
	Goal            string        `json:"goal,omitempty"`
	CompletedPomos  int           `json:"completed_pomos"`
	SkippedSessions int           `json:"skipped_sessions"`
	WorkTime        time.Duration `json:"work_time"`
	BreakTime       time.Duration `json:"break_time"`
	Duration        time.Duration `json:"duration"`
}

// Storage handles reading and writing session records to SQLite database.
type Storage struct {
	db *sql.DB
}

// NewStorage creates a new Storage instance.
// It stores data in ~/.pomo/pomo.db by default.
func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	pomoDir := filepath.Join(homeDir, ".pomo")
	if err := os.MkdirAll(pomoDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(pomoDir, "pomo.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	storage := &Storage{db: db}
	if err := storage.initDB(); err != nil {
		db.Close()
		return nil, err
	}

	return storage, nil
}

// initDB creates the necessary tables if they don't exist.
func (s *Storage) initDB() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATETIME NOT NULL,
		title TEXT,
		goal TEXT,
		completed_pomos INTEGER NOT NULL,
		skipped_sessions INTEGER NOT NULL,
		work_time INTEGER NOT NULL,
		break_time INTEGER NOT NULL,
		duration INTEGER NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_date ON sessions(date);
	CREATE INDEX IF NOT EXISTS idx_sessions_goal ON sessions(goal);
	`

	_, err := s.db.Exec(schema)
	return err
}

// Close closes the database connection.
func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// LoadRecords loads all session records from the database.
func (s *Storage) LoadRecords() ([]SessionRecord, error) {
	query := `
		SELECT id, date, title, goal, completed_pomos, skipped_sessions,
		       work_time, break_time, duration
		FROM sessions
		ORDER BY date DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []SessionRecord
	for rows.Next() {
		var r SessionRecord
		var dateStr string
		var workTimeNanos, breakTimeNanos, durationNanos int64

		err := rows.Scan(
			&r.ID,
			&dateStr,
			&r.Title,
			&r.Goal,
			&r.CompletedPomos,
			&r.SkippedSessions,
			&workTimeNanos,
			&breakTimeNanos,
			&durationNanos,
		)
		if err != nil {
			return nil, err
		}

		r.Date, _ = time.Parse(time.RFC3339, dateStr)
		r.WorkTime = time.Duration(workTimeNanos)
		r.BreakTime = time.Duration(breakTimeNanos)
		r.Duration = time.Duration(durationNanos)

		records = append(records, r)
	}

	return records, rows.Err()
}

// SaveRecord inserts a new session record into the database.
func (s *Storage) SaveRecord(record SessionRecord) error {
	query := `
		INSERT INTO sessions (date, title, goal, completed_pomos, skipped_sessions,
		                     work_time, break_time, duration)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		record.Date.Format(time.RFC3339),
		record.Title,
		record.Goal,
		record.CompletedPomos,
		record.SkippedSessions,
		int64(record.WorkTime),
		int64(record.BreakTime),
		int64(record.Duration),
	)

	return err
}

// GetRecordsSince returns all records since the given date.
func (s *Storage) GetRecordsSince(since time.Time) ([]SessionRecord, error) {
	query := `
		SELECT id, date, title, goal, completed_pomos, skipped_sessions,
		       work_time, break_time, duration
		FROM sessions
		WHERE date >= ?
		ORDER BY date DESC
	`

	rows, err := s.db.Query(query, since.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []SessionRecord
	for rows.Next() {
		var r SessionRecord
		var dateStr string
		var workTimeNanos, breakTimeNanos, durationNanos int64

		err := rows.Scan(
			&r.ID,
			&dateStr,
			&r.Title,
			&r.Goal,
			&r.CompletedPomos,
			&r.SkippedSessions,
			&workTimeNanos,
			&breakTimeNanos,
			&durationNanos,
		)
		if err != nil {
			return nil, err
		}

		r.Date, _ = time.Parse(time.RFC3339, dateStr)
		r.WorkTime = time.Duration(workTimeNanos)
		r.BreakTime = time.Duration(breakTimeNanos)
		r.Duration = time.Duration(durationNanos)

		records = append(records, r)
	}

	return records, rows.Err()
}

// GetRecordsInRange returns all records within the given date range.
func (s *Storage) GetRecordsInRange(start, end time.Time) ([]SessionRecord, error) {
	query := `
		SELECT id, date, title, goal, completed_pomos, skipped_sessions,
		       work_time, break_time, duration
		FROM sessions
		WHERE date >= ? AND date <= ?
		ORDER BY date DESC
	`

	rows, err := s.db.Query(
		query,
		start.Format(time.RFC3339),
		end.Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []SessionRecord
	for rows.Next() {
		var r SessionRecord
		var dateStr string
		var workTimeNanos, breakTimeNanos, durationNanos int64

		err := rows.Scan(
			&r.ID,
			&dateStr,
			&r.Title,
			&r.Goal,
			&r.CompletedPomos,
			&r.SkippedSessions,
			&workTimeNanos,
			&breakTimeNanos,
			&durationNanos,
		)
		if err != nil {
			return nil, err
		}

		r.Date, _ = time.Parse(time.RFC3339, dateStr)
		r.WorkTime = time.Duration(workTimeNanos)
		r.BreakTime = time.Duration(breakTimeNanos)
		r.Duration = time.Duration(durationNanos)

		records = append(records, r)
	}

	return records, rows.Err()
}
