
#################
# Builder image #
#################
FROM golang:1.10-alpine as builder

ENV depVersion=0.4.1
ENV projectName=go-release
ENV projectOwner=alexejk

RUN apk add --no-cache alpine-sdk zip

# Install Dep
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v${depVersion}/dep-linux-amd64 \
	&& chmod +x /usr/local/bin/dep

WORKDIR ${GOPATH}/src/github.com/${projectOwner}/${projectName}

# Pull Dependencies
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

# Required tools
RUN go get github.com/vektra/mockery/.../
#RUN go get github.com/wlbr/templify

# Get the source in and build
COPY . .
