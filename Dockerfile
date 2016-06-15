# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/kumarsarath588/chef

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/kumarsarath588/chef \ 
&& go install github.com/kumarsarath588/chef

# Run the outyet command by default when the container starts.
ADD ./key.pem /go/bin/key.pem
ENTRYPOINT /go/bin/chef
# Document that the service listens on port 8080.
EXPOSE 3002
