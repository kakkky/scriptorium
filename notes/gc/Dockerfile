FROM golang:1.25-alpine AS build
WORKDIR /src
COPY main.go .
RUN go build -o main main.go

FROM alpine:3.20
COPY --from=build /src/main /main
CMD ["/main"]
