# --------- ETAP 1------------------------
FROM scratch AS builder

ARG TARGETARCH

ARG ALPINE_ARCH=${TARGETARCH/amd64/x86_64}
ARG ALPINE_ARCH=${ALPINE_ARCH/arm64/aarch64}

ADD alpine-minirootfs-3.21.3-${ALPINE_ARCH}.tar /

LABEL org.opencontainers.image.authors="Kiryl Shkabara"

RUN apk update && \
    apk add --no-cache git ca-certificates wget tar && \
    wget https://go.dev/dl/go1.26.3.linux-${TARGETARCH}.tar.gz && \
    tar -C /usr/local -xzf go1.26.3.linux-${TARGETARCH}.tar.gz && \
    rm go1.26.3.linux-${TARGETARCH}.tar.gz && \
    rm -rf /var/cache/apk/*

ENV PATH=$PATH:/usr/local/go/bin
WORKDIR /src

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o /zad1 .

# --------- ETAP 2 ------------------------
FROM alpine:3.21 AS prod

LABEL org.opencontainers.image.authors="Kiryl Shkabara"
LABEL org.opencontainers.image.description="Weather application in Go - Secure Inline Cache Mode"
LABEL org.opencontainers.image.version="1.2.0"

WORKDIR /app

COPY --from=builder /zad1 ./zad1

EXPOSE 3000

HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD ["./zad1", "-check-health"]

ENTRYPOINT ["./zad1"]