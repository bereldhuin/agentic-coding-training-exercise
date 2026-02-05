---
layout: center
---

## Exercise: Implement the Price Filter Feature using **planning**

---

### Goal

Practice the **plan-then-act** workflow by implementing a new feature in the Lebonpoint platform. You will ensure that the implementation is consistent with the existing architecture and API contract by first creating and reviewing a detailed plan.

---

## Doc: Openspec

0. Make sure the `.claude` folder exists in your project root
1. `npm install -g @fission-ai/openspec@latest`
2. `openspec init`
3. Launch Claude
4. Choose between:  

```
1. /opsx:new - Create implementation plan
2. /opsx:ff <proposal-id> - Create the supporting documents
3. /opsx:apply <proposal-id> - Execute plan
```

---

## Doc: Superpowers

0. Make sure the `.claude` folder exists in your project root
1. Install Skill from Skills.sh
    ```bash
    npx skills add https://github.com/obra/superpowers --skill brainstorming
    ```
2. Restart claude
3. Use the brainstorming skill to brainstorm on the price filter feature:
    ```
    /brainstorming
    ```

---

### Hints

- Optionally brainstorm, write and review the implementation plan. Provide the API documentation to the brainstorming or planning tool to ensure alignment with the existing contract
- Review the plan _carefully_ before proceeding to implementation: **some things will most probably need adjustment or may be missing**
  - ❓ Which ones?
