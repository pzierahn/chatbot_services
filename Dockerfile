FROM golang:latest AS builder

WORKDIR /install

COPY . /install
RUN rm -rf go.sum; \
    go get all
RUN go install cmd/server/server.go

FROM ubuntu:latest

RUN apt-get update;  \
    apt-get upgrade -y; \
    apt-get dist-upgrade -y; \
    apt-get install -y xpdf; \
    apt-get autoremove -y; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/ /bin/
CMD ["server"]