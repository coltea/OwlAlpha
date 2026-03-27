package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/coltea/owlalpha/backend/internal/service"
)

func Auth(r *ghttp.Request) {
	header := r.Header.Get("Authorization")
	if header == "" {
		r.Response.WriteStatusExit(401, "missing authorization header")
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		r.Response.WriteStatusExit(401, "invalid authorization header")
		return
	}

	user, err := service.Auth().ValidateToken(r.Context(), parts[1])
	if err != nil {
		r.SetError(gerror.Wrap(err, "unauthorized"))
		r.Response.WriteStatusExit(401, "unauthorized")
		return
	}

	r.SetCtxVar("user", user)
	r.Middleware.Next()
}
