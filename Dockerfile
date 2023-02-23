## Build
FROM golang AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /urlshortener

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /urlshortener /urlshortener

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/urlshortener"]