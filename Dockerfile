FROM node:lts as web
WORKDIR /web
COPY ./server/app/ .
RUN npm install
RUN npm run build

FROM golang:latest as build
WORKDIR /go/src/
COPY . .
COPY --from=web /web/build/ ./build

# Force static linking so we don't need any deps
RUN go get github.com/rakyll/statik
RUN statik -src /go/src/build -dest /go/src/server

ENV CGO_ENABLED=0
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build

FROM gcr.io/distroless/static as app
WORKDIR /app
COPY --from=build /go/src/openbooks .

ENV IS_DOCKER true

EXPOSE 80
VOLUME [ "/books" ]

ENTRYPOINT ["./openbooks", "-name", "docker", "-dir", "/books", "-port", "80"]
