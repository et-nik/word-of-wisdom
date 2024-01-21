FROM golang:alpine as gobuilder

COPY --chown=1000:1000 .. /app

WORKDIR /app/cmd/client

RUN go build -o client


FROM scratch

COPY --from=gobuilder /app/cmd/client/client /app/run

WORKDIR /app

ENTRYPOINT ["/app/run"]