package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"github.com/coltea/owlalpha/backend/internal/service"
)

type ListReq struct {
	g.Meta `path:"/reports" method:"get" tags:"Reports" summary:"Report list"`
}

type ListRes struct {
	Items []service.ReportItem `json:"items"`
}
