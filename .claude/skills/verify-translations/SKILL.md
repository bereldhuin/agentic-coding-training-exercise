---
name: verify-translations
description: Verify i18n translation completeness across language files and propose missing translations. Use this skill when the user asks to check translations, verify i18n files, find missing translation keys, audit language files, or ensure translation completeness. Also trigger when they mention "traductions manquantes", "vérifier les traductions", or working with i18n/localization.
---

# Translation Verification Skill

This skill helps verify that all translation keys from the reference language (typically English) exist in all other language files, and provides interactive translation suggestions for any missing keys.

## When to Use This Skill

Use this skill when the user:
- Wants to check if all translations are complete
- Needs to find missing translation keys across language files
- Is working on i18n/localization and wants to audit their translation files
- Asks about translation coverage or completeness
- Mentions checking "traductions manquantes" or similar in French

## How It Works

The skill follows this workflow:

1. **Run the verification script** to detect missing translation keys
2. **If all translations are complete**: Report success and exit
3. **If translations are missing**: For each missing key, propose translation suggestions
4. **Get user validation**: Present suggestions interactively and wait for user approval
5. **Apply approved translations**: Update the language files with validated translations

## Step 1: Locate the i18n Directory

First, find where the translation files are stored. Common locations:
- `./client/i18n/`
- `./src/i18n/`
- `./locales/`
- `./translations/`

Use the Glob tool to search for i18n directories or JSON translation files:

```bash
# Look for i18n directories
Glob("**/i18n")

# Or look for translation files
Glob("**/*en.json")
Glob("**/*fr.json")
```

## Step 2: Run the Verification Script

Execute the bundled verification script to check for missing keys. The script is located at:

```
<skill-directory>/scripts/check_translations.py
```

Run it with:

```bash
python3 <skill-directory>/scripts/check_translations.py <i18n-directory> [reference-lang]
```

**Parameters:**
- `<i18n-directory>`: Path to the directory containing translation JSON files
- `[reference-lang]`: Optional. Reference language code (default: `en`)

**Example:**

```bash
python3 ~/.claude/skills/verify-translations/scripts/check_translations.py ./client/i18n en
```

**Script Output:**

If all translations are complete:
```
✓ All translations are complete! All language files have all keys from en.json
```
(Exit code 0)

If translations are missing:
```json
{
  "reference_language": "en",
  "missing_translations": {
    "fr": ["app.subtitle", "messages.errorHttp"],
    "it": ["app.subtitle"],
    "oc": ["app.subtitle"]
  },
  "reference_file": "./client/i18n/en.json"
}
```
(Exit code 1)

## Step 3: Handle Complete Translations

If the script reports all translations are complete (exit code 0), inform the user:

> ✓ Toutes les traductions sont complètes ! Tous les fichiers de langue contiennent toutes les clés de référence.

Then exit. No further action needed.

## Step 4: Handle Missing Translations

If the script reports missing translations (exit code 1), parse the JSON output and proceed to suggest translations.

For each language with missing keys:

1. **Read the reference language file** to get the English values for the missing keys
2. **Read the target language file** to understand existing translation style and context
3. **For each missing key**, propose 2-3 translation suggestions

### Reading Values from Nested JSON

The missing keys are reported in dot notation (e.g., `app.subtitle`). To retrieve the English value:

```python
# Example: get value for "app.subtitle" from en.json
import json

with open('client/i18n/en.json') as f:
    en_data = json.load(f)

keys = "app.subtitle".split('.')
value = en_data
for key in keys:
    value = value[key]

# value now contains: "Search, filter, and explore marketplace items in real time."
```

### Proposing Translation Suggestions

For each missing key, analyze:
- The English text to translate
- The context (what section of the app it belongs to)
- The translation style used in existing translations for that language
- Any interpolation variables (like `{{count}}` or `{{status}}`)

