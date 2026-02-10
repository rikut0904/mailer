package repository

type MailStorageRepository interface {
	ListKeys(prefix string, continuationToken *string, maxKeys int) (keys []string, nextToken *string, err error)
	GetObject(key string) ([]byte, error)
	DeleteObject(key string) error
}
