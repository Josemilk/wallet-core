#!/usr/bin/env bash
set -euo pipefail

: "${DATABASE_URL:?DATABASE_URL is required}"
: "${BACKUP_FILE:?BACKUP_FILE is required}"

mkdir -p "$(dirname "$BACKUP_FILE")"
pg_dump --format=custom --no-owner --no-privileges "$DATABASE_URL" > "$BACKUP_FILE"

# Integrity check: list archive contents and restore into an isolated database.
pg_restore --list "$BACKUP_FILE" > "${BACKUP_FILE}.toc"
: "${DR_RESTORE_DATABASE_URL:?DR_RESTORE_DATABASE_URL is required}"
dropdb --if-exists --force "$DR_RESTORE_DATABASE_URL" || true
createdb "$DR_RESTORE_DATABASE_URL"
pg_restore --exit-on-error --no-owner --no-privileges --dbname "$DR_RESTORE_DATABASE_URL" "$BACKUP_FILE"

echo "DR backup + restore verification completed successfully."
