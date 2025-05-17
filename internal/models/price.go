package models

type Price struct {
	Id           int           `db:"id" gorm:"primaryKey;autoIncrement"`
	ProductId    int           `db:"product_id" json:"product_id"`
	Unit         string        `db:"unit" json:"unit"`
	PriceDetails []PriceDetail `json:"price_details"`
}
