FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY main.go ./

RUN CGO_ENABLED=0 go build -o /server

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /server /server

EXPOSE 8087

USER nonroot:nonroot

ENTRYPOINT ["/server"]