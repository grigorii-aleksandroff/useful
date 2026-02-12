#!/bin/bash

DB_NAME="kta_ru_app"
DB_HOST="localhost"
DB_PORT=2000
DB_USER="teleport_rw"

DUMP_FILE="./dump_app_$(date +%Y%m%d).sql"

EXCLUDED_TABLES=("jobs" "exchange_rate" "activity_log" "balance_hold" "plus_request" "mobiles_healthcheck" "sms_hub" "counters" "callbacks" "receipt_validation_results" "scoring_queue" "scheduled_events" "slipok_queue" "sms_events")

RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

if [ ! -f "$DUMP_FILE" ]; then
  touch "$DUMP_FILE"
fi

tables=$(mysql -h $DB_HOST -u $DB_USER -P $DB_PORT -N -B -e "SHOW TABLES IN $DB_NAME;" | grep -Ev '[0-9]+$')

if [ -z "$tables" ]; then
  echo -e "${RED}No tables to dump!${NC}"
  exit 1
fi

existing_tables=$(grep -i '^CREATE TABLE' "$DUMP_FILE" | awk '{print $3}' | tr -d '`')

dump_table() {
  local table=$1
  local no_data=$2
  local retries=0
  local max_retries=3

  while [ $retries -lt $max_retries ]; do
    if [ "$no_data" = true ]; then
      mariadb-dump \
        -h $DB_HOST \
        -u $DB_USER \
        --port $DB_PORT \
        --single-transaction \
        --routines \
        --triggers \
        --no-data \
        --no-tablespaces \
        $DB_NAME $table >> "$DUMP_FILE"
    else
      mariadb-dump \
        -h $DB_HOST \
        -u $DB_USER \
        --port $DB_PORT \
        --single-transaction \
        --routines \
        --triggers \
        --no-tablespaces \
        $DB_NAME $table >> "$DUMP_FILE"
    fi

    if [ $? -eq 0 ]; then
      echo -e "${GREEN}Table $table dumped successfully${NC}"
      break
    else
      retries=$((retries + 1))
      echo -e "${YELLOW}Error dumping table $table, attempt $retries/$max_retries...${NC}"
      sleep 5
    fi
  done

  if [ $retries -eq $max_retries ]; then
    echo -e "${RED}Failed to dump table $table after $max_retries attempts, skipping${NC}"
  fi
}

for t in $tables; do
  if echo "$existing_tables" | grep -qw "$t"; then
    echo -e "${YELLOW}Table $t is already in the dump, skipping${NC}"
    continue
  fi

  if printf '%s\n' "${EXCLUDED_TABLES[@]}" | grep -qw "$t"; then
    echo -e "${YELLOW}Dumping structure only for table $t (no data)${NC}"
    dump_table "$t" true
  else
    echo -e "${YELLOW}Dumping table $t with data...${NC}"
    dump_table "$t" false
  fi
done

echo -e "${GREEN}Dump completed: $DUMP_FILE${NC}"