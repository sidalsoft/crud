package salePositions

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
type SalePositionsService struct {
	//db *sql.DB
	pool *pgxpool.Pool
}

//NewService ..
func NewSalePositionsService(pool *pgxpool.Pool) *SalePositionsService {
	return &SalePositionsService{pool: pool}
}

//SalePositions ...
type SalePositions struct {
	ID        int64     `json:"id"`
	SaleId    int64     `json:"saleId"`
	ProductId int64     `json:"productId"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Qty       int       `json:"qty"`
	Created   time.Time `json:"created"`
}

func (s *SalePositionsService) All(ctx context.Context) (cs []*SalePositions, err error) {

	sqlStatement := `select * from sale_positions`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &SalePositions{}
		err := rows.Scan(
			&item.ID,
			&item.SaleId,
			&item.ProductId,
			&item.Name,
			&item.Price,
			&item.Qty,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
		}
		cs = append(cs, item)
	}

	return cs, nil
}

func (s *SalePositionsService) ByID(ctx context.Context, id int64) (*SalePositions, error) {
	item := &SalePositions{}

	err := s.pool.QueryRow(ctx, `
SELECT * FROM sale_positions WHERE id=$1`, id).Scan(
		&item.ID,
		&item.SaleId,
		&item.ProductId,
		&item.Name,
		&item.Price,
		&item.Qty,
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

func (s *SalePositionsService) Delete(ctx context.Context, id int64) (*SalePositions, error) {
	item := &SalePositions{}

	err := s.pool.QueryRow(ctx, `
DELETE FROM sale_positions  WHERE id=$1 RETURNING *`, id).Scan(
		&item.ID,
		&item.SaleId,
		&item.ProductId,
		&item.Name,
		&item.Price,
		&item.Qty,
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

func (s *SalePositionsService) Save(ctx context.Context, customer *SalePositions) (c *SalePositions, err error) {

	item := &SalePositions{}

	if customer.ID == 0 {
		err = s.pool.QueryRow(ctx, `INSERT INTO sale_positions(sale_id, product_id, name, price, qty) values($1, $2, $3, $4, $5) RETURNING *`, customer.SaleId, customer.ProductId, customer.Name, customer.Price, customer.Qty).Scan(
			&item.ID,
			&item.SaleId,
			&item.ProductId,
			&item.Name,
			&item.Price,
			&item.Qty,
			&item.Created)
	} else {
		err = s.pool.QueryRow(ctx, `UPDATE sale_positions SET sale_id=$1, product_id=$2, name=$3, price=$4, qty=$5 where id=$6 RETURNING *`, customer.SaleId, customer.ProductId, customer.Name, customer.Price, customer.Qty, customer.ID).Scan(
			&item.ID,
			&item.SaleId,
			&item.ProductId,
			&item.Name,
			&item.Price,
			&item.Qty,
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
