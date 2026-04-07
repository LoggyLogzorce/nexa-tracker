package status

type Status struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ProjectID  uint   `gorm:"not null" json:"project_id"`
	Name       string `gorm:"not null;size:50" json:"name"`
	Color      string `gorm:"default:'#cccccc';size:7" json:"color"`
	OrderIndex int    `gorm:"default:0" json:"order_index"`
}
