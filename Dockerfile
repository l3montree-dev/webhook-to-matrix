# https://hub.docker.com/_/golang/tags
FROM registry.opencode.de/open-code/badgebackend/badge-api/golang:1.24.4-bookworm@sha256:ee7ff13d239350cc9b962c1bf371a60f3c32ee00eaaf0d0f0489713a87e51a67 AS build

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 go build


# https://console.cloud.google.com/artifacts/docker/distroless/us/gcr.io/static-debian12?inv=1&invt=Ab1PZQ
FROM gcr.io/distroless/static-debian12:nonroot

USER 53111

WORKDIR /app

COPY --from=build --chown=53111:53111 /go/src/app/webhook-to-matrix /usr/local/bin/webhook-to-matrix

EXPOSE 5001

CMD ["webhook-to-matrix"]