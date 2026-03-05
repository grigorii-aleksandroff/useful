#!/usr/bin/env bash
set -e

DUMP_BIN="$1"
BACKUP_DIR="$2"
DB_HOST="$3"
DB_PORT="$4"
DB_NAME="$5"
DB_USER="$6"
DB_PASS="$7"
KEEP="${8:-14}"

if [ -z "$DUMP_BIN" ] || [ -z "$BACKUP_DIR" ] || [ -z "$DB_PASS" ]; then
  echo "Usage: backup_db.sh <dump_bin> <backup_dir> <host> <port> <db> <user> <pass> [keep]"
  exit 1
fi

mkdir -p "$BACKUP_DIR"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
FILE="$BACKUP_DIR/dump_${DB_NAME}_${TIMESTAMP}.sql.gz"

echo "Starting DB backup: $DB_NAME"
echo "File: $FILE"

"$DUMP_BIN" -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" \
  --single-transaction --quick --triggers --no-tablespaces \
  "$DB_NAME" | gzip > "$FILE"

echo "Dump finished: $FILE"

OLD_FILES=$(find "$BACKUP_DIR" -name "dump_${DB_NAME}_*.sql.gz" -type f | sort -r | tail -n +"$((KEEP + 1))")
for old in $OLD_FILES; do
    rm -f "$old"
    echo "Deleted old backup: $old"
done

echo "Backup completed"