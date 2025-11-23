FROM docker.io/library/golang:1.25.4-alpine3.22 as build
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN mkdir /dist
RUN go build -o /dist/tux_triarch ./cmd/tux_triarch

FROM docker.io/library/alpine:3.22
WORKDIR /app
COPY --from=build /dist/tux_triarch .
ENTRYPOINT [ "./tux_triarch", "serve" ]
