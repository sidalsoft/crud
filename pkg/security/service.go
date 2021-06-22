package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrNoSuchUser      = errors.New("no such user")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInternal        = errors.New("internal error")
)

//Service ..
type AuthService struct {
	//db *sql.DB
	pool *pgxpool.Pool
}

//NewService ..
func NewAuthService(pool *pgxpool.Pool) *AuthService {
	return &AuthService{pool: pool}
}

func (as *AuthService) Auth(login, password string) (ok bool) {
	var item string

	err := as.pool.QueryRow(context.Background(), `SELECT name FROM managers WHERE login=$1 and password=$2`,
		login, password).Scan(&item)

	if err != nil || errors.Is(err, pgx.ErrNoRows) {
		return false
	}
	return true
}

func (as *AuthService) TokenForCustomer(ctx context.Context, phone string, password string) (token string, err error) {
	var hash string
	var id int64
	err = as.pool.QueryRow(ctx, `SELECT id, password FROM customers WHERE phone = $1`, phone).
		Scan(&id, &hash)
	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}
	if err != nil {
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}
	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}
	token = hex.EncodeToString(buffer)
	_, err = as.pool.Exec(ctx, `INSERT INTO customers_tokens(token, customer_id) VALUES($1,$2)`, token, id)
	if err != nil {
		return "", ErrInternal
	}
	return token, nil
}

func (as *AuthService) AuthenticateCustomer(ctx context.Context, token string) (id int64, expire *time.Time, err error) {
	err = as.pool.QueryRow(ctx, `SELECT customer_id, expire FROM customers_tokens WHERE token = $1`, token).Scan(&id, &expire)

	if err == pgx.ErrNoRows {
		return 0, nil, ErrNoSuchUser
	}
	if err != nil {
		return 0, nil, ErrInternal
	}
	return id, expire, nil
}
