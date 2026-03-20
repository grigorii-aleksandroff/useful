---
name: workflow-event-handler
description: Handle workflow events for a task. Performs git pre-check (creates remote branch if missing) and pushes before review.
argument-hint: <event,task_id>
---

# Аргументы
Event: $ARGUMENTS[0]  # например "planning→coding" или "coding→review"
Task ID: $ARGUMENTS[1]

## Steps

1. Log the received event
    - Run: echo "[Skill] Received event: $ARGUMENTS[0], task_id=$ARGUMENTS[1]"

2. Handle planning→coding
    - If: $ARGUMENTS[0] == "planning→coding"
    - Run: |
      echo "[Skill] Performing git pre-check for task $ARGUMENTS[1]"
      # Проверяем локальную ветку
      if ! git rev-parse --verify $ARGUMENTS[1] >/dev/null 2>&1; then
      git fetch origin
      git checkout -b $ARGUMENTS[1] origin/master
      echo "[Skill] Local branch $ARGUMENTS[1] created from origin/master."
      # Создаем удалённую ветку сразу
      git push -u origin $ARGUMENTS[1]
      echo "[Skill] Remote branch $ARGUMENTS[1] created."
      else
      git checkout $ARGUMENTS[1]
      echo "[Skill] Branch $ARGUMENTS[1] already exists locally, switched."
      # Проверяем есть ли удалённая ветка
      if ! git ls-remote --exit-code --heads origin $ARGUMENTS[1] >/dev/null 2>&1; then
      git push -u origin $ARGUMENTS[1]
      echo "[Skill] Remote branch $ARGUMENTS[1] created."
      fi

3. Handle coding→review
    - If: $ARGUMENTS[0] == "coding→review"
    - Run: |
      echo "[Skill] Committing and pushing changes for task $ARGUMENTS[1]"
      git add .
      git commit -m "Task $ARGUMENTS[1]: completed coding steps before review" || echo "No changes to commit"
      git push -u origin $ARGUMENTS[1]
      echo "[Skill] Changes pushed to remote branch $ARGUMENTS[1]"

4. Handle review→done
    - If: $ARGUMENTS[0] == "review→done"
    - Run: echo "[Skill] Task $ARGUMENTS[1] completed"