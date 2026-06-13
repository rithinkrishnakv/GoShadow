package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func main() {
	key := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	
	crypt := func(s string) string {
		b := []byte(s)
		for i := range b {
			b[i] ^= key[i%len(key)]
		}
		return base64.StdEncoding.EncodeToString(b)
	}
	
	// Get credentials from args or env for automation
	var token, ownerID, chanID string
	
	if len(os.Args) == 4 {
		// Command line arguments
		token = os.Args[1]
		ownerID = os.Args[2]
		chanID = os.Args[3]
	} else {
		// Interactive mode
		fmt.Print("Enter Bot Token: ")
		fmt.Scanln(&token)
		fmt.Print("Enter Owner ID: ")
		fmt.Scanln(&ownerID)
		fmt.Print("Enter Channel ID: ")
		fmt.Scanln(&chanID)
	}
	
	// Basic validation
	if token == "" || ownerID == "" || chanID == "" {
		fmt.Println("\n[!] Error: Missing required fields")
		os.Exit(1)
	}
	
	// Clean output for easy copy-paste
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("ENCRYPTED CONFIGURATION - COPY BELOW")
	fmt.Println(strings.Repeat("=", 60))
	
	fmt.Println("\n--- COPY START ---")
	fmt.Printf("const (\n")
	fmt.Printf("\tencToken   = \"%s\"\n", crypt(token))
	fmt.Printf("\tencOwnerID = \"%s\"\n", crypt(ownerID))
	fmt.Printf("\tencChanID  = \"%s\"\n", crypt(chanID))
	fmt.Printf(")\n")
	fmt.Println("--- COPY END ---")
	
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println(" Configuration ready! Paste into your agent's const block.")
	fmt.Println(strings.Repeat("=", 60))
	
	// Optional: Verify by showing first few chars
	fmt.Printf("\n[Verification] Token starts with: %s...\n", token[:min(15, len(token))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
