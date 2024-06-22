# Dockerfile
FROM timescale/timescaledb:latest-pg16

#ENV POSTGRES_PASSWORD=$POSTGRES_PASSWORD

RUN mkdir -p /tmp/psql_data/
COPY cpu_usage.csv /tmp/psql_data/

COPY cpu_usage.sql /docker-entrypoint-initdb.d/
#COPY cpu_usage.csv /docker-entrypoint-initdb.d/
COPY init-db.sh /docker-entrypoint-initdb.d/

RUN chmod +x /docker-entrypoint-initdb.d/init-db.sh

EXPOSE 5432
