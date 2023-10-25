package cmaptests

import "bytes"
import "sync"
import "testing"

import "github.com/sirgallo/cmap"


var parallelTestMap *cmap.CMap[uint32]
var pInputSize int
var initKeyValPairs []KeyVal
var pKeyValPairs []KeyVal


func init() {
	parallelTestMap = cmap.NewCMap[uint32]()

	pInputSize = 100000
	initKeyValPairs = make([]KeyVal, pInputSize)
	pKeyValPairs = make([]KeyVal, pInputSize)

	for idx := range initKeyValPairs {
		iRandomBytes, _ := GenerateRandomBytes(32)
		initKeyValPairs[idx] = KeyVal{ Key: iRandomBytes, Value: iRandomBytes }
	}

	for idx := range pKeyValPairs {
		pRandomBytes, _ := GenerateRandomBytes(32)
		pKeyValPairs[idx] = KeyVal{ Key: pRandomBytes, Value: pRandomBytes }
	}

	var initMapWG sync.WaitGroup

	for _, val := range initKeyValPairs {
		initMapWG.Add(1)
		go func(val KeyVal) {
			defer initMapWG.Done()

			parallelTestMap.Put(val.Key, val.Value)
		}(val)
	}

	initMapWG.Wait()
}


func TestParallelReadWrites(t *testing.T) {
	t.Run("test init key val pairs in map", func(t *testing.T) {
		t.Parallel()

		var retrieveWG sync.WaitGroup

		for _, val := range initKeyValPairs {
			retrieveWG.Add(1)
			go func(val KeyVal) {
				defer retrieveWG.Done()

				value := parallelTestMap.Get(val.Key)

				if ! bytes.Equal(value, val.Value) {
					t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", value, val.Value)
				}
			}(val)
		}

		retrieveWG.Wait()
	})

	t.Run("test write new key vals in map", func(t *testing.T) {
		t.Parallel()

		var insertWG sync.WaitGroup

		for _, val := range pKeyValPairs {
			insertWG.Add(1)
			go func(val KeyVal) {
				defer insertWG.Done()

				parallelTestMap.Put(val.Key, val.Value)
			}(val)
		}

		insertWG.Wait()
	})

	t.Log("Done")
}