package main

import (
	"flag"
	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"github.com/evan-buss/openbooks/search"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/soheilhy/cmux"
)

type server struct {
	searcher search.Searcher
	indexer  search.Indexer
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
	flag.StringVar(&server.config.port, "port", "8000", "Port to listen on.")
	flag.StringVar(&server.config.indexDir, "index-dir", "openbooks.bleve", "The place where files are downloaded and processed.")
	flag.BoolVar(&server.config.debug, "debug", true, "Enable addition debug functionality.")
	flag.Parse()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger := zapLogger.Sugar()
	server.log = logger
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	indexer := search.NewDirectoryIndexer(filepath.Join(dir, "sources"), logger)
	server.indexer = indexer

	searcher, err := search.NewBleveSeacher(server.config.indexDir)
	if err != nil {
		panic(err)
	}
	server.searcher = searcher

	grpcServer := server.StartGrpc()
	grpcGatewayMux := server.StartGrpcGateway()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Mount("/", grpcGatewayMux)
	router.Mount("/swagger", server.SwaggerDocHandlers())

	httpServer := http.Server{Handler: router}

	l, err := net.Listen("tcp", "0.0.0.0:"+server.config.port)
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
