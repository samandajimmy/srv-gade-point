############################################################
# Dockerfile to build Golang project with go tool
############################################################

FROM golang:1.11

ENV GO111MODULE=on

# Copy whole code to WORKDIR
COPY . /go/src/gade/srv-gade-point
WORKDIR /go/src/gade/srv-gade-point

# Set TimeZone to Jakarta Indonesia
ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Install dependencies based on go.mod and go.sum
RUN go get

# Build server binary
RUN go build

# Run the entrypoint.sh
ENTRYPOINT ["sh", "/go/src/gade/srv-gade-point/entrypoint.sh"]

