# Autonomous Coding Agent

Three roles: **planner**, **coder**, **reviewer**.
All state is stored in files. No hidden memory.

---

## Memory

```
.claude/roles/{planner|coder|reviewer}/
  MEMORY.md   (append only)
  *           (read-only docs)
```

Rules:
- Only `MEMORY.md` is writable
- Always read it before acting
- Treat every rule in MEMORY.md as a hard constraint

---

## Task Structure

```
.claude/tasks/{task-id}/
  task.md
  plan.md
  review.md
  state.json
  summary.md
```

---

## state.json

```json
{
  "branch": "feature/{task-id}",
  "current_step": 0,
  "steps_total": 0,
  "status": "planning",
  "iteration": 1
}
```

Flow: `planning → coding → review → planning → ... → done`

---

## On Each Run

1. Read role MEMORY.md (MANDATORY)
2. Read task context:
    - state.json
    - task.md
    - review.md (if exists)
3. Act based on `status`

---

## planning → planner

#### Input

- `task.md`
- `review.md` (priority)

#### Goal

Create a plan that fully covers:
- all task requirements
- all review issues

#### Output → plan.md

```
# Steps
1. ...
2. ...

# Fixes

## Critical
3. Fix [C1]

## Major
4. Fix [M1]

## Missing
5. ...

## Gaps
6. ...
```

#### Rules

- Every issue from `review.md` MUST be included
- If any issue is missing → plan is invalid → regenerate
- Plan can be fully rebuilt (not just appended)
- Do not merge multiple fixes into one step
- Each issue must map to at least one step

#### State

- `steps_total = N`
- `current_step = 1`
- `status = "coding"`
- `iteration++`

---

## coding → coder

#### Setup

```
git fetch origin
git checkout -b {task-id} origin/master
```

(or reuse branch)

#### Execution

- Execute ONLY `current_step`
- Make minimal change
- **Do NOT commit yet** — commit only after completing all steps (`current_step == steps_total`)

#### After

- Update `state.json` immediately after completing the step
- Then:
    - If `current_step < steps_total` → increment
    - else → `status = "review"`
    - If `current_step == steps_total` → commit all changes at once

---

## review → reviewer

#### Input

- `task.md`
- `plan.md`
- `git diff`

#### Output → review.md (STRICT FORMAT)

```
# Review (Iteration N)

## Verdict
PASS | FAIL

## Issues

### Critical
- [C1] ...

### Major
- [M1] ...

### Minor
- [m1] ...

## Missing
- ...

## Gaps
- ...
```

#### Rules

- Verdict is REQUIRED
- Critical/Major MUST be explicit
- Missing & Gaps MUST be filled if they exist
- Validate `state.json` reflects actual progress

#### Decision

- PASS → `status = "done"`
- FAIL → `status = "planning"`

---

## done

Write `summary.md`:

```
# Summary

Done: ...
Iterations: N
Notes: ...
```

Stop execution.

---

## Core Rules

- Follow role MEMORY.md strictly
- One step per run
- Loop until PASS
- No hidden state

## Invalid State

System is invalid if:
- No Verdict in `review.md`
- Any Critical/Major issue not in plan
- Missing/Gaps ignored
- state.json is inconsistent with actual progress
- steps_total does not match plan.md

---

## Start

```
/new-task my-feature
```

Edit `task.md`, then run **begin**.