FROM golang:1.11 as build-env
RUN apt-get update && apt-get install -y git=2.20.1
# All these steps will be cached

RUN mkdir /srv-gade-point
WORKDIR /srv-gade-point

# Force the go compiler to use modules
ENV GO111MODULE=on

# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# CHECK VERSION OF GIT
RUN git version

# Get dependancies - will also be cached if we won't change mod/sum
RUN go get -d -v ./...

# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/srv-gade-point

# Second step to build minimal image
FROM alpine:3.7
COPY --from=build-env /go/bin/srv-gade-point /go/bin/srv-gade-point
COPY --from=build-env /srv-gade-point/entrypoint.sh /srv-gade-point/entrypoint.sh
COPY --from=build-env /srv-gade-point/migrations /migrations
EXPOSE 8080
ENTRYPOINT ["sh", "/srv-gade-point/entrypoint.sh"]