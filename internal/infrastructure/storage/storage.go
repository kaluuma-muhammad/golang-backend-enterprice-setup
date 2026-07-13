package storage

type Storage interface {
	Upload(input UploadInput) (*File, error)
	Delete(path string) error
	GetURL(path string) string
}
