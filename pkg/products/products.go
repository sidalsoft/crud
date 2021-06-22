package products

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

//ErrNotFound ...
var ErrNotFound = errors.New("item not found")

//ErrInternal ...
var ErrInternal = errors.New("internal error")

//Service ..
type ProductService struct {
	//db *sql.DB
	pool *pgxpool.Pool
}

//NewService ..
func NewProductService(pool *pgxpool.Pool) *ProductService {
	return &ProductService{pool: pool}
}

//Product ...
type Product struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Price   int       `json:"price"`
	Qty     int       `json:"qty"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

func (s *ProductService) All(ctx context.Context) (cs []*Product, err error) {

	sqlStatement := `select * from products`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Product{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Price,
			&item.Qty,
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

func (s *ProductService) ByID(ctx context.Context, id int64) (*Product, error) {
	item := &Product{}

	err := s.pool.QueryRow(ctx, `
SELECT id, name, price, qty, active, created FROM products WHERE id=$1`, id).Scan(
		&item.ID,
		&item.Name,
		&item.Price,
		&item.Qty,
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

func (s *ProductService) Delete(ctx context.Context, id int64) (*Product, error) {
	item := &Product{}

	err := s.pool.QueryRow(ctx, `
DELETE FROM products  WHERE id=$1 RETURNING *`, id).Scan(
		&item.ID,
		&item.Name,
		&item.Price,
		&item.Qty,
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

func (s *ProductService) Save(ctx context.Context, customer *Product) (c *Product, err error) {

	item := &Product{}

	if customer.ID == 0 {
		err = s.pool.QueryRow(ctx, `INSERT INTO products(name, price, qty) values($1, $2, $3) RETURNING *`, customer.Name, customer.Price, customer.Qty).Scan(
			&item.ID,
			&item.Name,
			&item.Price,
			&item.Qty,
			&item.Active,
			&item.Created)
	} else {
		err = s.pool.QueryRow(ctx, `UPDATE products SET name=$1, price=$2, qty=$3 where id=$4 RETURNING *`, customer.Name, customer.Price, customer.Qty, customer.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Price,
			&item.Qty,
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