Then propose 2-3 natural translation options, explaining why each option works. Consider:
- Formal vs. informal tone (tu/vous in French)
- Regional variations (if applicable)
- Technical terminology consistency
- String length constraints (if it's a UI label)

**Important**: Preserve any interpolation variables exactly as they appear in the English text (e.g., `{{count}}`, `{{status}}`, `{{name}}`).

### Interactive Validation

Present the suggestions to the user **one key at a time** in this format:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Langue: [language-name] ([code])
Clé manquante: [key-path]
Valeur anglaise: "[english-value]"

Suggestions de traduction:
1. "[suggestion-1]"
   → [brief explanation why this works]

2. "[suggestion-2]"
   → [brief explanation why this works]

3. "[suggestion-3]"
   → [brief explanation why this works]

Quelle traduction souhaitez-vous utiliser ? (1/2/3, ou proposez votre propre traduction)
```

**Wait for the user's response** before moving to the next missing key. The user can:
- Choose 1, 2, or 3
- Provide their own translation
- Skip this key (if they want to handle it manually later)

## Step 5: Apply Approved Translations

Once the user validates a translation:

1. **Read the current target language file**
2. **Navigate to the correct nested location** using the dot-notation key path
3. **Insert the translated value**
4. **Write the updated JSON back to the file** (preserving formatting with 2-space indents)

**Example: Adding a translation**

```python
import json

# Read the file
with open('client/i18n/fr.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

# Navigate and set the value
keys = "app.subtitle".split('.')
target = data
for key in keys[:-1]:
    target = target[key]
target[keys[-1]] = "Recherchez, filtrez et explorez les annonces en temps réel."

# Write back with proper formatting
with open('client/i18n/fr.json', 'w', encoding='utf-8') as f:
    json.dump(data, f, indent=2, ensure_ascii=False)
    f.write('\n')  # Add trailing newline
```

**Important**: Use `ensure_ascii=False` to preserve special characters (accents, etc.).

## Step 6: Verify After Changes

After applying all validated translations, run the verification script again to confirm:

```bash
python3 <skill-directory>/scripts/check_translations.py <i18n-directory>
```

If it reports success, inform the user:

> ✓ Traductions mises à jour avec succès ! Toutes les clés manquantes ont été ajoutées.

## Tips for Quality Translations

When proposing translations:

1. **Maintain consistency**: Check how similar terms are translated elsewhere in the file
2. **Respect tone**: Match the formality level (informal "tu" vs formal "vous")
3. **Keep it natural**: Translations should sound natural to native speakers, not word-for-word literal
4. **Preserve technical terms**: Some technical terms might be kept in English or have standard translations
5. **Consider context**: A button label needs different phrasing than a paragraph of text
6. **Length matters**: UI labels need to be concise

## Example Translation Suggestions

**English**: `"Search, filter, and explore marketplace items in real time."`

**French suggestions**:
1. "Recherchez, filtrez et explorez les annonces en temps réel."
   → Direct translation, formal tone, natural flow
2. "Cherchez, filtrez et parcourez les articles en temps réel."
   → Alternative with "cherchez" (slightly more casual) and "parcourez" (browse)
3. "Explorez les annonces en temps réel grâce aux filtres et à la recherche."
   → Restructured for emphasis on real-time aspect

**Italian suggestions**:
1. "Cerca, filtra ed esplora gli articoli del marketplace in tempo reale."
   → Standard translation maintaining structure
2. "Ricerca, filtraggio ed esplorazione degli annunci in tempo reale."
   → Noun form, more formal
3. "Esplora gli annunci in tempo reale con ricerca e filtri."
   → Restructured, slightly more casual

## Handling Edge Cases

### Variables in translations

If the English text contains variables like `{{count}}` or `{{status}}`, ensure they are preserved exactly:

English: `"Loaded {{count}} items"`
French: `"{{count}} annonces chargées"` ✓ (correct)
French: `"XX annonces chargées"` ✗ (wrong - lost the variable)

### Pluralization

Some languages have complex pluralization rules. The verification script only checks key presence, not pluralization logic. If you notice plural forms in the English text, remind the user they may need to check their i18n library's pluralization syntax.

### Missing entire sections

If an entire section is missing (e.g., all `filters.*` keys), handle them as a batch but still validate each translation individually with the user.

## Error Handling

If you encounter errors:

- **Script fails to run**: Check that Python 3 is available and the script has correct permissions
- **JSON parse errors**: The translation files might have syntax errors - show the error to the user
- **File write errors**: Check file permissions in the i18n directory
- **Cannot find i18n directory**: Ask the user to provide the correct path

## Summary

This skill provides an interactive workflow to:
1. Automatically detect missing translation keys
2. Propose intelligent translation suggestions based on context
3. Get user validation before applying changes
4. Update translation files safely with validated translations
5. Verify completeness after updates

The goal is to make translation maintenance easy and ensure nothing is missed while keeping the user in control of translation quality.
