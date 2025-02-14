## Changelog 📌

### **v1.0.1 - 2025-05-15**
- **Fixed**
  - Undefined reference errors in `lib/function.go`
- **Changed**
  - Updated dependencies to latest `whatsmeow` version
  - Improved security checks and owner for exec
  - Simple Send Message
- **Added**
  - Status indicator system
  - Cross-platform installation guides
  - Add Function On `lib/function.go`

### **v1.0.0 - 2025-02-13**
- Initial Release
  - QR-based authentication
  - Basic command handlers
  - Session persistence
  - System monitoring via `/ping`

### **INDICATOR STATUS SCRIPT**
```bash
🔴 UNDER FIXING/ERROR
🟢 WORK
⚫ MAINTENANCE
```

---
<div align="center" style="margin-top: 20px;">
  <p>
    STATUS SCRIPT NOW: 🟢
  </p>
</div>

# Alya-Go 🤖
<div align="center">
  <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSk5zVAP3m15vO5A2fxPlFnGSyisCS16qdYrw&usqp=CAU" alt="Banner">
  <br>
  
  ![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue)
  ![License](https://img.shields.io/badge/License-MIT-green)
  ![Maintenance](https://img.shields.io/badge/Maintained-Yes-success)
  ![Stars](https://img.shields.io/github/stars/FahriAdison/Alya-Go?style=social)

</div>

A powerful WhatsApp bot built with Golang using [whatsmeow](https://github.com/tulir/whatsmeow) library. Features QR-based auth, session management, and smart message handling.

---

## Features ✨
- ✅ QR Code Authentication
- 🔄 Session Persistence
- 📱 Multi-Device Support
- 🛡️ Owner-Only Commands
- 💻 Cross-Platform (Termux/Linux/Windows/Mac)
- 📊 System Monitoring
- 🖼️ Image Menu Support

---

## Installation Guide 💻

### **Prerequisites**
- Go 1.20+
- SQLite3
- GCC (for CGO)

<details>
<summary><b>📱 Termux (Android)</b></summary>

```bash
pkg install golang git sqlite clang
git clone https://github.com/FahriAdison/Alya-Go.git
cd Alya-Go
CGO_ENABLED=1 go run main.go
```
</details>

<details>
<summary><b>🐧 Linux/VPS</b></summary>

```bash
# Debian/Ubuntu
sudo apt-get update && sudo apt-get install -y golang git sqlite3 gcc

# Arch/Manjaro
sudo pacman -S go git sqlite gcc

# CentOS/RHEL
sudo yum install golang git sqlite gcc

# Run
git clone https://github.com/FahriAdison/Alya-Go.git
cd Alya-Go
CGO_ENABLED=1 go run main.go
```
</details>

<details>
<summary><b>🪟 Windows</b></summary>

```powershell
# Using WSL (Recommended)
wsl --install
wsl
# Follow Linux instructions

# Native (Requires MinGW)
choco install mingw go sqlite
git clone https://github.com/FahriAdison/Alya-Go.git
cd Alya-Go
set CGO_ENABLED=1
go run main.go
```
</details>

<details>
<summary><b>🍎 MacOS</b></summary>

```bash
brew install go sqlite3
git clone https://github.com/FahriAdison/Alya-Go.git
cd Alya-Go
CGO_ENABLED=1 go run main.go
```
</details>

<details>
<summary><b>🦖 Pterodactyl Panel</b></summary>

1. Create new "Application" server
2. Install dependencies in startup:
```bash
apt-get update && apt-get install -y golang git sqlite3 gcc
```
3. Environment variable:
```ini
CGO_ENABLED=1
```
4. Startup command:
```bash
git clone https://github.com/FahriAdison/Alya-Go.git && cd Alya-Go && go run main.go
```
</details>

---

## Usage 🚀
```bash
CGO_ENABLED=1 go run main.go
```
1. Scan QR code via **WhatsApp → Linked Devices**
2. Session saved to `whatsapp-session.db`
3. Available commands:
   - `$ [command]` - Execute shell command
   - `=> [code]` - Evaluate Go code
   - `/ping` - System diagnostics
   - `/menu` - Show interactive menu

---

## Supported Platforms 🖥️
| Platform       | Status | 
|----------------|--------|
| Termux         | ✅     |
| Linux/VPS      | ✅     |
| Windows        | ✅     |
| MacOS          | ✅     |
| Pterodactyl    | ✅     |

---

## 🤝 Contribution
<div align="center">
  <a href="https://github.com/FahriAdison/Alya-Go/issues/new/choose">
    <img src="https://img.shields.io/badge/🚨-Report_Issue-red?style=for-the-badge">
  </a>
  
  <a href="https://github.com/FahriAdison/Alya-Go/compare">
    <img src="https://img.shields.io/badge/💡-Suggest_Improvement-green?style=for-the-badge">
  </a>
</div>

---

## 📜 License
MIT License - See [LICENSE](LICENSE) for details.

<div align="center" style="margin-top: 40px;">
  <h3 style="color: #FF69B4;">
    Crafted with ❤️ by Papah-Chan
  </h3>
  <p>
    ⭐ Star this repo if you find it awesome!
  </p>
  <p>
    🐞 Found a bug? <a href="https://github.com/FahriAdison/Alya-Go/issues">Let us know!</a>
  </p>
</div>
