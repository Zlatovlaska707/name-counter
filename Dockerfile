# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/namecounter ./cmd/run

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/namecounter /usr/local/bin/namecounter
ENTRYPOINT ["/usr/local/bin/namecounter"]
CMD ["--help"]
