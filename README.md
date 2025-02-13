# Alya-Go 🤖

A WhatsApp bot built with Golang using the [whatsmeow](https://github.com/tulir/whatsmeow) library. Features include QR-based authentication, session management, and basic message handling.

![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## Features ✨
- QR code login with session persistence
- Auto-reply to messages
- Multi-device support (no phone required after login)
- SQLite session storage

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