FROM previousnext/golang:1.8 as build
ADD . /go/src/github.com/previousnext/k8s-backup
WORKDIR /go/src/github.com/previousnext/k8s-backup
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/previousnext/k8s-backup/bin/k8s-backup_linux_amd64 /usr/local/bin/k8s-backup

ENTRYPOINT ["k8s-backup"]
