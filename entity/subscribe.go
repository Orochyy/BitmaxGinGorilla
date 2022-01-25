package entity

type Subsribe struct {
	ID       uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Response bool   `gorm:"type:boolean" json:"response"`
	UserID   uint64 `gorm:"not null" json:"-"`
	User     User   `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
}
