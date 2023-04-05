FROM golang:1.19-alpine AS build

WORKDIR /src/

RUN apk update && apk add git

COPY movieserver.go /src/
COPY movieapi /src/
RUN go mod init src
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/movieserver movieserver.go

FROM scratch
COPY --from=build /bin/movieserver /bin/movieserver
ENTRYPOINT ["/bin/movieserver"]