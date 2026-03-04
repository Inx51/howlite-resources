FROM golang:1.26.0-alpine3.23 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/google/go-licenses@latest && \
    go-licenses save ./... --save_path=/thirdparty-licenses --ignore github.com/inx51/howlite-resources && \
    CGO_ENABLED=0 GOOS=linux go build -o /howlite-resources

FROM alpine:3.23.2
COPY --from=build /howlite-resources /howlite-resources
COPY --from=build /thirdparty-licenses /thirdparty-licenses
RUN mkdir -p /tmp/howlite
EXPOSE 8080
CMD ["/howlite-resources"]