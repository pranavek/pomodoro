# Pomodoro

A minimal command-line timer for focused work and thoughtful reflection.

Inspired by Dieter Rams' design philosophy: less, but better.

## Philosophy

This is not just a timer. It's a tool that encourages you to pause, think, and reflect on your work. Each break presents a question designed to help you step back and consider your approach, your progress, and your priorities.

Work happens in the doing. Wisdom happens in the pausing.

## Features

- Minimal, distraction-free interface
- Thoughtful reflection prompts during breaks
- Configurable work and rest periods
- Desktop notifications
- Session tracking
- Cross-platform (Linux, macOS)

## Installation

### Build from Source

Requires Go 1.20 or higher:

```bash
git clone https://github.com/pranavek/pomodoro.git
cd pomodoro
go build -o pomo .
sudo mv pomo /usr/local/bin/
```

Or use the build script:

```bash
./build.sh
```

## Usage

### Start with defaults

25 minutes work, 5 minute pause, 30 minute rest:

```bash
pomo
```

### Customize durations

```bash
# 45 minute work sessions
pomo --work 45

# 10 minute pauses
pomo --short-break 10

# 20 minute rest periods
pomo --long-break 20

# Take rest after 3 work sessions (instead of 4)
pomo --count 3
```

### Combined example

```bash
pomo -w 30 -s 7 -l 15 -c 3
```

## Configuration

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--work` | `-w` | 25 | Work duration (minutes) |
| `--short-break` | `-s` | 5 | Pause duration (minutes) |
| `--long-break` | `-l` | 30 | Rest duration (minutes) |
| `--count` | `-c` | 4 | Work sessions before rest |
| `--countdown` | `-d` | true | Show countdown timer |

## During a Session

- Press `s` then Enter to skip current session
- After each cycle, choose to continue or review your summary
- During breaks, reflect on the question presented

## The Questions

During short pauses, you might be asked:
- What did you accomplish in this session?
- What challenged you most?
- Is your current approach working?
- Are you working on what matters?

During longer rest periods:
- What progress have you made today?
- Are you solving the right problem?
- What assumptions should you question?
- How can you approach this more simply?

These aren't rhetorical. Take a moment. Think.

## Principles

**Less, but better**
No clutter. No distractions. Just time, work, and thought.

**Pause to think**
The breaks aren't just rest. They're for reflection.

**Question your work**
Are you building the right thing? Could it be simpler?

**Respect the process**
Deep work requires time and space. This tool protects both.

## License

MIT License - see [LICENSE](LICENSE) file for details.
