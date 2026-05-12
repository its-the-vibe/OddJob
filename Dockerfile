FROM golang:1.26.3 AS builder
WORKDIR /src

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /oddjob ./cmd/oddjob

FROM scratch
COPY --from=builder /oddjob /oddjob
USER 65532:65532
ENTRYPOINT ["/oddjob"]
