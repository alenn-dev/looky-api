package graphql

import (
	"context"
	"errors"
	"looky/internal/models"

	"github.com/google/uuid"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func WithClaims(ctx context.Context, userID string, role string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, roleKey, role)
	return ctx
}

func extractClaims(ctx context.Context) (uuid.UUID, models.UserRole, error) {
	userIDStr, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return uuid.Nil, "", errors.New("unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, "", errors.New("invalid user id")
	}

	roleStr, ok := ctx.Value(roleKey).(string)
	if !ok {
		return uuid.Nil, "", errors.New("unauthorized")
	}

	return userID, models.UserRole(roleStr), nil
}
