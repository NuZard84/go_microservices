package order

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close() error
	PutOrder(ctx context.Context, order Order) error
	GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, order Order) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, `
	INSERT INTO orders (id, account_id, total_price, created_at)
	VALUES ($1, $2, $3, $4)
	`, order.ID, order.AccountID, order.TotalPrice, order.CreatedAt)
	if err != nil {
		return err
	}

	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))

	for _, product := range order.Products {
		_, err = stmt.ExecContext(ctx, order.ID, product.ID, product.Quantity)
		if err != nil {
			return err
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}
	stmt.Close()

	return nil
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
		o.id,
		o.created_at, 
		o.account_id, 
		o.total_price::money::numeric::float8, 
		op.product_id, 
		op.quantity 
		FROM orders o JOIN order_products op ON(o.id = op.order_id)
		WHERE o.account_id = $1
		ORDER BY o.id`,
		accountId, 
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}
	lastOrder := &Order{}
	orderedProduct := &OrderedProduct{}
	products := []OrderedProduct{}

	for rows.Next() {
		if err := rows.Scan(
			&lastOrder.ID,
			&lastOrder.CreatedAt,
			&lastOrder.AccountID,
			&lastOrder.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		if lastOrder.ID != "" && lastOrder.ID != lastOrder.ID {
			newOrder := Order{
				ID:         lastOrder.ID,
				CreatedAt:  lastOrder.CreatedAt,
				AccountID:  lastOrder.AccountID,
				TotalPrice: lastOrder.TotalPrice,
				Products:   products,
			}
			orders = append(orders, newOrder)
			products = []OrderedProduct{}
		}
		products = append(products, *orderedProduct)
	}

	if lastOrder.ID != "" && lastOrder.ID != lastOrder.ID {
		newOrder := Order{
			ID:         lastOrder.ID,
			CreatedAt:  lastOrder.CreatedAt,
			AccountID:  lastOrder.AccountID,
			TotalPrice: lastOrder.TotalPrice,
			Products:   products,
		}
		orders = append(orders, newOrder)
	}

	return orders, nil
}
