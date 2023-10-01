package cmap

import "math"
import "sync/atomic"
import "unsafe"

import "github.com/sirgallo/cmap/utils"


func NewCMap[T comparable, V uint32 | uint64]() *CMap[T, V] {
	var bitChunkSize int

	var v V
	switch any(v).(type) {
		case uint32:
			bitChunkSize = 5
		case uint64:
			bitChunkSize = 6
	}

	hashChunks := int(math.Pow(float64(2), float64(bitChunkSize))) / bitChunkSize

	rootNode := &CMapNode[T, V]{
		IsLeafNode: false,
		BitMap:     0,
		Children:   []*CMapNode[T, V]{},
	}

	return &CMap[T, V]{
		BitChunkSize: bitChunkSize,
		HashChunks:   hashChunks,
		Root:         unsafe.Pointer(rootNode),
	}
}

func (cMap *CMap[T, V]) NewLeafNode(key string, value T) *CMapNode[T, V] {
	return &CMapNode[T, V]{
		IsLeafNode: true,
		Key:        key,
		Value:      value,
	}
}

func (cMap *CMap[T, V]) NewInternalNode() *CMapNode[T, V] {
	return &CMapNode[T, V]{
		IsLeafNode: false,
		BitMap:     0,
		Children:   []*CMapNode[T, V]{},
	}
}

func (cMap *CMap[T, V]) CopyNode(node *CMapNode[T, V]) *CMapNode[T, V] {
	nodeCopy := &CMapNode[T, V]{
		Key:        node.Key,
		Value:      node.Value,
		IsLeafNode: node.IsLeafNode,
		BitMap:     node.BitMap,
		Children:   make([]*CMapNode[T, V], len(node.Children)),
	}

	copy(nodeCopy.Children, node.Children)

	return nodeCopy
}

func (cMap *CMap[T, V]) Put(key string, value T) bool {
	for {
		completed := cMap.putRecursive(&cMap.Root, key, value, 0)
		if completed { return true }
	}
}

func (cMap *CMap[T, V]) putRecursive(node *unsafe.Pointer, key string, value T, level int) bool {
	hash := cMap.CalculateHashForCurrentLevel(key, level)
	index := cMap.getSparseIndex(hash, level)

	currNode := (*CMapNode[T, V])(atomic.LoadPointer(node))
	nodeCopy := cMap.CopyNode(currNode)

	if ! IsBitSet(nodeCopy.BitMap, index) {
		newLeaf := cMap.NewLeafNode(key, value)
		nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
		pos := cMap.getPosition(nodeCopy.BitMap, hash, level)
		nodeCopy.Children = ExtendTable(nodeCopy.Children, nodeCopy.BitMap, pos, newLeaf)

		return cMap.compareAndSwap(node, currNode, nodeCopy)
	} else {
		pos := cMap.getPosition(nodeCopy.BitMap, hash, level)
		childNode := nodeCopy.Children[pos]

		if childNode.IsLeafNode {
			if key == childNode.Key {
				nodeCopy.Children[pos].Value = value
				return cMap.compareAndSwap(node, currNode, nodeCopy)
			} else {
				newINode := cMap.NewInternalNode()
				iNodePtr := unsafe.Pointer(newINode)

				cMap.putRecursive(&iNodePtr, childNode.Key, childNode.Value, level+1)
				cMap.putRecursive(&iNodePtr, key, value, level+1)

				nodeCopy.Children[pos] = (*CMapNode[T, V])(atomic.LoadPointer(&iNodePtr))
				return cMap.compareAndSwap(node, currNode, nodeCopy)
			}
		} else {
			childPtr := unsafe.Pointer(nodeCopy.Children[pos])
			cMap.putRecursive(&childPtr, key, value, level + 1)

			nodeCopy.Children[pos] = (*CMapNode[T, V])(atomic.LoadPointer(&childPtr))
			return cMap.compareAndSwap(node, currNode, nodeCopy)
		}
	}
}

func (cMap *CMap[T, V]) Get(key string) T {
	return cMap.getRecursive(&cMap.Root, key, 0)
}

func (cMap *CMap[T, V]) getRecursive(node *unsafe.Pointer, key string, level int) T {
	hash := cMap.CalculateHashForCurrentLevel(key, level)
	index := cMap.getSparseIndex(hash, level)
	currNode := (*CMapNode[T, V])(atomic.LoadPointer(node))

	if ! IsBitSet(currNode.BitMap, index) {
		return utils.GetZero[T]()
	} else {
		pos := cMap.getPosition(currNode.BitMap, hash, level)
		childNode := currNode.Children[pos]

		if childNode.IsLeafNode && key == childNode.Key {
			if childNode.Value == (*CMapNode[T, V])(atomic.LoadPointer(node)).Children[pos].Value {
				return childNode.Value
			} else { return utils.GetZero[T]() }
		} else {
			childPtr := unsafe.Pointer(currNode.Children[pos])
			return cMap.getRecursive(&childPtr, key, level + 1)
		}
	}
}

func (cMap *CMap[T, V]) Delete(key string) bool {
	for {
		completed := cMap.deleteRecursive(&cMap.Root, key, 0)
		if completed { return true }
	}
}

func (cMap *CMap[T, V]) deleteRecursive(node *unsafe.Pointer, key string, level int) bool {
	hash := cMap.CalculateHashForCurrentLevel(key, level)
	index := cMap.getSparseIndex(hash, level)

	currNode := (*CMapNode[T, V])(atomic.LoadPointer(node))
	nodeCopy := cMap.CopyNode(currNode)

	if ! IsBitSet(nodeCopy.BitMap, index) {
		return true
	} else {
		pos := cMap.getPosition(nodeCopy.BitMap, hash, level)
		childNode := nodeCopy.Children[pos]

		if childNode.IsLeafNode {
			if key == childNode.Key {
				nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
				nodeCopy.Children = ShrinkTable(nodeCopy.Children, nodeCopy.BitMap, pos)

				return cMap.compareAndSwap(node, currNode, nodeCopy)
			}

			return false
		} else {
			childPtr := unsafe.Pointer(nodeCopy.Children[pos])
			cMap.deleteRecursive(&childPtr, key, level + 1)

			popCount := calculateHammingWeight(nodeCopy.BitMap)
			if popCount == 0 { // if empty internal node, remove from the mapped array
				nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
				nodeCopy.Children = ShrinkTable(nodeCopy.Children, nodeCopy.BitMap, pos)
			}

			return cMap.compareAndSwap(node, currNode, nodeCopy)
		}
	}
}

func (cMap *CMap[T, V]) compareAndSwap(node *unsafe.Pointer, currNode *CMapNode[T, V], nodeCopy *CMapNode[T, V]) bool {
	if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
		return true
	} else { return false }
}