# 🍅 Pomodoro Timer CLI

A cross-platform command-line implementation of the [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique) to help you stay focused and productive.

## ✨ Features

- ⏱️ **Customizable timers** - Configure work sessions, short breaks, and long breaks
- 🔔 **Desktop notifications** - Get notified when sessions start and end
- ⏭️ **Skip sessions** - Press 's' to skip the current session
- 📊 **Session statistics** - Track completed pomodoros, work time, and break time
- 🎯 **Visual progress** - See your progress through the pomodoro cycle
- 🖥️ **Cross-platform** - Works on Linux and macOS (Intel & Apple Silicon)
- 💻 **Real-time countdown** - See time remaining updated every minute

## 📦 Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](../../releases):

**Linux (amd64):**
```bash
curl -LO https://github.com/pranavek/pomodoro/releases/latest/download/pomodoro-linux-amd64
chmod +x pomodoro-linux-amd64
sudo mv pomodoro-linux-amd64 /usr/local/bin/pomo
```

**Linux (arm64):**
```bash
curl -LO https://github.com/pranavek/pomodoro/releases/latest/download/pomodoro-linux-arm64
chmod +x pomodoro-linux-arm64
sudo mv pomodoro-linux-arm64 /usr/local/bin/pomo
```

**macOS (Intel):**
```bash
curl -LO https://github.com/pranavek/pomodoro/releases/latest/download/pomodoro-darwin-amd64
chmod +x pomodoro-darwin-amd64
sudo mv pomodoro-darwin-amd64 /usr/local/bin/pomo
```

**macOS (Apple Silicon):**
```bash
curl -LO https://github.com/pranavek/pomodoro/releases/latest/download/pomodoro-darwin-arm64
chmod +x pomodoro-darwin-arm64
sudo mv pomodoro-darwin-arm64 /usr/local/bin/pomo
```

### Build from Source

Requires Go 1.20 or higher:

```bash
git clone https://github.com/pranavek/pomodoro.git
cd pomodoro
go build -o pomo .
sudo mv pomo /usr/local/bin/
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
