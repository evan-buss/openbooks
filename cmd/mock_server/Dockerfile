FROM golang as build
WORKDIR /go/src/
COPY . .

ENV CGO_ENABLED=0
WORKDIR /go/src/cmd/mock_server/
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build

FROM gcr.io/distroless/static as app
WORKDIR /app
COPY --from=build /go/src/cmd/mock_server/ .

EXPOSE 6667
EXPOSE 6668
EXPOSE 6669

ENTRYPOINT ["./mock_server"]
