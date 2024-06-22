#!/bin/bash
set -e

echo "Running init-db.sh script..."
#psql -U postgres < /docker-entrypoint-initdb.d/cpu_usage.sql
psql -U postgres -d homework -c "\COPY cpu_usage FROM /tmp/psql_data/cpu_usage.csv CSV HEADER"
