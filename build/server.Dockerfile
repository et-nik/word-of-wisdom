FROM golang:alpine as gobuilder

COPY --chown=1000:1000 .. /app

WORKDIR /app/cmd/server

RUN go build -o server


FROM scratch

COPY --from=gobuilder /app/cmd/server/server /app/run

WORKDIR /app

EXPOSE 9100

ENTRYPOINT ["/app/run"]