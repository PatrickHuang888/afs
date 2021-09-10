package daemon



type MetaStore interface {
	InsertOrAppend(key string, meta string) error
	Put(key string, meta *string) error
	Get(key string) (bids []string, err error)
}
