FROM golang:1.12.2 as builder
LABEL maintainer="Henrique Vicente <henriquevicente@gmail.com>"

COPY . /go/src/github.com/henvic/galaxy

# disable CGO so we can use multi-stage with alpine. Otherwise, this error happens:
# standard_init_linux.go:207: exec user process caused "no such file or directory"
ENV CGO_ENABLED="0"

RUN [ "go", "build", "-o", "/bin/galaxy", "/go/src/github.com/henvic/galaxy/cmd/server" ]

FROM alpine

COPY --from=builder /bin/galaxy /bin
RUN [ "chmod", "+x", "/bin/galaxy" ]

EXPOSE 9000
ENTRYPOINT [ "/bin/galaxy", "-addr=0.0.0.0:9000" ]
