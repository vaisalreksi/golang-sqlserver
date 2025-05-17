package models

type Tier string

const (
	NonMember Tier = "Non Member"
	Basic     Tier = "Basic"
	Premium   Tier = "Premium"
)

type PriceDetail struct {
	Id      int  `db:"id" gorm:"primaryKey;autoIncrement"`
	PriceId int  `db:"price_id" json:"price_id"`
	Tier    Tier `db:"tier" json:"tier" gorm:"type:ENUM('Non Member', 'Basic', 'Premium')"`
	Price   int  `db:"price" json:"price"`
}
