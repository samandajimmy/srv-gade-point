
FROM artifactory.pegadaian.co.id:8084/golang:1.11 as build-env

# add ssl certificate
ADD ssl_certificate.crt /usr/local/share/ca-certificates/ssl_certificate.crt
RUN chmod 644 /usr/local/share/ca-certificates/ssl_certificate.crt && update-ca-certificates

RUN apt-get update && apt-get install git
# All these steps will be cached

RUN mkdir /srv-gade-point
WORKDIR /srv-gade-point

# Force the go compiler to use modules
ENV GO111MODULE=on

# Force to download lib from nexus pgdn
ENV GOPROXY="https://artifactory.pegadaian.co.id/repository/go-group-01/"

# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# CHECK VERSION OF GIT
RUN git version

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/srv-gade-point

# Second step to build minimal image
FROM artifactory.pegadaian.co.id:8084/alpine:3.7
COPY --from=build-env /go/bin/srv-gade-point /go/bin/srv-gade-point
COPY --from=build-env /srv-gade-point/entrypoint.sh /srv-gade-point/entrypoint.sh
COPY --from=build-env /srv-gade-point/migrations /migrations

# add apk ca certificate
RUN apk add ca-certificates

# set timezone
RUN apk add tzdata
RUN ls /usr/share/zoneinfo
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
RUN echo "Asia/Jakarta" > /etc/timezone
RUN apk del tzdata

EXPOSE 8080
ENTRYPOINT ["sh", "/srv-gade-point/entrypoint.sh"]