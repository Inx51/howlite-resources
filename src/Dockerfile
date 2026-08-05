FROM golang:1.26.5-alpine3.24 AS build
WORKDIR /app
RUN apk add --no-cache gcc musl-dev pkgconfig zeromq-dev czmq-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/google/go-licenses@latest && \
    CGO_ENABLED=1 go-licenses save ./... --save_path=/thirdparty-licenses --ignore github.com/inx51/howlite-resources && \
    CGO_ENABLED=1 go-licenses report ./... --ignore github.com/inx51/howlite-resources > /THIRD_PARTY_LICENSES.csv && \
    CGO_ENABLED=1 GOOS=linux go build -o /howlite-resources

FROM alpine:3.24.1
RUN apk add --no-cache zeromq czmq
COPY --from=build /howlite-resources /howlite-resources
COPY --from=build /thirdparty-licenses /thirdparty-licenses
COPY --from=build /THIRD_PARTY_LICENSES.csv /THIRD_PARTY_LICENSES.csv
COPY NOTICE /NOTICE
# Record the license of every apk package present in the runtime image. The Go
# license report above only covers Go modules; the native libraries reach the
# binary through cgo and are invisible to it. libzmq (LGPL-3.0+ with linking
# exception) and czmq (MPL-2.0) both require their terms travel with the image.
RUN mkdir -p /thirdparty-licenses/native && \
    awk -F: '/^P:/{p=$2} /^L:/{print p ": " $2}' /lib/apk/db/installed \
      | sort -u > /thirdparty-licenses/native/APK_PACKAGE_LICENSES.txt
RUN mkdir -p /tmp/howlite
EXPOSE 8080
CMD ["/howlite-resources"]