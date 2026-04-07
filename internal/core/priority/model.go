package priority

type Priority struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProjectID uint   `gorm:"not null" json:"project_id"`
	Title     string `gorm:"not null;size:50" json:"title"`
	Color     string `gorm:"default:'#cccccc';size:7" json:"color"`
}
