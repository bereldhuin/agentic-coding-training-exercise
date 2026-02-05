---
layout: center
---

## Exercise: Explore the Lebonpoint Multi-Stack Platform and Implmement a new Feature

---

### Goal

Learn how to configure your project for agentic coding and write a new feature.

---

### Context

Lebonpoint is a proof-of-concept application demonstrating a unified items API implemented in five different programming languages: TypeScript, Python, Go, Kotlin, and Swift.

The project emphasizes API consistency, where all implementations share the same OpenAPI specification and a single SQLite database, verified by an automated test suite.

---

### Task

Explore the Lebonpoint architecture and implement the search by price feature on your chosen backend. The feature should allow users to filter items based on a minimum and maximum price range.

#### Step 1: Introducing Lebonpoint

Familiarize yourself with the Lebonpoint platform by reviewing the different backend implementations in the `server/` directory:

* **TypeScript**: Node.js + Express
* **Python**: FastAPI
* **Go**: Gin
* **Kotlin**: Ktor
* **Swift**: Vapor

---

#### Step 2: Claude.md / Agents.md Initialization

1. Write a comprehensive `CLAUDE.md` or `AGENTS.md` file that describes the overall architecture of the Lebonpoint project.

2. Add _another_ `CLAUDE.md` or `AGENTS.md` file in the `server/<your-stack>` directory that details the specific implementation of the backend you chose (e.g., TypeScript, Python, Go, Kotlin, Swift).

---

#### Step 3: Generate Price Filter Feature

- Use Claude, Codex or another AI assistant to help you implement the price filter feature in your chosen backend(s).
