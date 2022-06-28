package main

import (
	"flag"
	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/cmd/service/library"
	"github.com/go-chi/chi/v5"
	"github.com/soheilhy/cmux"
)

type server struct {
	searcher library.Searcher
	indexer  library.Indexer
	config   *serverConfig
	log      *zap.SugaredLogger

	openbooksv1.UnimplementedOpenBooksServiceServer
	openbooksv1.UnimplementedAdminServiceServer
}

type serverConfig struct {
	indexDir string
	port     string
	debug    bool
}

func main() {
	server := &server{config: &serverConfig{}}
	flag.StringVar(&server.config.port, "port", "8080", "Port to listen on.")
	flag.StringVar(&server.config.indexDir, "index-dir", "openbooks.bleve", "The place where files are downloaded and processed.")
	flag.BoolVar(&server.config.debug, "debug", true, "Enable addition debug functionality.")
	flag.Parse()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger := zapLogger.Sugar()
	server.log = logger
	defer logger.Sync()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	indexer := &library.DirectoryIndexer{Path: filepath.Join(dir, "sources")}
	server.indexer = indexer

	searcher, err := library.NewBleveSeacher(server.config.indexDir)
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

	server.log.Infow("OpenBooks Service Started", "address", l.Addr().String())

	m := cmux.New(l)
	httpL := m.Match(cmux.HTTP1Fast())
	grpcL := m.Match(cmux.HTTP2())

	w := errgroup.Group{}

	w.Go(func() error {
		return grpcServer.Serve(grpcL)
	})

	w.Go(func() error {
		return httpServer.Serve(httpL)
	})

	w.Go(func() error {
		return m.Serve()
	})

	server.log.Fatal(w.Wait())
}
