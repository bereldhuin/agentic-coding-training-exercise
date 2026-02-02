## Exercise: Implement Seller Rating Feature on Lebonpoint Items

### Goal

Extend the Lebonpoint platform by implementing a new feature that adds seller rating metadata to marketplace items.

By the end, you will have designed and implemented a feature that allows items to display the reputation of the seller (e.g., a rating from 1 to 5), improving trust and helping buyers make better decisions.

---

### Context

Lebonpoint enables users to browse and buy items from various sellers. As the platform grows, it becomes important to provide feedback on seller reliability.

Currently, items have metadata like `title`, `description`, and `price`. This exercise focuses on extending the Item schema to support **Seller Ratings**:
  - Add a `seller_rating` field (integer, 1-5).
  - Add a `seller_review_count` field (integer).

---

### Task

Implement the seller rating feature for items in your chosen backend implementation.

#### Step 1: Analyze Current Item Schema

Review how items are currently structured and stored:

- Locate the Item DTO/Entity in your backend (e.g., `server/typescript/src/domain/entity/Item.ts`).
- Review the database schema in `server/database/db.sqlite` (or the initialization script).
- Identify where changes need to be made (OpenAPI spec, database, domain models, persistence logic).

---

#### Step 2: Design the `seller_rating` Feature

Plan the implementation considering these requirements:

**Updated Item Schema:**
```json
{
  "title": "iPhone 13",
  "price_cents": 65000,
  "seller_rating": 4,
  "seller_review_count": 12
}
```

---

#### Step 3: Implement Using Claude's Agentic Workflow

Use Claude to **implement, test and verify** the feature:
1. Update the OpenAPI specification.
2. Update the database schema (SQL).
3. Implement the changes in your chosen backend stack.
4. Verify using the automated test suite (you may need to update the tests!).

---

### Hints

- Start by updating the `server/api/openapi.yaml` to ensure the contract is updated first.
- Use Claude's planning capabilities to identify all affected areas across the stack.

---

### Expected Deliverables

1. Updated OpenAPI spec supporting `seller_rating`.
2. Database migration or updated initialization script.
3. Backend implementation passing health checks.
4. Unit or integration tests demonstrating the new fields.