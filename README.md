---

```markdown
<div align="center">

# 🎭 GoShadow C2

### A Stealthy Discord-Based Command & Control Framework

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-red?style=for-the-badge)](LICENSE)
[![Python](https://img.shields.io/badge/Python-3.6+-3776AB?style=for-the-badge&logo=python)](https://python.org)

**Silent Execution. Instant Results. Discord C2, Refined.**

</div>

---

## 📋 Table of Contents

- [⚠️ Disclaimer](#️-disclaimer)
- [✨ Features](#-features)
- [🏗️ Architecture](#️-architecture)
- [📁 Project Structure](#-project-structure)
- [🛠️ Prerequisites](#️-prerequisites)
- [⚙️ Installation & Bot Setup](#️-installation--bot-setup)
- [🚀 Quick Start Workflow](#-quick-start-workflow)
- [🔧 Production Compiling & Optimization](#-production-compiling--optimization)
- [📜 License](#-license)

---

## ⚠️ Disclaimer

> **This framework is developed strictly for authorized security testing, red team simulations, and educational research purposes ONLY.**
>
> Unauthorized monitoring or control of remote computing environments without explicit, written permission is strictly prohibited by law. The authors and contributors assume absolute liability for downstream misuse, operational failures, or collateral infrastructure damage caused by this software. Execute responsibly.

---

## ✨ Features

| Capability | Technical Implementation | Tactical Advantage |
| :--- | :--- | :--- |
| **Encrypted Transport** | Symmetric XOR stream cipher encapsulated within Base64 payloads | Bypasses perimeter network string match heuristics and packet inspection |
| **Timeout Guard** | Dynamic `context.WithTimeout` enforcement capped at 60 seconds | Prevents orphaned commands from causing perpetual agent hangs or freezes |
| **Memory Protection** | High-ceiling file inspection prior to buffer read sequences (`10MB Limit`) | Eradicates host crash telemetry triggers or memory spikes during exfiltration |
| **Evasion Mechanics** | Integrated Win32 API runtime checks & subsystem manipulation | Sidesteps casual debug deployment attaching tools and user visual awareness |
| **Resilience Control** | Automated cross-thread `sync.Map` tracking and message chunking | Enforces strict API rate compliance while maintaining high payload integrity |

---

## 🏗️ Architecture

```text
┌─────────────────────────────────────────────────────────────┐
│                    OPERATOR WORKSTATION                     │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐      ┌──────────────┐      ┌───────────┐  │
│  │  builder.go  │ ───▶ │   main.go    │ ───▶ │ goshadow  │  │
│  │ (Config Gen) │      │ (Agent Code) │      │  (.exe)   │  │
│  └──────────────┘      └──────────────┘      └───────────┘  │
│          │                                         │        │
│          ▼                                         │        │
│  ┌──────────────┐      ┌──────────────┐            │        │
│  │goshadow-crypt│ ◀──▶ │   Discord    │ ◀──────────┘        │
│  │ (Python CLI) │      │   Bot API    │     HTTPS (TLS)     │
│  └──────────────┘      └──────────────┘                     │
└───────────────────────────────▲─────────────────────────────┘
                                │
                                │ Outbound Connection via TLS
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                       TARGET MACHINE                        │
├─────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────────┐  │
│  │ goshadow.exe (Persistent Background Process)          │  │
│  │   • Decrypts encrypted credentials at initialization  │  │
│  │   • Polls target Discord channel for inbound text     │  │
│  │   • Executes sub-tasks under a strict 60s timeout     │  │
│  │   • Chunks and encrypts outputs back to the panel     │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

```

---

## 📁 Project Structure

```text
GoShadow/
├── main.go               # Core execution loop, session handling, and command routing
├── config.example.go     # Unpopulated reference placeholder for deployment setups
├── go.mod                # Explicit compilation instructions and version pins
├── go.sum                # cryptographically verified module dependencies tree
└── tools/
    ├── builder.go        # Hardcoded cryptographic seed string injector tool
    └── goshadow-crypt.py # Python handling interface for payload encoding/decryption

```

---

## 🛠️ Prerequisites

### Development Environment

* **Go Compiler Stack:** `Version 1.26+` (Required to support optimized runtime configuration).
* **Python Runtime:** `Version 3.6+` (Required for operator interaction parsing).

### External Dependency Synchronization

Before modifying source layouts, populate the target dependency tree from the workspace root:

```bash
go mod download
go mod tidy

```

---

## ⚙️ Installation & Bot Setup

### 1. Project Initialization

```bash
git clone 
cd GoShadow

```

### 2. Infrastructure Provisioning

1. Authenticate into the **[Discord Developer Portal](https://discord.com/developers/applications)**.
2. Select **New Application**, assign an alias (e.g., `WinUpdateService`), and navigate to the **Bot** submenu.
3. Generate a fresh secret string via **Reset Token** and store it securely.
4. Scroll to the **Privileged Gateway Intents** panel and strictly toggle on:
* ✅ **Message Content Intent**
* ✅ **Server Members Intent**



### 3. Guild Access Authorization

Navigate to **OAuth2 URL Generator**, select the `bot` scope, enable the `Send Messages` and `Read Message History` permissions, and use the compiled URL to authorize the daemon inside a target testing server channel.

---

## 🚀 Quick Start Workflow

### Step 1: Encrypt Infrastructure Coordinates

Run the standalone builder utility from your workspace directory:

```bash
go run tools/builder.go

```

*Provide the exact raw parameters when prompted by the initialization engine.*

The tool will return an encoded static constant string tree layout:

```go
const (
    encToken   = ""  // YOUR ENCRYPTED TOKEN
	encOwnerID = ""  // YOUR ENCRYPTED OWNER ID
	encChanID  = ""  // YOUR ENCRYPTED CHANNEL ID  
)

```

### Step 2: Inject Coordinates into Context

Open `main.go` and map the generated values directly into the configuration block:

```go
// ============================================================
// HARDCODED ENCRYPTED CONFIGURATION
// Generate these using: go run tools/builder.go
// ============================================================
const (
	encToken   = "" // Inject token here
	encOwnerID = ""                 // Inject operator account ID
	encChanID  = ""                 // Inject interaction channel ID
)

```

### Step 3: Compile for Target Subsystem

To target a Windows machine silently while building inside a Linux framework like Kali, enforce the compilation cross-compiler flags:

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o goshadow.exe main.go

```

---

## 🔧 Production Compiling & Optimization

### Binary Footprint Reduction via Symbol Stripping

Standard Go binaries frequently retain explicit metadata maps. To minimize the final signature envelope size, utilize custom flags during the generation stage:

* `-s`: Disables DWARF generation mapping trees completely.
* `-w`: Disables internal symbol table generation tracking vectors.

### Deployment Compilation Matrix

| Target Operating Environment | Target CPU Structure | Command Execution String |
| --- | --- | --- |
| **Windows Background Process** | `x86_64 / 64-Bit` | `GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui -s -w" -o goshadow.exe main.go` |
| **Windows Debug Shell** | `x86_64 / 64-Bit` | `GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o goshadow_debug.exe main.go` |
| **Linux Headless Standalone** | `x86_64 / 64-Bit` | `GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o goshadow_linux main.go` |
| **macOS Native Architecture** | `Apple Silicon` | `GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o goshadow_mac main.go` |

---

## 📜 License

Distributed under the MIT License. See the explicit structure below:

```text
MIT License

Copyright (c) 2026 Rimu (GoShadow Contributors)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```

**Built for specialized security operations research configurations.**