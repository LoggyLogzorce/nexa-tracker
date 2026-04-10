package status

import "github.com/google/uuid"

type Status struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProjectID  uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"project_id"`
	Name       string    `gorm:"not null;size:50" json:"name"`
	Color      string    `gorm:"default:'#cccccc';size:7" json:"color"`
	OrderIndex int       `gorm:"default:0" json:"order_index"`
}
