package dto

type Subscribe struct {
	ID       uint64 `json:"id" form:"id" binding:"required"`
	Response bool   `json:"response" form:"response" binding:"required"`
	UserID   uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
}
