## Exercise: Create a Lebonpoint Marketplace CLI - with an Unfamiliar Stack

### Goal

Expand your technical versatility by building a **CLI** that interacts with the **Lebonpoint platform** using a programming language or framework that is outside your comfort zone.

By the end, you will have a functional **CLI** that provides a composable search interface for querying and managing marketplace items, proving that AI agents can help you bridge the gap into unfamiliar technology stacks while building useful tools.

---

### Prerequisites

* You have identified a "target" stack (e.g., Python, TypeScript, Go).
* You have the Lebonpoint backends running.

---

## The "Lebonpoint Marketplace CLI"

Your task is to implement a CLI in your **target stack** that interacts with the Lebonpoint Items API or directly with the SQLite database.

The CLI must provide a composable search DSL that allows users to query the marketplace effectively.

---

## Constraints (1/2)

* **Unfamiliar Stack**: You must NOT use your primary programming language.
* **Idiomatic Code & Best Practices**: You must follow the best practices and idiomatic patterns of your target stack.
* **Automated Testing**: You must include unit tests for your tools and integration logic using the standard testing framework for your chosen stack.
* **Network Connectivity**: The CLI should be able to connect to APIs running on different machines.

---

## Constraints (2/2)

* **Integration**: Your CLI MUST call the existing Lebonpoint APIs (e.g., the TypeScript server on port 3000) rather than duplicating database logic.
* **Native Tooling**: Use the idiomatic package manager and project structure for the chosen stack.

---

## Hints

- Brainstorm the architecture and design before coding.
- Use the **skills** we created in previous exercises to help you verify the quality of your implementation or improve your workflow.
- Use other **skills** to debug your CLI output: `npx skills add https://github.com/viteinfinite/skills --skill request-recorder`
- Use **skills** to **automate testing and debugging**: e.g., tell the agent to run the CLI and send testing requests.

---

## Peer Review & Feedback

- Find a peer who masters the stack you've chosen.
- Show your implementation, discuss any challenges you faced and ask for feedback.

---
layout: center
---

Full specification available at:
<br/>
[hymaia/agentic-coding-training-exercise/issues/1](https://github.com/hymaia/agentic-coding-training-exercise/issues/1)
