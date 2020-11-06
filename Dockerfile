FROM node:15.1.0 as web
WORKDIR /web
COPY ./server/app/ .
RUN npm install
RUN npm run build

FROM golang:1.15.3 as build
WORKDIR /go/src/
COPY . .
COPY --from=web /web/build/ ./build

# Force static linking so we don't need any deps
ENV CGO_ENABLED=0
RUN go get -d -v ./...
RUN go install -v ./...
RUN go get github.com/rakyll/statik
RUN statik -src ./build -dest /server
RUN go build

FROM alpine:3.12.1 as app
# Add SSL certs
RUN apk add ca-certificates && update-ca-certificates

WORKDIR /app
COPY --from=build /go/src/openbooks .
EXPOSE 80
ENTRYPOINT ["./openbooks", "-name", "openbooks", "-port", "80"]
