package auth

import (
	"context"
	"errors"

	"encore.app/services/auth/db"
	"encore.app/services/auth/password"
	"encore.dev/rlog"
)

var secrets struct {
	FirstAdminUsername string
	FirstAdminPassword string
}

func createFirstAdmin(query db.Querier) {
	if secrets.FirstAdminUsername == "" || secrets.FirstAdminPassword == "" {
		return
	}

	ctx := context.Background()
	id, err := query.CheckUserExists(ctx, secrets.FirstAdminUsername)
	if err != nil && !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to check if first admin exists", "error", err)
		panic(err)
	}
	if id != 0 {
		return
	}

	hashed, err := password.HashPassword(secrets.FirstAdminPassword)
	if err != nil {
		rlog.Error("failed to hash first admin password", "error", err)
		return
	}

	_, err = query.RegisterAdmin(ctx, db.RegisterAdminParams{
		Username:     secrets.FirstAdminUsername,
		PasswordHash: hashed + string(hashed),
	})
	if err != nil {
		rlog.Error("failed to create first admin user", "error", err)
		return
	}
	rlog.Info("created first admin user", "username", secrets.FirstAdminUsername)
}
