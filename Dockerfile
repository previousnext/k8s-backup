FROM previousnext/golang:1.8 as build
ADD . /go/src/github.com/previousnext/k8s-backup
WORKDIR /go/src/github.com/previousnext/k8s-backup
RUN make build

FROM alpine:3.7
RUN apk -v --update add ca-certificates python py-pip groff less mailcap mariadb-client && \
    pip install --upgrade awscli==1.14.5 s3cmd==2.0.1 python-magic && \
    apk -v --purge del py-pip && \
    rm /var/cache/apk/*
COPY --from=build /go/src/github.com/previousnext/k8s-backup/bin/k8s-backup_linux_amd64 /usr/local/bin/k8s-backup
ENTRYPOINT ["k8s-backup"]
