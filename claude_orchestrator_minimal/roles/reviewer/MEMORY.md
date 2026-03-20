# Reviewer Memory

Main memory file for the Reviewer role.  
Update after each review (recurring issues, quality standards, decisions).

---

## Review Scope

- Always review changes against:
    - `task.md`
    - `plan.md`
    - `state.json`

- Ensure implementation strictly follows the approved plan

---

## Critical Checks (AUTO REJECT)

Reject immediately if:

- Any plan step was skipped
- Multiple steps were executed in one iteration (batched)
- `state.json` was not updated after a completed step
- `current_step` does not match actual progress
- Work was done outside the current step

---

## Plan Compliance

- Each change must correspond to the current step only
- No extra or “hidden” work beyond the plan
- Steps must be completed in order
- No deviations without explicit plan update

---

## Code Quality

- Changes are minimal and scoped to the step
- No unrelated refactoring
- No dead code or temporary hacks left behind
- Code is consistent with existing project patterns

---

## Consistency Checks

- Naming matches existing conventions
- No duplicated logic introduced
- No breaking changes outside the step scope

---

## Definition of Done

A step is considered complete only if:

- Its intent (from `plan.md`) is fully achieved
- Changes are committed
- `state.json` is updated correctly
- No side effects outside the step scope

---

## Template Rules

- Always preserve all section headers from `review.md` template (Verdict, Issues, Critical, Major, Minor, Missing, Gaps)
- Only fill in content — never remove headers even if sections are empty

---

## Past Decisions

- Always append new iteration review to `review.md` — never remove or rewrite existing iteration content

<!-- append entries here -->
