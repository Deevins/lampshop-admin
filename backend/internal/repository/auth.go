// internal/repository/auth_repository.go

package repository

import (
	"context"
	"errors"
	"github.com/Deevins/lampshop-admin-backend/internal/repository/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// AuthRepository описывает метод получения пользователя (для входа).
type AuthRepository interface {
	VerifyUser(ctx context.Context, username, password string) error
}

// authRepoPG реализует AuthRepository через базу данных (sqlc + pgx).
type authRepoPG struct {
	queries *sql.Queries
}

// NewAuthRepository создаёт новую реализацию AuthRepository, используя pgxpool.
func NewAuthRepository(pool *pgxpool.Pool) AuthRepository {
	return &authRepoPG{
		queries: sql.New(pool),
	}
}

// VerifyUser читает хеш пароля из БД и сверяет его с введённым.
func (r *authRepoPG) VerifyUser(ctx context.Context, username, password string) error {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return errors.New("invalid username or password")
	}
	// user.HashedPassword — это строка вида bcrypt хеш
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return errors.New("invalid username or password")
	}
	return nil
}
