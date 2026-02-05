## Exercise: Generate 3 Production-Ready Skills Using Your Skill Generator Command

### Goal

Use your `/new-skill` (or equivalent) command to generate **three reusable Claude skills** that will help you ship Lebonpoint features with higher reliability:

1. **Commit Style Enforcer** (enforce commit message conventions)
2. **Project Consistency Verifier** (style/architecture/conventions)

---

### Prerequisites

* You have already created the command from the previous exercise:

  * a slash command that scaffolds a new skill given `[skill-name] [description]`

---

## Skill 1: Commit Style Enforcer

Your **Commit Style Enforcer** skill must:
* Inspect the staged changes
* Generate a commit message that adheres to the project's commit message conventions, including:
  * tags within square brackets indicating affected components or areas of the codebase (e.g., `[api]`, `[ui]`, `[docs]`)
  * A concise summary line (under 50 characters)
  * A fairly detailed description of the changes

---

## Skill 2: Project Consistency Verifier

Your **Consistency Verifier** skill must:

* Inspect the change (diff/files touched)
* Compare against repo conventions:
  * naming patterns
  * architecture conventions
  * error-handling conventions

---

## Output format requirements 

Each generated skill should include:

- **Frontmatter with clear name and description**
- Instructions
- Best practices
- Constraints
- Acceptance criteria
