package utils

type Container struct {
	pos int
	buf []byte
}

const kMaxContainerSize = 500 * 1024

func (container *Container) Get(size int) []byte {
	if kMaxContainerSize-container.pos < size {
		container.pos = 0
		container.buf = make([]byte, kMaxContainerSize)
	}
	b := container.buf[container.pos : container.pos+size]
	container.pos += size
	return b
}

func NewContainer() *Container {
	return &Container{
		buf: make([]byte, kMaxContainerSize),
	}
}
