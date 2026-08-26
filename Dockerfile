# https://hub.docker.com/_/golang/tags
FROM golang:1.26.5@sha256:705e964a93a2fd2e75c7d59bb7d781b57e30f12293ffde5175c69229e18fb678 AS build

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