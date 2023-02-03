FROM golang:1.19-alpine3.16 as buildStage

WORKDIR /app
COPY . .

RUN go build -o wat .

FROM alpine:latest

WORKDIR /
COPY --from=buildStage /app/wat /wat

EXPOSE 6660

ENTRYPOINT ["/wat"]
