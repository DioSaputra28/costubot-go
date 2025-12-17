package middlewares

import (
	"contact-management/src/helpers"
	"contact-management/src/utils"
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helpers.UnauthorizedResponse(w, "Unauthorized")
			return
		}

		token := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		username, err := utils.VerifyToken(token)
		if err != nil {
			helpers.UnauthorizedResponse(w, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}
