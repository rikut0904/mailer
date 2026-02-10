package entity

type ThreadGroup struct {
	ParentUUID string `json:"parent_uuid" gorm:"column:parent_uuid;primaryKey"`
	GroupName  string `json:"group_name" gorm:"column:group_name"`
}

func (ThreadGroup) TableName() string {
	return "thread_groups"
}
