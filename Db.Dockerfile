FROM postgres

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD qwerty
ENV POSTGRES_DB postgres
COPY ./conf/schema.sql /docker-entrypoint-initdb.d/

EXPOSE 5432