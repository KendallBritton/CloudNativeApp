FROM golang:1.18-alpine AS build

WORKDIR /src
COPY Webpage ../Webpage
COPY GoCode/Server.go .

RUN go mod init src
RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /bin/server

FROM scratch
COPY --from=build /bin/server /bin/server
COPY Webpage /Webpage
ENTRYPOINT [ "/bin/server" ]