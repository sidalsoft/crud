package sales

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
type SalesService struct {
	//db *sql.DB
	pool *pgxpool.Pool
}

//NewService ..
func NewSalesService(pool *pgxpool.Pool) *SalesService {
	return &SalesService{pool: pool}
}

//Sales ...
type Sales struct {
	ID         int64     `json:"id"`
	ManagerId  int64     `json:"managerId"`
	CustomerId *int64    `json:"customerId"`
	Created    time.Time `json:"created"`
}

func (s *SalesService) All(ctx context.Context) (cs []*Sales, err error) {

	sqlStatement := `select * from sales`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Sales{}
		err := rows.Scan(
			&item.ID,
			&item.ManagerId,
			&item.CustomerId,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
		}
		cs = append(cs, item)
	}

	return cs, nil
}

func (s *SalesService) ByID(ctx context.Context, id int64) (*Sales, error) {
	item := &Sales{}

	err := s.pool.QueryRow(ctx, `
SELECT id, manager_id, customer_id, created FROM sales WHERE id=$1`, id).Scan(
		&item.ID,
		&item.ManagerId,
		&item.CustomerId,
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

func (s *SalesService) Delete(ctx context.Context, id int64) (*Sales, error) {
	item := &Sales{}

	err := s.pool.QueryRow(ctx, `
DELETE FROM sales  WHERE id=$1 RETURNING *`, id).Scan(
		&item.ID,
		&item.ManagerId,
		&item.CustomerId,
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

func (s *SalesService) Save(ctx context.Context, customer *Sales) (c *Sales, err error) {

	item := &Sales{}

	if customer.ID == 0 {
		err = s.pool.QueryRow(ctx, `INSERT INTO sales(manager_id, customer_id) values($1, $2) RETURNING *`, customer.ManagerId, customer.CustomerId).Scan(
			&item.ID,
			&item.ManagerId,
			&item.CustomerId,
			&item.Created)
	} else {
		err = s.pool.QueryRow(ctx, `UPDATE sales SET manager_id=$1, customer_id=$2 where id=$3 RETURNING *`, customer.ManagerId, customer.CustomerId, customer.ID).Scan(
			&item.ID,
			&item.ManagerId,
			&item.CustomerId,
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

func (s *SalesService) TotalByManager(ctx context.Context, managerId int64) (int, error) {
	var total int32

	err := s.pool.QueryRow(ctx, `
SELECT sum(price*qty) as s from sale_positions sp join sales s on s.id = sp.sale_id where s.manager_id=$1;`, managerId).Scan(
		&total)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	}

	//if err != nil {
	//	log.Println(err)
	//	return 0, ErrInternal
	//}
	return int(total), nil

}
