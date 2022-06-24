package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/cmd/service/library"
	"github.com/evan-buss/openbooks/cmd/service/library/bleve"
	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"github.com/go-chi/chi/v5"
	"github.com/soheilhy/cmux"
)

type server struct {
	searcher library.Searcher
	config   struct {
		indexDir string
		port     string
		debug    bool
	}
	openbooksv1.UnimplementedOpenBooksServiceServer
}

func main() {
	var server server
	flag.StringVar(&server.config.port, "port", "8080", "Port to listen on.")
	flag.StringVar(&server.config.indexDir, "index-dir", "index", "The place where files are downloaded and processed.")
	flag.BoolVar(&server.config.debug, "debug", true, "Enable addition debug functionality.")
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	indexer := library.DirectoryIndexer{Path: filepath.Join(dir, "sources")}
	for b := range indexer.Index() {
		b.GetAuthor()
		// log.Println(b.GetServer(), b.GetTitle())
	}

	searcher, err := bleve.NewBleveSeacher()
	if err != nil {
		panic(err)
	}
	server.searcher = searcher

	grpcServer := server.StartGrpc()
	grpcGatewayMux := server.StartGrpcGateway()

	router := chi.NewRouter()
	router.Mount("/", grpcGatewayMux)
	router.Mount("/swagger", server.SwaggerDocHandlers())

	httpServer := http.Server{
		Handler: router,
	}

	l, err := net.Listen("tcp", ":"+server.config.port)
	if err != nil {
		panic(err)
	}

	log.Println("Listening on ", l.Addr().String())

	m := cmux.New(l)
	httpL := m.Match(cmux.HTTP1Fast())
	grpcL := m.Match(cmux.HTTP2())

	go grpcServer.Serve(grpcL)
	go httpServer.Serve(httpL)
	m.Serve()
}
