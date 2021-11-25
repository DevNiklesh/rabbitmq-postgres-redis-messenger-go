FROM public.ecr.aws/bitnami/golang:1.17 AS backend
WORKDIR /go/src/app
COPY go.* ./
COPY internal ./internal
RUN go mod download
COPY microservices/backend ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s' -o backend .

FROM scratch
COPY --from=backend /go/src/app/backend /backend
CMD ["/backend"]