package priority

import "github.com/google/uuid"

type Priority struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"project_id"`
	Title     string    `gorm:"not null;size:50" json:"title"`
	Color     string    `gorm:"default:'#cccccc';size:9" json:"color"`
}
