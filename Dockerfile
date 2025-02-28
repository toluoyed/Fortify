FROM golang:1.21.0

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .

RUN go get -d -v ./...

# Build
RUN go build -o /fortify-app

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/fortify-app"]