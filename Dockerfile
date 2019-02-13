FROM golang:1.11

COPY . /src

RUN set -ex \
    && cd /src \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w -extldflags "-static"' -o /bin/ingress-merge ./cmd/ingress-merge

FROM scratch
COPY --from=0 /bin/ingress-merge /ingress-merge
ENTRYPOINT ["/ingress-merge"]
CMD ["--logtostderr"]
