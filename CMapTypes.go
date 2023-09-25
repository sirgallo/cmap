package cmap

import "unsafe"


type CMapNode[T comparable, V uint32 | uint64] struct {
	Key        string
	Value      T
	IsLeafNode bool
	BitMap     V
	Children   []*CMapNode[T, V]
}

type CMap[T comparable, V uint32 | uint64] struct {
	BitChunkSize int
	HashChunks   int
	Root         unsafe.Pointer
}