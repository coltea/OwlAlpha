package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"github.com/coltea/owlalpha/backend/internal/service"
)

type StatusReq struct {
	g.Meta `path:"/health" method:"get" tags:"Health" summary:"Health status"`
}

type StatusRes struct {
	service.HealthStatus
}
