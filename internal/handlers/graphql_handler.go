package handlers

import (
	"context"
	gql "looky/internal/graphql"
	"looky/internal/utils"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
)

func GraphQLHandler() fiber.Handler {
	h := handler.NewDefaultServer(
		gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{}}),
	)

	h.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		return next(ctx)
	})

	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token != "" {
			claims, err := utils.ValidateJWT(token)
			if err == nil {
				ctx := gql.WithClaims(r.Context(), claims.UserID, string(claims.Role))
				r = r.WithContext(ctx)
			}
		}
		h.ServeHTTP(w, r)
	})
}

func GraphQLPlayground() fiber.Handler {
	h := playground.Handler("GraphQL", "/graphql")
	return adaptor.HTTPHandler(h)
}
