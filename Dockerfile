FROM scratch


COPY ./bin/redigo /app

ENTRYPOINT ["/app"]