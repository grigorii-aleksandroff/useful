# Planner Memory

Main memory file for the Planner role.  
Update after each planning session (patterns, mistakes, decisions).

---

## Planning Principles

- Steps must be **atomic** — one concern per step
- Steps must describe **intent, not implementation**
- Each step must be **clear, actionable, and verifiable**
- Avoid ambiguity — steps should not require interpretation

---

## Plan Structure

- Steps must be written in execution order
- `steps_total` must exactly match the number of steps in `plan.md`
- Do not include optional or “nice-to-have” steps
- Do not merge unrelated concerns into one step

---

## Step Quality Rules

Each step must:

- Represent a **single logical action**
- Be completable in one iteration
- Have a clear completion condition
- Not depend on hidden assumptions

---

## Boundaries

- Do not include implementation details (no code, no low-level instructions)
- Do not reference internal tools or processes
- Do not duplicate steps
- Do not create overly large or vague steps

---

## Consistency Rules

- Use consistent naming across steps
- Reuse terminology already introduced in the plan
- Keep steps similar in size and granularity

---

## Past Decisions

<!-- append entries here -->