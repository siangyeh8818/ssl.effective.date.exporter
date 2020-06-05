FROM golang:1.14.4-stretch as builder

COPY go.mod /go/src/github.com/siangyeh8818/ssl.effective.date.exporter/go.mod
COPY go.sum /go/src/github.com/siangyeh8818/ssl.effective.date.exporter/go.sum

# Run golang at any directory, not neccessary $GOROOT, $GOPATH
ENV GO111MODULE=on
WORKDIR /go/src/github.com/siangyeh8818/ssl.effective.date.exporter

# RUN go mod init github.com/pnetwork/sre.monitor.metrics
RUN go mod download
COPY cmd /go/src/github.com/siangyeh8818/ssl.effective.date.exporter/cmd
COPY internal /go/src/github.com/siangyeh8818/ssl.effective.date.exporter/internal
#COPY pkg /go/src/github.com/pnetwork/sre.monitor.metrics/pkg

# Build the Go app
RUN env GOOS=linux GOARCH=amd64 go build -o ssl-effective-date-exporter -v -ldflags "-s" github.com/siangyeh8818/ssl.effective.date.exporter/cmd/

##### To reduce the final image size, start a new stage with alpine from scratch #####
FROM centos:7
RUN yum install -y openssl

# Run as root
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/siangyeh8818/ssl.effective.date.exporter/ssl-effective-date-exporter /usr/local/bin/ssl-effective-date-exporter

# EXPOSE 8081

ENTRYPOINT [ "ssl-effective-date-exporter" ] 
