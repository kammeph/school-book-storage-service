##
## BUILD
##
FROM golang:1.18-alpine3.15 AS build

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /school-book-storage-service

##
## Deploy
##
FROM alpine:3.15

RUN adduser -D nonroot

WORKDIR /

COPY --from=build /school-book-storage-service /school-book-storage-service

EXPOSE 9090

USER nonroot

ENTRYPOINT [ "/school-book-storage-service" ]