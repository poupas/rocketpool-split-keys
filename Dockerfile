FROM golang
ADD . /go/src/github.com/poupas/rocketpool-split-keys
RUN go install github.com/poupas/rocketpool-split-keys@latest
ENTRYPOINT /go/bin/rocketpool-split-keys
