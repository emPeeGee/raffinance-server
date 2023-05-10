FROM postgres:latest
ENV POSTGRES_PASSWORD 1111
ENV POSTGRES_DB raffinance
COPY ./schema/init.sql /docker-entrypoint-initdb.d/
