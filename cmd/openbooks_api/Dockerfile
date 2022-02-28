FROM golang as build
WORKDIR /go/src/
COPY . .

ENV CGO_ENABLED=0
RUN go get -d -v ./...
RUN go install -v ./...
WORKDIR /go/src/cmd/openbooks_api/
RUN go build

FROM gcr.io/distroless/static as app
WORKDIR /app
COPY --from=build /go/src/cmd/openbooks_api/openbooks_api .

EXPOSE 80
VOLUME [ "/app/sources" ]
VOLUME [ "/app/openbooks.bleve" ]

ENTRYPOINT ["./openbooks_api", "--port", "80", "--log"]
