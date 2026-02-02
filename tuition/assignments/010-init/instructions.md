---
layout: center
---

## Exercise: Explore the Lebonpoint Multi-Stack Platform and Implmement a new Feature

---


---

### Goal

Learn how to configure your project for agentic coding and write a new feature.

---

### Context

Lebonpoint is a proof-of-concept application demonstrating a unified items API implemented in five different programming languages: TypeScript, Python, Go, Kotlin, and Swift.

The project emphasizes API consistency, where all implementations share the same OpenAPI specification and a single SQLite database, verified by an automated test suite.

---

### Task

Explore the Lebonpoint architecture and implement the search by price feature on your chosen backend. The feature should allow users to filter items based on a minimum and maximum price range. Refer to the [openapi.yaml](/server/api/openapi.yaml) file for the exact schema and endpoint details.

#### Step 1: Introducing Lebonpoint

Familiarize yourself with the Lebonpoint platform by reviewing the different backend implementations in the `server/` directory:

* **TypeScript**: Node.js + Express
* **Python**: FastAPI
* **Go**: Gin
* **Kotlin**: Ktor
* **Swift**: Vapor

---

#### Step 2: Claude Init

Initialize Claude in the project directory to enable AI-assisted development:

```bash
claude
/init
```

This will allow Claude to index the multi-language structure and understand how the different components interact.

Understand the project structure and the role of the different tiers:
- backend
- database
- tests
- client

---

#### Step 3: Generate Price Filter Feature

- Use Claude, Codex or another AI assistant to help you implement the price filter feature in your chosen backend(s).
- Ensure that the new feature adheres to the existing API contract defined in `openapi.yaml`.

---

### Hints

- Examine `[openapi.yaml](/server/api/openapi.yaml)` to understand the contract all servers must follow.
