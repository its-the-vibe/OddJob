FROM golang:1.24 AS builder
WORKDIR /src

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o /oddjob ./cmd/oddjob

FROM scratch
COPY --from=builder /oddjob /oddjob
USER 65532:65532
ENTRYPOINT ["/oddjob"]
