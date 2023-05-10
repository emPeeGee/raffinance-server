package entity

type TransactionTag struct {
	TransactionID uint `gorm:"primaryKey"`
	TagID         uint `gorm:"primaryKey"`

	// TODO: Looks like it doesn't help
	// Define foreign key constraints
	// Ensure that transaction_id references transactions(id) and
	// tag_id references tags(id)
	Transaction Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:TransactionID"`
	Tag         Tag         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:TagID"`
}
