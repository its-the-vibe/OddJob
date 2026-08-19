FROM --platform=$BUILDPLATFORM golang:1.26.6-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags='-s -w' -o /oddjob ./cmd/oddjob

FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=builder /oddjob /oddjob

USER nonroot:nonroot

ENTRYPOINT ["/oddjob"]
