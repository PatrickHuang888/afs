package daemon

type BlockService interface {
	UploadBlock(b *[]byte) (string, error)

	GetBlock(bid string, b *[]byte)  error

	ObsoleteBlock(bid string) error
}
