#!/usr/bin/env python3
# goshadow-crypt.py - Encryption/Decryption Tool for GoShadow C2

import base64
import sys
import argparse
from datetime import datetime

VERSION = "1.0"
KEY = b'\xDE\xAD\xBE\xEF'

# Logo
LOGO = """
    ╔═══════════════════════════════════════╗
    ║         GoShadow Crypt v1.0           ║
    ║      Discord C2 Encryption Tool       ║
    ╚═══════════════════════════════════════╝
"""

def xor_cipher_bytes(data_bytes, key):
    """XOR cipher that works on bytes directly"""
    return bytes([b ^ key[i % len(key)] for i, b in enumerate(data_bytes)])

def encrypt(plaintext):
    """Encrypt a string to base64"""
    encrypted_bytes = xor_cipher_bytes(plaintext.encode(), KEY)
    return base64.b64encode(encrypted_bytes).decode()

def decrypt(ciphertext_b64):
    """Decrypt a base64 string to plaintext"""
    xor_bytes = base64.b64decode(ciphertext_b64)
    original_bytes = xor_cipher_bytes(xor_bytes, KEY)
    return original_bytes.decode('utf-8')

def interactive_mode():
    """Interactive menu mode"""
    while True:
        print(LOGO)
        print("1. Encrypt a Message (Plaintext → Discord)")
        print("2. Decrypt a Message (Discord → Plaintext)")
        print("3. Batch Encrypt (one per line)")
        print("4. Exit")
        
        choice = input("\nSelect option (1-4): ")

        if choice == '1':
            msg = input("Enter plaintext to encrypt: ")
            if msg:
                result = encrypt(msg)
                print(f"\n[+] Encrypted Result:")
                print(f"{result}\n")
            else:
                print("\n[!] Empty message\n")
        
        elif choice == '2':
            msg = input("Enter base64 to decrypt: ")
            if msg:
                try:
                    result = decrypt(msg)
                    print(f"\n[+] Decrypted Result:")
                    print(f"{result}\n")
                except Exception as e:
                    print(f"\n[!] Decryption failed: {e}\n")
            else:
                print("\n[!] Empty input\n")
        
        elif choice == '3':
            print("\n Batch Mode - Enter messages (one per line, empty line to finish):")
            messages = []
            while True:
                line = input()
                if not line:
                    break
                messages.append(line)
            
            if messages:
                print("\n" + "="*60)
                print("ENCRYPTED RESULTS")
                print("="*60)
                for msg in messages:
                    encrypted = encrypt(msg)
                    print(f"{msg} → {encrypted}")
                print("="*60)
            else:
                print("\n[!] No messages entered\n")
            input("\nPress Enter to continue...")
        
        elif choice == '4':
            print("\nGoodbye from GoShadow! Stay in the shadows...")
            break
        
        else:
            print("\n[!] Invalid selection. Try again.\n")
            input("Press Enter to continue...")

def main():
    parser = argparse.ArgumentParser(
        description='GoShadow Crypt - Encryption/Decryption Tool for GoShadow C2',
        epilog='Examples:\n  %(prog)s -e "ping"\n  %(prog)s -d "l+P4k+vM2N+i2teBusLJnKLM04vomQ=="'
    )
    parser.add_argument('-e', '--encrypt', help='Encrypt a message')
    parser.add_argument('-d', '--decrypt', help='Decrypt a base64 string')
    parser.add_argument('-v', '--version', action='store_true', help='Show version')
    parser.add_argument('-i', '--interactive', action='store_true', help='Start interactive mode')
    
    args = parser.parse_args()
    
    # Show version
    if args.version:
        print(f"GoShadow Crypt v{VERSION}")
        print("Silent Execution, Instant Results")
        return
    
    # Command line encryption
    if args.encrypt:
        result = encrypt(args.encrypt)
        print(result)
        return
    
    # Command line decryption
    if args.decrypt:
        try:
            result = decrypt(args.decrypt)
            print(result)
        except Exception as e:
            print(f"Error: {e}", file=sys.stderr)
            sys.exit(1)
        return
    
    # Interactive mode (default if no args)
    if args.interactive or len(sys.argv) == 1:
        interactive_mode()
    else:
        parser.print_help()

if __name__ == "__main__":
    main()