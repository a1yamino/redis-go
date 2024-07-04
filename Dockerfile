FROM ubuntu:22.04-slim

WORKDIR /app
COPY ./bin /app

CMD ["./app/bin/redigo"]