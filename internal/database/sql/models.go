package sql

type Url struct {
	Id      int64 `gorm:"primaryKey;AUTO_INCREMENT"`
	FullUrl string
}
