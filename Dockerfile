# argument for Go version
ARG GO_VERSION=1.11

# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
# from https://medium.com/@pierreprinetti/the-go-1-11-dockerfile-a3218319d191
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# git required for go mod
RUN apk add --no-cache git

# Working directory will be created if it does not exist
WORKDIR /src

# We use go modules; copy go.mod and go.sum
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import code
COPY ./ ./

# Build the executable
RUN CGO_ENABLED=0 go build \
	-installsuffix 'static' \
	-o /app .

# STAGE 2: build the container to run
FROM scratch AS finale

# import user and group files
COPY --from=build /user/group /user/passwd /etc/

# copy compiled app
COPY --from=build /app /app

# copy assets
COPY /assets /assets

# expose port 8888
EXPOSE 8888

# run as unpriviledged user
USER nobody:nobody

# run binary
ENTRYPOINT ["/app"]