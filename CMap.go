package cmap

import "bytes"
import "math"
import "sync/atomic"
import "unsafe"


/*
NewCMap initializes a new hash array mapped trie

Returns:

The newly initialized hash array mapped trie
*/
func NewCMap[T uint32 | uint64]() *CMap[T] {
	var bitChunkSize int

	var t T
	switch any(t).(type) {
		case uint32:
			bitChunkSize = 5
		case uint64:
			bitChunkSize = 6
	}

	hashChunks := int(math.Pow(float64(2), float64(bitChunkSize))) / bitChunkSize

	rootNode := &CMapNode[T]{
		IsLeaf: false,
		Bitmap: 0,
		Children: []*CMapNode[T]{},
	}

	return &CMap[T]{
		BitChunkSize: bitChunkSize,
		HashChunks: hashChunks,
		Root: unsafe.Pointer(rootNode),
	}
}

/*
NewLeafNode creates a new leaf node in the hash array mapped trie, which stores a key value pair

Parameters:

- key: the incoming key to be inserted

- value: the incoming value associated with the key

Returns:

A new leaf node in the hash array mapped trie
*/
func (cMap *CMap[T]) NewLeafNode(key []byte, value []byte) *CMapNode[T] {
	return &CMapNode[T]{
		IsLeaf: true,
		Key: key,
		Value: value,
	}
}

/*
NewInternalNode creates a new internal node in the hash array mapped trie, which is essentially a branch node that contains pointers to child nodes

Returns:

A new internal node with bitmap initialized to 0 and an empty array of child nodes
*/
func (cMap *CMap[T]) NewInternalNode() *CMapNode[T] {
	return &CMapNode[T]{
		IsLeaf: false,
		Bitmap: 0,
		Children: []*CMapNode[T]{},
	}
}

/*
CopyNode creates a copy of an existing node. This is used for path copying, so on operations that modify the trie, a copy is created instead of modifying the existing node. This
makes the data structure essentially immutable. If an operation succeeds, the copy replaces the existing node, otherwise the copy is discarded.

Parameters:

- node: the existing node to create a copy of

Returns:

A copy of the existing node within the hash array mapped trie, which the operation will modify
*/
func (cMap *CMap[T]) CopyNode(node *CMapNode[T]) *CMapNode[T] {
	nodeCopy := &CMapNode[T]{
		Key: node.Key,
		Value: node.Value,
		IsLeaf: node.IsLeaf,
		Bitmap: node.Bitmap,
		Children: make([]*CMapNode[T], len(node.Children)),
	}

	copy(nodeCopy.Children, node.Children)

	return nodeCopy
}

/*
Put inserts or updates key-value pair into the hash array mapped trie. The operation begins at the root of the trie and traverses through the tree until the correct location is found. If the 
operation fails, the copied and modified path is discarded and the operation retries back at the root until completed.

Parameters:

- key: the key in the key-value pair

- value: the value in the key-value pair

Returns:

truthy on successful completion
*/
func (cMap *CMap[T]) Put(key []byte, value []byte) bool {
	for {
		completed := cMap.putRecursive(&cMap.Root, key, value, 0)
		if completed { return true }
	}
}

