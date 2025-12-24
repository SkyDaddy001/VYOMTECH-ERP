#!/usr/bin/env python3
"""
Migrate Go models from int/int64 IDs to string ULIDs
Run: python3 migrate_models.py
"""

import re
import os
from pathlib import Path

def convert_model_file(file_path):
    """Convert int IDs to string in a Go model file"""
    with open(file_path, 'r') as f:
        content = f.read()
    
    original = content
    
    # Convert ID field from int/int64 to string
    # Pattern: ID <space> int/int64 <backtick>
    content = re.sub(
        r'(\s+ID\s+)(int|int64)(\s+`json:"id")',
        r'\1string\3',
        content
    )
    
    # Convert UserID from int/int64 to string
    content = re.sub(
        r'(\s+UserID\s+)(int|int64)(\s+`)',
        r'\1string\3',
        content
    )
    
    # Convert *int to *string for optional IDs
    content = re.sub(
        r'(\s+\w+ID\s+)\*?(int|int64)(\s+`)',
        r'\1*string\3',
        content
    )
    
    # Convert CreatedBy from int/int64 to string
    content = re.sub(
        r'(\s+CreatedBy\s+)(int|int64)(\s+`)',
        r'\1*string\3',
        content
    )
    
    # Convert UpdatedBy from int/int64 to string
    content = re.sub(
        r'(\s+UpdatedBy\s+)(int|int64)(\s+`)',
        r'\1*string\3',
        content
    )
    
    if content != original:
        with open(file_path, 'w') as f:
            f.write(content)
        return True
    return False

def main():
    models_dir = Path('internal/models')
    if not models_dir.exists():
        print(f"Models directory not found: {models_dir}")
        return
    
    updated_files = []
    for file_path in models_dir.glob('*.go'):
        if convert_model_file(file_path):
            updated_files.append(file_path.name)
            print(f"✓ Updated {file_path.name}")
        else:
            print(f"- No changes needed in {file_path.name}")
    
    if updated_files:
        print(f"\n✓ Updated {len(updated_files)} files:")
        for f in updated_files:
            print(f"  - {f}")
    else:
        print("\nNo files needed updating.")

if __name__ == '__main__':
    main()
