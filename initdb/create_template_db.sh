#!/usr/bin/env bash

set -e

psql -v ON_ERROR_STOP=1 --username="$POSTGRES_USER" --dbname="$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE template_db;
  ALTER DATABASE template_db WITH is_template TRUE;
EOSQL

psql -v ON_ERROR_STOP=1 --username="$POSTGRES_USER" --dbname=template_db < /docker-entrypoint-initdb.d/init.sql
