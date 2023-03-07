FROM golang:1.15-alpine AS build

WORKDIR /src/

COPY movieserver.go /src/
COPY movieapi /src/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /bin/movieserver


FROM scratch
COPY --from=build /bin/movieserver /bin/movieserver
ENTRYPOINT ["/bin/movieserver"]