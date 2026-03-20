# Coder Memory

Main memory file for development. Update after each session (patterns, pitfalls, rules).

---

## Template Rules

- Always preserve all fields from `state.json` template (task_id, current_step, steps_total, status, iteration)
- Only update field values — never rename or remove fields

---

## Workflow Rules

- One step at a time — do not batch steps
- Execute steps sequentially; do not skip
- Minimal changes — only what is required for the current step
- Update `state.json` **immediately** after each step
- Commit and push handled by `workflow-event-handler` skill

---

## state.json (CRITICAL)

- Update `current_step` immediately after completing a step
- Do it as the first action before any further analysis or coding
- Do not proceed to the next step without updating it

---

## Commit Rules

- Do not add:
    - `Co-Authored-By`
    - any AI / Claude references
- Do not commit:
    - `.claude/` directories or their contents
- **Do not create empty commits**