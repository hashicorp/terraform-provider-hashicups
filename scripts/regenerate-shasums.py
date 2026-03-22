#!/usr/bin/env python3
"""Regenerate SHA256SUMS with correct terraform-provider prefixed filenames."""

import sys
import re

def regenerate_shasums(input_file: str, output_file: str, old_prefix: str, new_prefix: str) -> None:
    """Regenerate SHA256SUMS file with renamed filenames."""
    try:
        with open(input_file, 'r') as f:
            lines = f.readlines()
        
        updated_lines = []
        for line in lines:
            # Match SHA256 hash and filename
            match = re.match(r'^([a-f0-9]{64})\s+(.+)$', line.strip())
            if match:
                hash_val = match.group(1)
                filename = match.group(2)
                # Only update filenames that match the old prefix
                if filename.startswith(old_prefix):
                    new_filename = filename.replace(old_prefix, new_prefix, 1)
                    updated_lines.append(f"{hash_val}  {new_filename}\n")
                else:
                    updated_lines.append(line)
            else:
                # Keep non-matching lines as-is
                if line.strip():
                    updated_lines.append(line)
        
        with open(output_file, 'w') as f:
            f.writelines(updated_lines)
        
        print(f"Regenerated {output_file} with {new_prefix} prefixed filenames")
    except Exception as e:
        print(f"Error regenerating SHA256SUMS: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 5:
        print(f"Usage: {sys.argv[0]} <input_file> <output_file> <old_prefix> <new_prefix>", file=sys.stderr)
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    old_prefix = sys.argv[3]
    new_prefix = sys.argv[4]
    
    regenerate_shasums(input_file, output_file, old_prefix, new_prefix)
