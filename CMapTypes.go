package cmap

import "unsafe"


// CMapNode represents a singular node within the hash array mapped trie data structure. Can be either 32 or 64 bits
type CMapNode[T uint32 | uint64] struct {
	// Key: The key associated with a value. Keys are in byte array representation. Keys are only stored within leaf nodes of the hamt
	Key	[]byte
	// Value: The value associated with a key, in byte array representation. Values are only stored within leaf nodes
	Value []byte
	// IsLeaf: flag indicating if the current node is a leaf node or an internal node
	IsLeaf bool
	// Bitmap: a 32 bit sparse index that indicates the location of each hashed key within the array of child nodes. Only stored in internal nodes
	Bitmap T
	// Children: an array of child nodes, which are CMapNodes. Location in the array is determined by the sparse index
	Children []*CMapNode[T]
}

// CMap is the root of the hash array mapped trie
type CMap[T uint32 | uint64] struct {
	// Root: the root CMapNode within the hash array mapped trie. Stored as a pointer to the location in memory of the root
	Root unsafe.Pointer
	// BitChunkSize: the size of each chunk in the 32 bit or 64 bit hash. Example, with a 32 bit hash total size is 2^5, so each chunk will be 5 bits long
	BitChunkSize int
	// HashChunks: the total chunks of the 32 bit or 64 bit hash determining the levels within the hash array mapped trie
	HashChunks int
}