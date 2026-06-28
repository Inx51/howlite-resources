FROM golang:1.26.4-alpine3.23 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/google/go-licenses@latest && \
    CGO_ENABLED=0 go-licenses save ./... --save_path=/thirdparty-licenses --ignore github.com/inx51/howlite-resources && \
    CGO_ENABLED=0 go-licenses report ./... --ignore github.com/inx51/howlite-resources > /THIRD_PARTY_LICENSES.csv && \
    CGO_ENABLED=0 GOOS=linux go build -o /howlite-resources

FROM alpine:3.23.5
COPY --from=build /howlite-resources /howlite-resources
COPY --from=build /thirdparty-licenses /thirdparty-licenses
COPY --from=build /THIRD_PARTY_LICENSES.csv /THIRD_PARTY_LICENSES.csv
RUN mkdir -p /tmp/howlite
EXPOSE 8080
CMD ["/howlite-resources"]