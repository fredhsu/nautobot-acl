FROM golang:1.16

WORKDIR /go/src/github.com/fredhsu/nautobot-buildacl
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a cmd/getips/getIpAddr.go
RUN ls cmd/getips/
FROM alpine
COPY --from=0 /go/src/github.com/fredhsu/nautobot-buildacl/getIpAddr /getIpAddr
CMD ["/getIpAddr"]