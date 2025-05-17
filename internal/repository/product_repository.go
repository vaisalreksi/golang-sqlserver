package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-sqlserver/internal/models"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id int) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, keyword, category, tier string) ([]models.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO Product (Name, Product_Category, Description)
		VALUES (@name, @product_category, @description);
		SELECT SCOPE_IDENTITY();
	`

	row := r.db.QueryRowContext(ctx, query,
		sql.Named("name", product.Name),
		sql.Named("product_category", product.ProductCategory),
		sql.Named("description", product.Description),
	)

	return row.Scan(&product.Id)
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT Id, Name, Product_Category, Description
		FROM Product
		WHERE Id = @id
	`

	product := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&product.Id,
		&product.Name,
		&product.ProductCategory,
		&product.Description,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT Id, Name, Product_Category, Description
		FROM Product
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.ProductCategory,
			&product.Description,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE Product
		SET Name = @name,
			Product_Category = @product_category,
			Description = @description
		WHERE Id = @id
	`

	result, err := r.db.ExecContext(ctx, query,
		sql.Named("name", product.Name),
		sql.Named("product_category", product.ProductCategory),
		sql.Named("description", product.Description),
		sql.Named("id", product.Id),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", product.Id)
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM Product WHERE Id = @id`

	result, err := r.db.ExecContext(ctx, query, sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	return nil
}

func (r *productRepository) Search(ctx context.Context, keyword, category, tier string) ([]models.Product, error) {
	query := `
        SELECT DISTINCT
            p.Id,
            p.Name,
            p.Product_Category,
            p.Description,
            pr.Id as PriceId,
            pr.Unit,
            pd.Id as PriceDetailId,
            pd.Tier,
            pd.Price
        FROM Product p
        LEFT JOIN Price pr ON p.Id = pr.Product_Id
        LEFT JOIN PriceDetail pd ON pr.Id = pd.Price_Id
        WHERE 1=1
    `
	params := make([]interface{}, 0)
	paramCount := 1

	if keyword != "" {
		query += fmt.Sprintf(" AND p.Name LIKE @p%d", paramCount)
		params = append(params, sql.Named(fmt.Sprintf("p%d", paramCount), "%"+keyword+"%"))
		paramCount++
	}

	if category != "" {
		query += fmt.Sprintf(" AND p.Product_Category = @p%d", paramCount)
		params = append(params, sql.Named(fmt.Sprintf("p%d", paramCount), category))
		paramCount++
	}

	if tier != "" {
		query += fmt.Sprintf(" AND pd.Tier = @p%d", paramCount)
		params = append(params, sql.Named(fmt.Sprintf("p%d", paramCount), tier))
	}

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productMap := make(map[int]*models.Product)

	for rows.Next() {
		var (
			product  models.Product
			priceId  sql.NullInt64
			unit     sql.NullString
			detailId sql.NullInt64
			tier     sql.NullString
			price    sql.NullInt64
		)

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.ProductCategory,
			&product.Description,
			&priceId,
			&unit,
			&detailId,
			&tier,
			&price,
		)
		if err != nil {
			return nil, err
		}

		if existingProduct, ok := productMap[product.Id]; ok {
			if priceId.Valid {
				priceExists := false
				for _, p := range existingProduct.Prices {
					if p.Id == int(priceId.Int64) {
						if detailId.Valid {
							p.PriceDetails = append(p.PriceDetails, models.PriceDetail{
								Id:      int(detailId.Int64),
								PriceId: int(priceId.Int64),
								Tier:    models.Tier(tier.String),
								Price:   int(price.Int64),
							})
						}
						priceExists = true
						break
					}
				}
				if !priceExists && unit.Valid {
					prices := models.Price{
						Id:        int(priceId.Int64),
						ProductId: product.Id,
						Unit:      unit.String,
					}
					if detailId.Valid {
						prices.PriceDetails = []models.PriceDetail{{
							Id:      int(detailId.Int64),
							PriceId: int(priceId.Int64),
							Tier:    models.Tier(tier.String),
							Price:   int(price.Int64),
						}}
					}
					existingProduct.Prices = append(existingProduct.Prices, prices)
				}
			}
		} else {
			product.Prices = make([]models.Price, 0)
			if priceId.Valid && unit.Valid {
				prices := models.Price{
					Id:        int(priceId.Int64),
					ProductId: product.Id,
					Unit:      unit.String,
				}
				if detailId.Valid {
					prices.PriceDetails = []models.PriceDetail{{
						Id:      int(detailId.Int64),
						PriceId: int(priceId.Int64),
						Tier:    models.Tier(tier.String),
						Price:   int(price.Int64),
					}}
				}
				product.Prices = append(product.Prices, prices)
			}
			productMap[product.Id] = &product
		}
	}

	// Convert map to slice
	result := make([]models.Product, 0, len(productMap))
	for _, product := range productMap {
		// If tier is specified, only include products that have matching tier
		if tier != "" {
			hasTier := false
			for _, price := range product.Prices {
				for _, detail := range price.PriceDetails {
					if detail.Tier == models.Tier(tier) {
						hasTier = true
						break
					}
				}
				if hasTier {
					break
				}
			}
			if !hasTier {
				continue
			}
		}
		result = append(result, *product)
	}

	return result, nil
}
