package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/coltea/owlalpha/backend/internal/bootstrap"
	"github.com/coltea/owlalpha/backend/internal/controller/auth"
	"github.com/coltea/owlalpha/backend/internal/controller/health"
	"github.com/coltea/owlalpha/backend/internal/controller/report"
	"github.com/coltea/owlalpha/backend/internal/controller/settings"
	"github.com/coltea/owlalpha/backend/internal/logic"
	"github.com/coltea/owlalpha/backend/internal/middleware"
)

func main() {
	ctx := context.Background()
	deps, err := bootstrap.New(ctx)
	if err != nil {
		panic(fmt.Errorf("bootstrap failed: %w", err))
	}

	logic.RegisterServices(deps)

	s := g.Server()
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Use(middleware.CORS)

	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			health.NewV1(),
			auth.NewV1(),
		)

		group.Middleware(middleware.Auth)
		group.Bind(
			report.NewV1(),
			settings.NewV1(),
		)
	})

	s.Run()
}
