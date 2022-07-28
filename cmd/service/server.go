package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"github.com/evan-buss/openbooks/search"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/go-chi/chi/v5"
	"github.com/soheilhy/cmux"
)

type server struct {
	config   *Config
	log      *zap.SugaredLogger
	searcher search.Searcher
	indexer  search.Indexer
	store    Store

	openbooksv1.UnimplementedOpenBooksServiceServer
	openbooksv1.UnimplementedAdminServiceServer
}

func NewServer(config *Config) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger := zapLogger.Sugar()
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

	// searcher, err := search.NewBleveSeacher(config.Indexing.StoreDirectory)
	// if err != nil {
	// 	panic(err)
	// }

	searcher, err := search.NewTypesenseSearcher()
	if err != nil {
		panic(err)
	}

	// searcher, err := search.NewMeiliSearcher()
	// if err != nil {
	// 	panic(err)
	// }

	server := &server{
		config:   config,
		log:      logger,
		indexer:  indexer,
		searcher: searcher,
		store:    NewConfigStore(config),
	}

	grpcServer := server.StartGrpc()
	grpcGatewayMux := server.StartGrpcGateway()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Mount("/", grpcGatewayMux)
	router.Mount("/swagger", server.SwaggerDocHandlers())

	httpServer := http.Server{Handler: router}

	l, err := net.Listen("tcp", "0.0.0.0:"+server.config.Port)
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
