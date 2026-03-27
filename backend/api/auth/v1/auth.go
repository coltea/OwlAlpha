package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"github.com/coltea/owlalpha/backend/internal/service"
)

type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"Auth" summary:"Admin login"`
	Username string `json:"username" v:"required#用户名不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
}

type LoginRes struct {
	service.LoginOutput
}
