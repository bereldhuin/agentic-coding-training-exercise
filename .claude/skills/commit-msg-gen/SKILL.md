---
name: commit-msg-gen
description: Generate structured commit messages from staged files. Use this skill whenever the user asks to commit changes, create a commit, or wants help writing a commit message. Triggers on phrases like "commit this", "create a commit", "commit these changes", "generate a commit message", or any variation of committing work.
allowed-tools: Bash(git diff *)
---

# Commit Message Generator

Generate well-structured commit messages following project conventions by analyzing staged files and organizing changes by functional components.

## When to Use

Trigger this skill whenever the user wants to commit changes, regardless of how they phrase it:
- "commit this"
- "create a commit"
- "commit these changes"
- "make a commit"
- "generate a commit message"

## Workflow

### 1. Analyze Staged Changes

First, check what files are staged for commit:

```bash
git diff --cached --name-only
```

Then examine the actual changes to understand what was modified:

```bash
git diff --cached
```

Read through the changes carefully to understand:
- What functionality was added, modified, or removed
- Which components or areas of the codebase were affected
- The intent behind the changes

### 2. Map Files to Tags

Based on the file paths, determine the appropriate tags:

**Backend/API changes:**
- `server/python/*` → `[api][python]`
- `server/go/*` → `[api][go]`
- `server/kotlin/*` → `[api][kotlin]`
- `server/swift/*` → `[api][swift]`
- `server/typescript/*` → `[api][typescript]`

**Frontend changes:**
- `client/*` → `[ui]`

**Documentation changes:**
- `README.md`, `CLAUDE.md`, `*.md` files → `[docs]`

**Multiple areas:**
If changes span multiple areas, include all relevant tags. For example:
- Changes to both `server/python/` and `client/` → `[api][python][ui]`
- Changes to `server/go/` and `README.md` → `[api][go][docs]`

### 3. Generate Commit Message

Create a commit message with this structure:

```
[tag1][tag2]... Concise title (< 50 characters)

Detailed description organized by functional component:
- Component A: description of changes
- Component B: description of changes
- Component C: description of changes
```

**Title guidelines:**
- Keep it under 50 characters
- Use imperative mood ("Add", "Fix", "Update", not "Added", "Fixed", "Updated")
- Be specific but concise
- Focus on WHAT changed, not HOW

**Description guidelines:**
- Organize by functional components (API, UI, Database, Auth, etc.)
- Each component gets a bullet point with a clear description
- Focus on the functional impact and purpose of changes
- Include relevant details like endpoints added, features modified, bugs fixed

### 4. Present and Confirm

Display the generated commit message to the user in a formatted code block so they can review it.

Then ask if they would like to create the commit with this message. Wait for their confirmation before proceeding.

### 5. Create Commit (if confirmed)

If the user confirms, create the commit using:

```bash
git commit -m "$(cat <<'EOF'
[tag1][tag2] Title

- Component A: changes
- Component B: changes

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

Note: Always include the Co-Authored-By line as shown above.

## Examples

### Example 1: Python API + UI Changes

**Staged files:**
- `server/python/src/main.py`
- `server/python/src/i18n.py`
- `client/index.html`
- `client/i18n/en.json`

**Generated message:**
```
[api][python][ui] Add internationalization system

- API: Implemented language detection endpoint and translation loading
- UI: Added language selector dropdown in header
- Translations: Created English translation base file
```

### Example 2: Documentation Update

**Staged files:**
- `README.md`
- `CLAUDE.md`

**Generated message:**
```
[docs] Update project documentation

- README: Added i18n system setup instructions
- CLAUDE: Documented translation file conventions
```

### Example 3: Multi-language API Changes

**Staged files:**
- `server/python/src/main.py`
- `server/go/main.go`

**Generated message:**
```
[api][python][go] Fix CORS headers

- Python: Added Access-Control-Allow-Origin header to all responses
- Go: Updated middleware to include CORS headers
```

## Important Notes

- Always analyze the actual diff content, not just file names
- Group related changes together under logical component names
- If unsure about a change's purpose, describe what you observe
- Keep the title short and actionable
- Make the description detailed enough to understand the changes without reading the diff
