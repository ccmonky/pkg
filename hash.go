package pkg

import "hash"

const size = 4

// JavaStringHashCode simultes the java.io.String.hashCode method
type JavaStringHashCode interface {
	hash.Hash
	HashCode() int32
}

// NewHash return the new hash code
// copy from `https://gist.github.com/giautm/d79994acd796f3065903eccbc8d6e09b`
func NewHash() JavaStringHashCode {
	var s sum32
	return &s
}

// HashCode return the hash value of bytes slice p
func HashCode(p []byte) int32 {
	h := NewHash()
	h.Write(p)
	return h.HashCode()
}

type sum32 int32

// BlockSize implement hash.Hash
func (sum32) BlockSize() int {
	return 1
}

// Size implement hash.Hash
func (sum32) Size() int {
	return size
}

// Reset implement hash.Hash
func (h *sum32) Reset() {
	*h = 0
}

// Sum implement hash.Hash
func (h sum32) Sum(in []byte) []byte {
	s := h.HashCode()
	return append(in, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}

// Write implement hash.Hash
func (h *sum32) Write(p []byte) (n int, err error) {
	s := h.HashCode()
	for _, pp := range p {
		s = 31*s + int32(pp)
	}
	*h = sum32(s)
	return len(p), nil
}

// HashCode implement the JavaStringHashCode
func (h sum32) HashCode() int32 {
	return int32(h)
}
