## Exercise: Create a Claude Skill Generator Skill

### Goal

Create a **reusable Claude skill generator** that helps you scaffold new skills quickly and consistently, using prompt-engineering best practices.

By the end, you will have a (meta) **skill** usable via **slash command**.

When triggered, this command:
	* Inspects your current project to identify local coding patterns and conventions.
	* Generates a new, custom skill tailored to those specific project requirements.

---

### Task

The new skill must:
1. Be invocable via a **slash command** (e.g., `/new-skill`).
2. **Explore the current project structure and conventions** to gather context necessary for generating a skill that fits well within the existing ecosystem.

---

#### Inputs

Your command must accept **two arguments**:

1. `skill-name`
2. `description`

Example invocation:

* `/new-skill code-review "Review a PR"`

---

### Hints 

- Feel free to use https://code.claude.com/docs/en/skills
- Feel free to use tool permission documentation: https://code.claude.com/docs/en/iam#tool-specific-permission-rules

---

# Example

```
.claude/skills/
└── my-skill/
    └── SKILL.md
```
