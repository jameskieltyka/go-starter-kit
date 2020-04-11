package internal

type ApplicationContext interface {
	Version() string
	JWTKey() auth.Something
	Auth() string
	EndpointScope() []string
}

func SetupServer(cfg ApplicationContext) *config.Server{
	basePath := fmt.Sprintf("/%v", cfg.Version())

	server := config.NewServer().
		WithMiddleWare(
			middlewares.CorrelationId
		).
		WithRoutes(
			ProductRoute(),
		)
	return server
}