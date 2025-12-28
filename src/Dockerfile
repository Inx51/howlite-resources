FROM golang:1.25.5-alpine3.23 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN go install github.com/google/go-licenses@latest
RUN go-licenses save ./... --save_path=/thirdparty-licenses --ignore github.com/inx51/howlite-resources
RUN CGO_ENABLED=0 GOOS=linux go build -o /howlite-resources

FROM alpine:3.23.2
COPY --from=build /howlite-resources /howlite-resources
COPY --from=build /thirdparty-licenses /thirdparty-licenses
RUN mkdir -p /tmp/howlite
EXPOSE 8080
CMD ["/howlite-resources"]