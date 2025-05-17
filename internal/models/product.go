package models

type ProductCategory string

const (
	Rokok   ProductCategory = "Rokok"
	Obat    ProductCategory = "Obat"
	Lainnya ProductCategory = "Lainnya"
)

type Product struct {
	Id              int             `db:"id" gorm:"primaryKey;autoIncrement"`
	Name            string          `db:"name" json:"name"`
	ProductCategory ProductCategory `db:"product_category" json:"product_category" gorm:"type:ENUM('Rokok', 'Obat', 'Lainnya')"`
	Description     string          `db:"description" json:"description"`
	Prices          []Price         `json:"prices"`
}
