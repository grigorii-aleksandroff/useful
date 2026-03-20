# Coder Memory

Main memory file for development. Update after each session (patterns, pitfalls, rules).

---

## Git setup (once per repo)

1. `git fetch origin`
2. `git checkout -b {task-id} origin/master` (or reuse existing branch)

---

## Workflow Rules

- One step at a time — do not batch steps
- Execute steps sequentially; do not skip
- Minimal changes — only what is required for the current step
- Update `state.json` **immediately** after each step
- **Commit only after all steps are completed** (`current_step == steps_total`)

---

## state.json (CRITICAL)

- Update `current_step` immediately after completing a step
- Do it as the first action before any further analysis or coding
- Do not proceed to the next step without updating it

---

## .claude Rules (CRITICAL)

- Before working in:
    - `asiatix/`
    - `admin-panel/`
    - `admin-panel-2/`  
      → always read `.claude/CLAUDE.md` inside that directory
- All `CLAUDE.md` files:
    - must be stored only inside `.claude/`
    - must never be created outside `.claude/`

---

## Commit Rules

- Do not add:
    - `Co-Authored-By`
    - any AI / Claude references
- Do not commit:
    - `.claude/` directories or their contents
- **Commit all coder changes only after completing all steps**
- **Do not create empty commits**

---

## Past Decisions

<!-- append entries here -->