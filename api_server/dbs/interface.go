package dbs

// PrimaryKey 获得主键的接口, 需要主键获得主键名和主键值的方法
type PrimaryKey interface {
	GetKeyTag() string
	GetKeyValue() string
}

type MongodbData interface {
	PrimaryKey
	Dump() (err error)
	Load() (err error)
	Delete() (err error)
}

type RedisData interface {
	PrimaryKey
	Save() (err error)
	Retrieve() (err error)
	Remove() (err error)
}