/*
putRecursive attempts to traverse through the trie, locating the node at a given level to modify for the key-value pair. It first hashes the key, determines the sparse index in the bitmap to modify, 
and createsa copy of the current node to be modified. If the bit in the bitmap of the node is not set, a new leaf node is created, the bitmap of the copy is modified to reflect the position of the 
new leaf node, and the child node array is extended to include the new leaf node. Then, an atomic compare and swap operation is performed where the operation attempts to replace the current node with 
the modified copy. If the operation succeeds the response is returned by moving back up the tree. If it fails, the copy is discarded and the operation returns to the root to be reattempted. If the 
current bit is set in the bitmap, the operation checks if the node at the location in the child node array is a leaf node or an internal node. If it is a leaf node and the key is the same as the 
incoming key, the copy is modified with the new value and we attempt to compare and swap the current child leaf node with the new copy. If the leaf node does not contain the same key, the operation creates 
a new internal node, and inserts the new leaf node for the incoming key and value as well as the existing child node into the new internal node, and attempts to compare and swap the current leaf node 
with the new internal node containing the existing child node and the new leaf node for the incoming key and value. If the node is an internal node, the operation traverses down the tree to the internal
node and the above steps are repeated until the key-value pair is inserted.

Parameters:

- node: the node in the tree where the operation is currently

- key: the key for the incoming key-value pair

- value: the value for the incoming key-value pair

- level: the current level in the tree the operation is at

Returns:

truthy value from successful or failed compare and swap operations
*/
func (cMap *CMap[T]) putRecursive(node *unsafe.Pointer, key []byte, value []byte, level int) bool {
	hash := cMap.CalculateHashForCurrentLevel(key, level)
	index := cMap.getSparseIndex(hash, level)

	currNode := (*CMapNode[T])(atomic.LoadPointer(node))
	nodeCopy := cMap.CopyNode(currNode)

	if ! IsBitSet(nodeCopy.Bitmap, index) {
		newLeaf := cMap.NewLeafNode(key, value)
		nodeCopy.Bitmap = SetBit(nodeCopy.Bitmap, index)
		pos := cMap.getPosition(nodeCopy.Bitmap, hash, level)
		nodeCopy.Children = ExtendTable(nodeCopy.Children, nodeCopy.Bitmap, pos, newLeaf)

		return cMap.compareAndSwap(node, currNode, nodeCopy)
	} else {
		pos := cMap.getPosition(nodeCopy.Bitmap, hash, level)
		childNode := nodeCopy.Children[pos]

		if childNode.IsLeaf {
			if bytes.Equal(key, childNode.Key) {
				nodeCopy.Children[pos].Value = value
				return cMap.compareAndSwap(node, currNode, nodeCopy)
			} else {
				newINode := cMap.NewInternalNode()
				iNodePtr := unsafe.Pointer(newINode)

				cMap.putRecursive(&iNodePtr, []byte(childNode.Key), childNode.Value, level + 1)
				cMap.putRecursive(&iNodePtr, key, value, level + 1)

				nodeCopy.Children[pos] = (*CMapNode[T])(atomic.LoadPointer(&iNodePtr))
				return cMap.compareAndSwap(node, currNode, nodeCopy)
			}
		} else {
			childPtr := unsafe.Pointer(nodeCopy.Children[pos])
			cMap.putRecursive(&childPtr, key, value, level + 1)

			nodeCopy.Children[pos] = (*CMapNode[T])(atomic.LoadPointer(&childPtr))
			return cMap.compareAndSwap(node, currNode, nodeCopy)
		}
	}
}

/*
Get attempts to retrieve the value for a key within the hash array mapped trie. The operation begins at the root of the trie and traverses down the path to the key.

Returns:

either the value for the key in byte array representation or nil if the key does not exist
*/
func (cMap *CMap[T]) Get(key []byte) []byte {
	return cMap.getRecursive(&cMap.Root, key, 0)
}

/*
getRecursive attempts to recursively retrieve a value for a given key within the hash array mapped trie. For each node traversed to at each level the operation travels to, the sparse
index is calculated for the hashed key. If the bit is not set in the bitmap, return nil since the key has not been inserted yet into the trie. Otherwise, determine the position in the child node
array for the sparse index. If the child node is a leaf node and the key to be searched for is the same as the key of the child node, the value has been found. Since the trie utilizes path copying, 
any threads modifying the trie are modifying copies so it the get operation returns the value at the point in time of the get operation. If the node is node a leaf node, but instead an internal node, 
recurse down the path to the next level to the child node in the position of the child node array and repeat the above.

Parameters:

- node: the pointer to the node to be checked for the key-value pair

- key: the key being searched for

- level: the current level within the trie the operation is at

Returns:

either the value for the given key or nil if non-existent or if the node is being modified
*/
func (cMap *CMap[T]) getRecursive(node *unsafe.Pointer, key []byte, level int) []byte {
	hash := cMap.CalculateHashForCurrentLevel(key, level)
	index := cMap.getSparseIndex(hash, level)
	currNode := (*CMapNode[T])(atomic.LoadPointer(node))

	if ! IsBitSet(currNode.Bitmap, index) {
		return nil
	} else {
		pos := cMap.getPosition(currNode.Bitmap, hash, level)
		childNode := currNode.Children[pos]

		if childNode.IsLeaf && bytes.Equal(key, childNode.Key) {
			return childNode.Value
		} else {
			childPtr := unsafe.Pointer(currNode.Children[pos])
			return cMap.getRecursive(&childPtr, key, level + 1)
		}
	}
}

