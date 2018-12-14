# First step: build Go app
FROM golang:1.11.2

RUN mkdir -p /go/src/realtime
WORKDIR /go/src/realtime
COPY . .

ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go build -o realtime -v ./... 

# Second step: build image from scratch
FROM scratch
COPY --from=0 /go/src/realtime /
ADD assets /assets

CMD ["/realtime"]