package domain

type File struct {
	Name        string
	ContentType string
	Body        []byte
}

func NewFile(name string, bytes []byte, contentType string) *File {
	return &File{
		Name:        name,
		ContentType: contentType,
		Body:        bytes,
	}
}
