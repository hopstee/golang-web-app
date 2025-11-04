package filestorage

type FileStorage interface {
	SaveFile(name string, data []byte) (string, error)
	DeleteFile(path string) error
}
