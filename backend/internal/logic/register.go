package logic

import (
	"github.com/coltea/owlalpha/backend/internal/bootstrap"
	"github.com/coltea/owlalpha/backend/internal/logic/auth"
	"github.com/coltea/owlalpha/backend/internal/logic/health"
	"github.com/coltea/owlalpha/backend/internal/logic/report"
	"github.com/coltea/owlalpha/backend/internal/logic/settings"
	"github.com/coltea/owlalpha/backend/internal/service"
)

func RegisterServices(deps *bootstrap.Dependencies) {
	service.RegisterHealth(health.New())
	service.RegisterAuth(auth.New(deps))
	service.RegisterReport(report.New(deps))
	service.RegisterSettings(settings.New(deps))
}
