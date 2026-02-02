## Exercise: Create a Claude Skill Generator Skill

### Goal

Create a **reusable Claude skill generator** that helps you scaffold new skills quickly and consistently, using prompt-engineering best practices.

By the end, you will have a **skill** usable via **slash command** that produces _another skill_ based on a name + description, while first inspecting the current project to stay aligned with local conventions.

---

### Task

Implement a new skill that generates _another Claude skill_.

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
