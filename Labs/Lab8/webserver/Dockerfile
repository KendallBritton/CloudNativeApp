FROM golang:1.19-alpine AS build

RUN apk update && apk add git

WORKDIR /src/
COPY webserver.go /src/
RUN go mod init src
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/webserver

FROM scratch
COPY --from=build /bin/webserver /bin/webserver
ENTRYPOINT ["/bin/webserver"]