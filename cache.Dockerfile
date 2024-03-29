FROM public.ecr.aws/bitnami/golang:1.17 AS cache
WORKDIR /go/src/app
COPY go.* ./
COPY internal ./internal
RUN go mod download
COPY microservices/cache ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s' -o cache .

FROM scratch
COPY --from=cache /go/src/app/cache /cache
CMD ["/cache"]