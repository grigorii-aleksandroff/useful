# Autonomous Coding Agent

Roles: planner, coder, reviewer
All state is stored in files. No hidden memory.

---

## Memory

.claude/roles/{planner|coder|reviewer}/MEMORY.md (append-only)  
Rules: always read MEMORY.md before acting

---

## Task Directory

.claude/tasks/{task_id}/
task.md
plan.md
review.md
state.json
summary.md

---

## State

status: planning | coding | review | done  
current_step: integer  
steps_total: integer  
iteration: integer

---

## New Task Handling

- When a new task is created (`/new-task {task_id}`):
    1. Call the `new-task` skill with `{task_id}`  
       → scaffolds `.claude/tasks/{task_id}/` with all templates
    2. User edits `task.md` to describe the task
    3. User says `"begin"` to start the agent workflow

---

## Flow

planning → coding → review → (loop until PASS) → done

---

## Actions by Status

### planning → planner
- Read `task.md` + `review.md` (if exists)
- Generate `plan.md` covering:
    - all task requirements
    - all review issues
- Update `state.json`:
    - `steps_total` = number of steps in plan
    - `current_step = 1`
    - `status = "coding"`
    - `iteration++`

### coding → coder
- Execute **current_step**
- Update `state.json`:
    - if `current_step < steps_total` → increment `current_step`
    - else → `status = "review"` and commit all changes

### review → reviewer
- Read `plan.md` + `git diff`
- Write `review.md` with Verdict + Issues
- If `Verdict = PASS` → `status = "done"`
- If `Verdict = FAIL` → `status = "planning"`

### done
- Write `summary.md`
- Stop execution

---

## Core Rules

- One step per run
- Always follow `MEMORY.md`
- No hidden state
- Validate consistency between `state.json` and actual progress
- **On every status transition**, automatically call the `workflow-event-handler` skill with:
    - Event: `"<previous_status>→<next_status>"`
    - Task ID: `$TASK_ID` (from `state.json`)