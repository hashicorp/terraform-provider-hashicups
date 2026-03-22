#!/usr/bin/env python3
"""Generate signing-keys.json for boring-registry from GPG public key."""

import json
import sys

def generate_signing_keys(key_id: str, public_key_file: str, output_file: str) -> None:
    """Generate signing-keys.json from GPG public key."""
    try:
        with open(public_key_file, 'r') as f:
            ascii_armor = f.read()
        
        data = {
            "gpg_public_keys": [
                {
                    "key_id": key_id,
                    "ascii_armor": ascii_armor
                }
            ]
        }
        
        with open(output_file, 'w') as f:
            json.dump(data, f, indent=2)
        
        print(f"Generated {output_file} with key ID: {key_id}")
    except Exception as e:
        print(f"Error generating signing-keys.json: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print(f"Usage: {sys.argv[0]} <key_id> <public_key_file> <output_file>", file=sys.stderr)
        sys.exit(1)
    
    key_id = sys.argv[1]
    public_key_file = sys.argv[2]
    output_file = sys.argv[3]
    
    generate_signing_keys(key_id, public_key_file, output_file)
