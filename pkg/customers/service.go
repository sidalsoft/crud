package customers

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

//ErrNotFound ...
var ErrNotFound = errors.New("item not found")

//ErrInternal ...
var ErrInternal = errors.New("internal error")

//Service ..
type Service struct {
	//db *sql.DB
	pool *pgxpool.Pool
}

//NewService ..
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

//Customer ...
type Customer struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

func (s *Service) All(ctx context.Context) (cs []*Customer, err error) {

	sqlStatement := `select * from customers`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
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

func (s *Service) AllActive(ctx context.Context) (cs []*Customer, err error) {
	rows, err := s.pool.Query(ctx, `select * from customers where active=true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		if err != nil {
			log.Println(err)
			return nil, ErrInternal
		}
		cs = append(cs, item)
	}

	return cs, nil
}

func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
SELECT id, name, phone, active, created FROM customers WHERE id=$1`, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
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

func (s *Service) ChangeActive(ctx context.Context, id int64, active bool) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
UPDATE customers SET active=$2 WHERE id=$1 RETURNING *`, id, active).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
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

func (s *Service) Delete(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
DELETE FROM customers  WHERE id=$1 RETURNING *`, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
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

func (s *Service) Save(ctx context.Context, customer *Customer) (c *Customer, err error) {

	item := &Customer{}

	if customer.ID == 0 {
		err = s.pool.QueryRow(ctx, `INSERT INTO customers(name, phone, password) values($1, $2, $3) RETURNING *`, customer.Name, customer.Phone, hashPassword(customer.Password)).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Password,
			&item.Active,
			&item.Created)
	} else {
		err = s.pool.QueryRow(ctx, `UPDATE customers SET name=$1, phone=$2, password=$3 where id=$4 RETURNING *`, customer.Name, customer.Phone, hashPassword(customer.Password), customer.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
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

func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, "SELECT customer_id FROM customers_tokens WHERE token = $1", token).Scan(&id)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func hashPassword(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
	}
	return string(hash)
}
