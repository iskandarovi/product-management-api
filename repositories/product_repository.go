package repositories

import (
	"context"
	"fmt"
	"product-management-api/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

func (r *ProductRepository) Create(ctx context.Context, product models.Product) (models.Product, error) {
	const query = `
		insert into products (name, price, quantity)
		values ($1, $2, $3)
		returning id, name, price, quantity, created_at;
	`

	var createdProduct models.Product

	err := r.DB.QueryRow(ctx, query,
		product.Name,
		product.Price,
		product.Quantity,
	).Scan(
		&createdProduct.ID,
		&createdProduct.Name,
		&createdProduct.Price,
		&createdProduct.Quantity,
		&createdProduct.CreatedAt,
	)

	if err != nil {
		return models.Product{}, err
	}

	return createdProduct, nil
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	const query = `
		select id, name, price, quantity, created_at FROM products;
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productslar []models.Product

	for rows.Next() {
		var prod models.Product
		err := rows.Scan(&prod.ID, &prod.Name, &prod.Price, &prod.Quantity, &prod.CreatedAt)
		if err != nil {
			return nil, err
		}
		productslar = append(productslar, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productslar, nil
}
func (r *ProductRepository) Search(ctx context.Context, keyword string) ([]models.Product, error) {
	const query = `
	SELECT * FROM products
	WHERE name ILIKE $1;
	`
	qidiruv := "%" + keyword + "%"
	rows, err := r.DB.Query(ctx, query, qidiruv)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productslar []models.Product

	for rows.Next() {
		var prod models.Product
		err := rows.Scan(
			&prod.ID,
			&prod.Name,
			&prod.Price,
			&prod.Quantity,
			&prod.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		productslar = append(productslar, prod)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productslar, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (models.Product, error) {
	const query = `
	select id, name, price, quantity, created_at from products where id = $1 
	`
	var product models.Product

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Product{}, pgx.ErrNoRows
		}
		return models.Product{}, fmt.Errorf("id boyicha olishda xatolik: %w", err)
	}
	return product, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int64, name string, price float64, quantity int) (models.Product, error) {
	const query = `
	update products 
	set name = $1, price = $2, quantity = $3 
	WHERE id = $4
	RETURNING id, name, price, quantity, created_at;
	`
	var updateProduct models.Product
	err := r.DB.QueryRow(ctx, query, name, price, quantity, id).Scan(
		&updateProduct.ID,
		&updateProduct.Name,
		&updateProduct.Price,
		&updateProduct.Quantity,
		&updateProduct.CreatedAt,
	)
	if err != nil {
		return models.Product{}, fmt.Errorf("Update qilishda xatolik > > %w", err)
	}
	return updateProduct, nil
}
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	const query = `
	delete from products where id =$1;
	`
	deleteProduct, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("xatolik >> %w", err)
	}
	if deleteProduct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) Sort(ctx context.Context, order string) ([]models.Product, error) {
	var query string

	switch order {
	case "asc":
		query = `
			SELECT id, name, price, quantity, created_at
			FROM products
			ORDER BY price ASC;
		`
	case "desc":
		query = `
			SELECT id, name, price, quantity, created_at
			FROM products
			ORDER BY price DESC;
		`
	default:
		return nil, fmt.Errorf("noto'g'ri sort tartibi: %s", order)
	}

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("xatolik: %w", err)
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("xatolik: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("xatolik: %w", err)
	}

	return products, nil
}
func (r *ProductRepository) Pagination(ctx context.Context, page, limit int) ([]models.Product, error) {
	if page < 1 || limit < 1 {
		return nil, fmt.Errorf("page va limit 1 dan katta bo'lishi kerak")
	}
	offset := (page - 1) * limit
	const query = `
	select * from products
	order by id
	limit $1 offset $2
	`
	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("xatolik > > %w", err)
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("xatolik: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("xatolik: %w", err)
	}

	return products, nil
}
