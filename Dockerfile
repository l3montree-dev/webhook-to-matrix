# https://hub.docker.com/_/golang/tags
FROM golang:1.24.4

WORKDIR /go/src/app
COPY . .

RUN go build


# https://console.cloud.google.com/artifacts/docker/distroless/us/gcr.io/static-debian12?inv=1&invt=Ab1PZQ
FROM gcr.io/distroless/static-debian12:nonroot

USER 53111

WORKDIR /

COPY --from=build --chown=53111:53111 /go/src/app/webhook-to-matrix /

EXPOSE 5001

CMD ["/webhook-to-matrix"]