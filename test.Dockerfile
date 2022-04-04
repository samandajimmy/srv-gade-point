
FROM artifactory.pegadaian.co.id:8084/golang:1.16.9 

# add ssl certificate
ADD ssl_certificate.crt /usr/local/share/ca-certificates/ssl_certificate.crt
RUN chmod 644 /usr/local/share/ca-certificates/ssl_certificate.crt && update-ca-certificates

RUN apt-get update && apt-get install git
# All these steps will be cached

RUN mkdir /srv-gade-point
WORKDIR /srv-gade-point

# Force to download lib from nexus pgdn
ENV GOPROXY="https://artifactory.pegadaian.co.id/repository/go-group-01/"

# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# install executable ginkgo
RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.1.3

# COPY the source code as the last step
COPY . .

# Run test
CMD [ "ginkgo", "-r", "--randomize-all", "--randomize-suites", "--fail-on-pending", "--cover" ]
