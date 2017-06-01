FROM golang:1.8

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/mehrat/train-app/api

# Install our dependencies
RUN go get github.com/go-martini/martini
RUN go get github.com/martini-contrib/binding
RUN go get github.com/martini-contrib/render
RUN go get labix.org/v2/mgo
RUN go get labix.org/v2/mgo/bson

# Install api binary globally within container
RUN go install github.com/mehrat/train-app/api

# Set binary as entrypoint
ENTRYPOINT /go/bin/api

# Expose default port (3000)
EXPOSE 3000
