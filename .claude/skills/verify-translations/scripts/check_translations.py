#!/usr/bin/env python3
"""
Translation verification script for i18n JSON files.

Compares all translation files against the reference language (English)
to find missing keys or structural mismatches.
"""

import json
import sys
from pathlib import Path
from typing import Dict, List, Tuple, Any


def load_json_file(file_path: Path) -> Dict:
    """Load and parse a JSON file."""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            return json.load(f)
    except Exception as e:
        print(f"Error loading {file_path}: {e}", file=sys.stderr)
        sys.exit(1)


def flatten_dict(d: Dict, parent_key: str = '', sep: str = '.') -> Dict[str, Any]:
    """
    Flatten a nested dictionary into a flat dict with dot-separated keys.

    Example: {"app": {"title": "Hello"}} -> {"app.title": "Hello"}
    """
    items = []
    for k, v in d.items():
        new_key = f"{parent_key}{sep}{k}" if parent_key else k
        if isinstance(v, dict):
            items.extend(flatten_dict(v, new_key, sep=sep).items())
        else:
            items.append((new_key, v))
    return dict(items)


def find_missing_keys(reference: Dict, target: Dict) -> List[str]:
    """Find keys present in reference but missing in target."""
    ref_flat = flatten_dict(reference)
    target_flat = flatten_dict(target)

    missing = []
    for key in ref_flat.keys():
        if key not in target_flat:
            missing.append(key)

    return sorted(missing)


def check_translations(i18n_dir: Path, reference_lang: str = 'en') -> Dict[str, List[str]]:
    """
    Check all translation files against the reference language.

    Returns a dict mapping language code to list of missing keys.
    """
    i18n_path = Path(i18n_dir)

    if not i18n_path.exists():
        print(f"Error: Directory {i18n_dir} does not exist", file=sys.stderr)
        sys.exit(1)

    # Load reference language
    reference_file = i18n_path / f"{reference_lang}.json"
    if not reference_file.exists():
        print(f"Error: Reference file {reference_file} not found", file=sys.stderr)
        sys.exit(1)

    reference_data = load_json_file(reference_file)

    # Find all translation files
    translation_files = list(i18n_path.glob("*.json"))

    # Check each language file
    results = {}
    for trans_file in translation_files:
        lang_code = trans_file.stem

        # Skip the reference language
        if lang_code == reference_lang:
            continue

        trans_data = load_json_file(trans_file)
        missing_keys = find_missing_keys(reference_data, trans_data)

        if missing_keys:
            results[lang_code] = missing_keys

    return results


def get_value_from_nested_dict(data: Dict, key_path: str) -> Any:
    """Get value from nested dict using dot-separated key path."""
    keys = key_path.split('.')
    value = data
    for key in keys:
        if isinstance(value, dict):
            value = value.get(key)
            if value is None:
                return None
        else:
            return None
    return value


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 check_translations.py <i18n_directory> [reference_lang]")
        print("Example: python3 check_translations.py ./client/i18n en")
        sys.exit(1)

    i18n_dir = sys.argv[1]
    reference_lang = sys.argv[2] if len(sys.argv) > 2 else 'en'

    results = check_translations(i18n_dir, reference_lang)

    if not results:
        print(f"✓ All translations are complete! All language files have all keys from {reference_lang}.json")
        sys.exit(0)

    # Output results as JSON for easy parsing
    output = {
        "reference_language": reference_lang,
        "missing_translations": results,
        "reference_file": str(Path(i18n_dir) / f"{reference_lang}.json")
    }

    print(json.dumps(output, indent=2, ensure_ascii=False))

    # Exit with error code to indicate missing translations
    sys.exit(1)


if __name__ == "__main__":
    main()
