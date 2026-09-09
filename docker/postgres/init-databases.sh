#!/bin/bash
set -e

# Create the three databases used by the telemetry microservices.
# This script runs once on the first start of the postgres container
# (mounted at /docker-entrypoint-initdb.d/).
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
  CREATE DATABASE telemetry_vehicle;
  CREATE DATABASE geo_service;
  CREATE DATABASE telemetry_alert;
EOSQL