/*
Delete attempts to delete a key-value pair within the hash array mapped trie. It starts at the root of the trie and recurses down the path to the key to be deleted. If the operation succeeds truthy value 
is returned, otherwise the operation returns to the root to retry the operation.

Parameters:

- key: the key to attempt to delete

Returns:

truthy on successful completion
*/
func (cMap *CMap[T]) Delete(key []byte) bool {
	for {
		completed := cMap.deleteRecursive(&cMap.Root, key, 0)
		if completed { return true }
	}
}

/*
deleteRecursive attempts to recursively move down the path of the trie to the key-value pair to be deleted. The hash for the key is calculated, the sparse index in the bitmap is determined 
for the given level, and a copy of the current node is created to be modifed. If the bit in the bitmap is not set, the key doesn't exist so truthy is returned since there is nothing to delete and the
operation completes. If the bit is set, the child node for the position within the child node array is found. If the child node is a leaf node and the key of the child node is equal to the 
key of the key to delete, the copy is modified to update the bitmap and shrink the table and remove the given node. A compare and swap operation is performed, and if successful traverse back up the trie
and complete, otherwise the operation is returned to the root to retry. If the child node is an internal node, the operation recurses down the trie to the next level. On return, if the internal node 
is empty, the copy modified so the bitmap is updated and table is shrunk. A compare and swap operation is performed on the current node with the new copy.

Parameters:

- node: a pointer to the node that is being modified

- key: the key to be deleted

- level: the current level within the trie

Returns:

truthy response on success and falsey on failure
*/
func (cMap *CMap[T]) deleteRecursive(node *unsafe.Pointer, key []byte, level int) bool {
	hash := cMap.CalculateHashForCurrentLevel(key, level)
	index := cMap.getSparseIndex(hash, level)

	currNode := (*CMapNode[T])(atomic.LoadPointer(node))
	nodeCopy := cMap.CopyNode(currNode)

	if ! IsBitSet(nodeCopy.Bitmap, index) {
		return true
	} else {
		pos := cMap.getPosition(nodeCopy.Bitmap, hash, level)
		childNode := nodeCopy.Children[pos]

		if childNode.IsLeaf {
			if bytes.Equal(key, childNode.Key) {
				nodeCopy.Bitmap = SetBit(nodeCopy.Bitmap, index)
				nodeCopy.Children = ShrinkTable(nodeCopy.Children, nodeCopy.Bitmap, pos)

				return cMap.compareAndSwap(node, currNode, nodeCopy)
			}

			return false
		} else {
			childPtr := unsafe.Pointer(nodeCopy.Children[pos])
			cMap.deleteRecursive(&childPtr, key, level + 1)

			popCount := calculateHammingWeight(nodeCopy.Bitmap)
			if popCount == 0 { // if empty internal node, remove from the mapped array
				nodeCopy.Bitmap = SetBit(nodeCopy.Bitmap, index)
				nodeCopy.Children = ShrinkTable(nodeCopy.Children, nodeCopy.Bitmap, pos)
			}

			return cMap.compareAndSwap(node, currNode, nodeCopy)
		}
	}
}

/*
compareAndSwap performs CAS opertion

Parameters:

- node: the node to be updated

- currNode: the original node to be updated

- nodeCopy: the copy to swap the original with

Returns:

true on successful CAS and false on failure
*/
func (cMap *CMap[T]) compareAndSwap(node *unsafe.Pointer, currNode *CMapNode[T], nodeCopy *CMapNode[T]) bool {
	if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
		return true
	} else { return false }
}