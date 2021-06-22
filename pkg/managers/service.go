package managers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sidalsoft/crud/cmd/app/middleware"
	"github.com/sidalsoft/crud/pkg/security"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

//ErrNotFound ...
var ErrNotFound = errors.New("item not found")

//ErrInternal ...
var ErrInternal = errors.New("internal error")

//Service ..
type ManagersService struct {
	//db *sql.DB
	pool *pgxpool.Pool
}

//NewService ..
func NewManagersService(pool *pgxpool.Pool) *ManagersService {
	return &ManagersService{pool: pool}
}

//Managers ...
type Managers struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Salary     int       `json:"salary"`
	Plan       int       `json:"plan"`
	BossId     *int64    `json:"bossId"`
	Department string    `json:"department"`
	Phone      string    `json:"phone"`
	Password   string    `json:"password"`
	Roles      []string  `json:"roles"`
	Active     bool      `json:"active"`
	Created    time.Time `json:"created"`
}

func (s *ManagersService) All(ctx context.Context) (cs []*Managers, err error) {

	sqlStatement := `select * from managers`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Managers{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Salary,
			&item.Plan,
			&item.BossId,
			&item.Department,
			&item.Phone,
			&item.Password,
			&item.Roles,
			&item.Active,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
		}
		cs = append(cs, item)
	}

	return cs, nil
}

func (s *ManagersService) ByID(ctx context.Context, id int64) (*Managers, error) {
	item := &Managers{}

	err := s.pool.QueryRow(ctx, `
SELECT id, name, salary, plan, boss_id, department, login, roles, active, created FROM managers WHERE id=$1`, id).Scan(
		&item.ID,
		&item.Name,
		&item.Salary,
		&item.Plan,
		&item.BossId,
		&item.Department,
		&item.Phone,
		&item.Password,
		&item.Roles,
		&item.Active,
		&item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil

}

func (s *ManagersService) ChangeActive(ctx context.Context, id int64, active bool) (*Managers, error) {
	item := &Managers{}

	err := s.pool.QueryRow(ctx, `
UPDATE managers SET active=$2 WHERE id=$1 RETURNING *`, id, active).Scan(
		&item.ID,
		&item.Name,
		&item.Salary,
		&item.Plan,
		&item.BossId,
		&item.Department,
		&item.Phone,
		&item.Password,
		&item.Roles,
		&item.Active,
		&item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil

}

func (s *ManagersService) Delete(ctx context.Context, id int64) (*Managers, error) {
	item := &Managers{}

	err := s.pool.QueryRow(ctx, `
DELETE FROM managers  WHERE id=$1 RETURNING *`, id).Scan(
		&item.ID,
		&item.Name,
		&item.Salary,
		&item.Plan,
		&item.BossId,
		&item.Department,
		&item.Phone,
		&item.Password,
		&item.Roles,
		&item.Active,
		&item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil

}

func (s *ManagersService) Save(ctx context.Context, customer *Managers) (c *Managers, err error) {
	item := &Managers{}
	if customer.ID == 0 {
		err = s.pool.QueryRow(ctx, `INSERT INTO managers(name, phone, roles, password) values($1, $2, $3, $4) RETURNING *`, customer.Name, customer.Phone, customer.Roles, hashPassword(customer.Password)).Scan(
			&item.ID,
			&item.Name,
			&item.Salary,
			&item.Plan,
			&item.BossId,
			&item.Department,
			&item.Phone,
			&item.Password,
			&item.Roles,
			&item.Active,
			&item.Created)
	} else {
		err = s.pool.QueryRow(ctx, `UPDATE managers SET name=$1, phone=$2, roles=$3 where id=$4 RETURNING *`, customer.Name, customer.Phone, customer.Roles, customer.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Salary,
			&item.Plan,
			&item.BossId,
			&item.Department,
			&item.Phone,
			&item.Password,
			&item.Roles,
			&item.Active,
			&item.Created)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil

}

func (s *ManagersService) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, "SELECT managers_id FROM managers_tokens WHERE token = $1", token).Scan(&id)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ManagersService) HasAnyRole(ctx context.Context, role string) bool {
	var aRoles []string
	id, err := middleware.Authentication(ctx)
	if err != nil {
		return false
	}
	err = s.pool.QueryRow(ctx, "SELECT roles FROM managers WHERE id = $1", id).Scan(&aRoles)
	if err == pgx.ErrNoRows {
		return false
	}
	if err != nil {
		return false
	}
	for _, v := range aRoles {
		if v == role {
			return true
		}
	}
	return false
}

func (s *ManagersService) TokenForManager(ctx context.Context, phone string, password string) (token string, err error) {
	var hash string
	var id int64
	err = s.pool.QueryRow(ctx, `SELECT id, password FROM managers WHERE phone = $1`, phone).
		Scan(&id, &hash)
	if err == pgx.ErrNoRows {
		return "", security.ErrNoSuchUser
	}
	if err != nil {
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", security.ErrInvalidPassword
	}
	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}
	token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO managers_tokens(token, managers_id) VALUES($1,$2)`, token, id)
	if err != nil {
		return "", ErrInternal
	}
	return token, nil
}
func hashPassword(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
	}
	return string(hash)
}
