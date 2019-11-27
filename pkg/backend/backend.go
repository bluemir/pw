package backend

type Backend struct {
	invFilePath string
}

func New(invFilePath string) (*Backend, error) {
	return &Backend{invFilePath}, nil
}
