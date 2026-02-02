---
layout: center
---

## Exercise: Implement the Price Filter Feature using **planning**

---

### Goal

Practice the **plan-then-act** workflow by implementing a new feature in the Lebonpoint platform. You will ensure that the implementation is consistent with the existing architecture and API contract by first creating and reviewing a detailed plan.

---

## Doc: Openspec

1. `npm install -g openspec`
2. `openspec init`
3. Launch Claude
4. Choose between:  

```
/openspec:proposal - Create implementation plan
/openspec:apply - Execute plan
```

---

## Doc: Superpowers
1. Install Skill from Skills.sh
    ```bash
    npx skills add https://github.com/obra/superpowers --skill brainstorming \
        writing-plans \
        executing-plans \
        test-driven-development \
        finishing-a-development-branch
    ```
2. Launch Claude
3. Find the skills using `/<skill-name>` commands.

---

### Hints

- Optionally brainstorm, write and review the implementation plan. Provide the API documentation to the brainstorming or planning tool to ensure alignment with the existing contract.
- Examine `[openapi.yaml](/server/api/openapi.yaml)` to understand the contract all servers must follow.
