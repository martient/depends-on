FROM golang:1.23.2-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o depends-on .


FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source=https://github.com/martient/depends-on
LABEL org.opencontainers.image.description="Depends-on"
LABEL org.opencontainers.image.licenses=APACHE2.0

WORKDIR /app
COPY --from=builder /app/depends-on .
USER 65532:65532

ENTRYPOINT ["/app/depends-on"]