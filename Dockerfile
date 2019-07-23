FROM golang:1.12-stretch AS build

RUN go get -u github.com/golang/dep/cmd/dep

ARG GOARCH="amd64"
ARG GOOS="linux"
ARG CGO_ENABLED=0

WORKDIR /go/src/kubectl-dynamicns

COPY . /go/src/kubectl-dynamicns

RUN apt-get install git gcc g++ binutils
RUN dep ensure --vendor-only
RUN CGO_ENABLED=$CGO_ENABLED GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags "-X main.VERSION=v0.0.1" -o /go/bin/kubectl-dynamicns .

FROM alpine:3.9
RUN apk add --no-cache ca-certificates
COPY --from=build /go/bin/kubectl-dynamicns /kubectl-dynamicns
ENTRYPOINT ["/kubectl-dynamicns"]