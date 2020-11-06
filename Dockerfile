FROM node:15.1.0 as web
WORKDIR /web
COPY ./server/app/ .
RUN npm install
RUN npm run build

FROM golang:1.15.3 as build
WORKDIR /go/src/
COPY . .
COPY --from=web /web/build/ ./build
RUN go get -d -v ./...
RUN go install -v ./...
RUN go get github.com/rakyll/statik
RUN statik -src ./build -dest /server
RUN go build

FROM alpine:3.12.1 as app
WORKDIR /app
COPY --from=build /go/src/openbooks .
EXPOSE 80
ENTRYPOINT ["./openbooks"]
