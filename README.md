# Alya-Go 🤖

A WhatsApp bot built with Golang using the [whatsmeow](https://github.com/tulir/whatsmeow) library. Features include QR-based authentication, session management, and basic message handling.

![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## Features ✨
- QR code login with session persistence
- Auto-reply to messages
- Multi-device support (no phone required after login)
- Cross-platform (Termux, Linux, Windows, Mac, VPS)

---

## Installation 💻

### **Prerequisites**
- Go 1.20+
- SQLite3
- GCC (for CGO)

---

### **Termux (Android)**
```bash
# Install dependencies
pkg install golang git sqlite clang

# Clone & run
git clone https://github.com/FahriAdison/Alya-Go.git
cd Alya-Go
CGO_ENABLED=1 go run main.go
```

### **Linux/VPS**
```bash
# Debian/Ubuntu
sudo apt-get update
sudo apt-get install golang git sqlite3 gcc -y

# Arch/Manjaro
sudo pacman -S go git sqlite gcc

# CentOS/RHEL
sudo yum install golang git sqlite gcc

# Run (all distros)
git clone https://github.com/FahriAdison/Alya-Go.git
cd Alya-Go
CGO_ENABLED=1 go run main.go
```

### **Windows**
1. **Option 1: WSL (Recommended)**
   ```powershell
   wsl --install # Install WSL
   wsl # Enter Linux environment
   # Follow Linux instructions above
   ```

2. **Option 2: Native (MinGW)**
   ```powershell
   choco install mingw go sqlite # Requires Chocolatey
   git clone https://github.com/FahriAdison/Alya-Go.git
   cd Alya-Go
   set CGO_ENABLED=1
   go run main.go
   ```

### **MacOS**
```bash
# Install dependencies
brew install go sqlite3

# Run
git clone https://github.com/FahriAdison/Alya-Go.git
cd Alya-Go
CGO_ENABLED=1 go run main.go
```

### **Pterodactyl Panel**
1. Create new "Application" server
2. Install dependencies in startup command:
   ```bash
   apt-get update && apt-get install -y golang git sqlite3 gcc
   ```
3. Add environment variable:
   ```
   Key: CGO_ENABLED
   Value: 1
   ```
4. Use this as startup command:
   ```bash
   git clone https://github.com/FahriAdison/Alya-Go.git && cd Alya-Go && go run main.go
   ```

---

## Usage 🚀
1. Run the bot:
   ```bash
   CGO_ENABLED=1 go run main.go
   ```
2. Scan the QR code via **WhatsApp → Linked Devices**
3. Session will save to `whatsapp-session.db`

---
