# 🍅 Pomodoro Timer CLI

A cross-platform command-line implementation of the [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique) to help you stay focused and productive.

## ✨ Features

- ⏱️ **Customizable timers** - Configure work sessions, short breaks, and long breaks
- 🔔 **Desktop notifications** - Get notified when sessions start and end
- ⏭️ **Skip sessions** - Press 's' to skip the current session
- 📊 **Session statistics** - Track completed pomodoros, work time, and break time
- 🎯 **Visual progress** - See your progress through the pomodoro cycle
- 🖥️ **Cross-platform** - Works on Linux, macOS (Apple Silicon), and Windows
- 💻 **Real-time countdown** - See time remaining updated every minute
- 💭 **Reflection prompts** - Thoughtful questions during breaks to help you pause and think

## 📦 Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](../../releases):

### Build from Source

Requires Go 1.20 or higher:

**Linux/macOS:**
```bash
git clone https://github.com/pranavek/pomodoro.git
cd pomodoro
go build -o pomo .
sudo mv pomo /usr/local/bin/
```

**Windows (PowerShell):**
```powershell
git clone https://github.com/pranavek/pomodoro.git
cd pomodoro
go build -o pomo.exe .
# Add to PATH or move to a directory in your PATH
```

## 🚀 Usage

### Basic Usage

Start a pomodoro timer with default settings (25min work, 5min short break, 30min long break):

```bash
pomo
```

### Custom Configuration

```bash
# Custom work duration (45 minutes)
pomo --work 45

# Custom short break (10 minutes)
pomo --short-break 10

# Custom long break (20 minutes)
pomo --long-break 20

# Change number of pomodoros before long break (default: 4)
pomo --count 3

# Disable real-time countdown
pomo --countdown=false
```

### Combined Options

```bash
# 30-minute work sessions with 7-minute breaks
pomo -w 30 -s 7 -l 15 -c 3
```

## ⚙️ Configuration Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--work` | `-w` | 25 | Work session duration in minutes (1-120) |
| `--short-break` | `-s` | 5 | Short break duration in minutes (1-60) |
| `--long-break` | `-l` | 30 | Long break duration in minutes (1-120) |
| `--count` | `-c` | 4 | Number of pomodoros before a long break (1-10) |
| `--countdown` | `-d` | true | Show real-time countdown during sessions |

## 🎮 Interactive Controls

During a session:
- **Press 's' + Enter** - Skip the current session
- **Continue prompt** - After each pomodoro cycle, choose whether to continue (y/n)

## 💭 Reflection Prompts

During breaks, the app presents thoughtful questions to help you pause and reflect:

**Short breaks:**
- What did you accomplish in this session?
- What challenged you most?
- Is your current approach working?
- Are you working on what matters?

**Long breaks:**
- What progress have you made today?
- Are you solving the right problem?
- What assumptions should you question?
- How can you approach this more simply?

## 📊 Session Summary

At the end of your session, you'll see statistics including:
- Total pomodoros completed
- Sessions skipped
- Total work time
- Total break time
- Session duration

## 🛠️ The Pomodoro Technique

The Pomodoro Technique is a time management method:

1. **Work** for 25 minutes (one "pomodoro")
2. **Short break** for 5 minutes
3. After 4 pomodoros, take a **long break** (15-30 minutes)
4. Repeat the cycle

This tool implements this technique with customizable durations to fit your workflow.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
