package cmaptests

import "testing"

import "github.com/sirgallo/cmap"


//=================================== 32 bit

func TestMurmur32(t *testing.T) {
	key := "hello"
	seed := uint32(1)

	hash := cmap.Murmur32(key, seed)
	t.Log("hash:", hash)
}

func TestMurmur32ReSeed(t *testing.T) {
	key := "hello"
	levels := make([]int, 17)
	totalLevels := 6
	chunkSize := 5

	cMap := cmap.NewCMap[string, uint32]()

	for idx := range levels {
		hash := cMap.CalculateHashForCurrentLevel(key, idx)
		index := cmap.GetIndexForLevel(hash, chunkSize, idx, totalLevels)
		t.Logf("hash: %d, index: %d", hash, index)
	}
}


//=================================== 64 bit

func TestMurmur64(t *testing.T) {
	key := "hello"
	seed := uint64(1)

	hash := cmap.Murmur64(key, seed)
	t.Log("hash:", hash)
}

func TestMurmur64ReSeed(t *testing.T) {
	key := "hello"
	levels := make([]int, 33)
	totalLevels := 10
	chunkSize := 6

	cMap := cmap.NewCMap[string, uint64]()

	for idx := range levels {
		hash := cMap.CalculateHashForCurrentLevel(key, idx)
		index := cmap.GetIndexForLevel(hash, chunkSize, idx, totalLevels)
		t.Logf("hash: %d, index: %d", hash, index)
	}
}