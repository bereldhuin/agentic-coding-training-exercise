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

1. Launch Claude
2. `/plugin marketplace add obra/superpowers-marketplace`
3. `/plugin install superpowers@superpowers-marketplace`
4. Restart claude
5. Choose between:  

```
/superpowers:brainstorm - Interactive design refinement
/superpowers:write-plan - Create implementation plan
/superpowers:execute-plan - Execute plan in batches
```

---

### Hints

- Optionally brainstorm, write and review the implementation plan. Provide the API documentation to the brainstorming or planning tool to ensure alignment with the existing contract.
- Examine `[openapi.yaml](/server/api/openapi.yaml)` to understand the contract all servers must follow.
