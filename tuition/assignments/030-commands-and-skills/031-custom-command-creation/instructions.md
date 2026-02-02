## Exercise: Create a Claude Skill Generator Command

### Goal

Create a **reusable Claude skill generator** that helps you scaffold new skills quickly and consistently, using prompt-engineering best practices.

By the end, you will have a **slash command** that produces a well-structured skill prompt based on a name + description, while first inspecting the current project to stay aligned with local conventions.

---

### Context

You are working in a repo that uses Claude skills and custom commands. You’ll add a new command whose purpose is:

* take a skill name + a short description
* improve the description into something actionable
* inspect the repo to understand conventions and requirements
* output a new skill prompt aligned with the project

---

### Task

Implement a new custom command that generates a Claude skill.

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
