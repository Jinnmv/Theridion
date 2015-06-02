package StorageManager

type Storage interface {
	Write([]*interface{}) (int, error)
	Truncate() error
	Drop() error
	Close() error
}

func New() {}
