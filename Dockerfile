FROM golang:1.19

COPY . /src

RUN set -ex \
    && cd /src \
    && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags='-s -w -extldflags "-static"' -o /arm64bins/ingress-merge ./cmd/ingress-merge

FROM scratch
COPY --from=0 /arm64bins/ingress-merge /ingress-merge
ENTRYPOINT ["/ingress-merge"]
CMD ["--logtostderr"]
