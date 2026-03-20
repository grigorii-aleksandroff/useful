---
name: new-task
description: Scaffold a new agent task directory in .claude/tasks/. Use when starting a new coding task with the Planner/Coder/Reviewer workflow.
argument-hint: <task_id>
---

Scaffold a new task directory for the autonomous agent workflow.

Task ID: $ARGUMENTS

## Steps

1. Create the directory `.claude/tasks/$ARGUMENTS/`
2. Create `.claude/tasks/$ARGUMENTS/task.md` from the template at `.claude/templates/task.md`, replacing `{TASK_NAME}` and `{TASK_ID}` with `$ARGUMENTS`
3. Create `.claude/tasks/$ARGUMENTS/plan.md` from `.claude/templates/plan.md`, replacing `{TASK_NAME}` with `$ARGUMENTS`
4. Create `.claude/tasks/$ARGUMENTS/state.json` from `.claude/templates/state.json`, replacing `{TASK_ID}` with `$ARGUMENTS`
5. Create `.claude/tasks/$ARGUMENTS/summary.md` from `.claude/templates/summary.md`, replacing `{TASK_NAME}` with `$ARGUMENTS`
6. Create `.claude/tasks/$ARGUMENTS/review.md` from `.claude/templates/review.md`

Then tell the user:
- Task directory created at `.claude/tasks/$ARGUMENTS/`
- Edit `.claude/tasks/$ARGUMENTS/task.md` to describe the task
- Say "begin" when ready to start the agent workflow