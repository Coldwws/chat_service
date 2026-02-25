FROM alpine:3.13

RUN apk add --no-cache bash postgresql-client ca-certificates

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose

RUN chmod +x /bin/goose

WORKDIR /root

COPY postgres/migrations/ /migrations/
COPY migration.sh /migration.sh
RUN chmod +x /migration.sh

ENTRYPOINT [ "bash", "/migration.sh" ]