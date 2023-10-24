package cmaptests

import "bytes"
import "fmt"
import "crypto/rand"
import "sync"
import "testing"

import "github.com/sirgallo/cmap"


type KeyVal struct {
	Key []byte
	Value []byte
}


var sKeyValPairs, lKeyValPairs []KeyVal
var lKeyValChan chan KeyVal
var fillArrWG sync.WaitGroup


func init() {
	sInputSize := 100000
	sKeyValPairs = make([]KeyVal, sInputSize)

	for idx := range sKeyValPairs {
		randomBytes, _ := GenerateRandomBytes(32)
		sKeyValPairs[idx] = KeyVal{ Key: randomBytes, Value: randomBytes }
	}

	lInputSize := 10000000
	lKeyValPairs = make([]KeyVal, lInputSize)
	lKeyValChan = make(chan KeyVal, lInputSize)

	for range lKeyValPairs {
		fillArrWG.Add(1)
		go func () {
			defer fillArrWG.Done()

			randomBytes, _ := GenerateRandomBytes(32)
			lKeyValChan <- KeyVal{ Key: randomBytes, Value: randomBytes }
		}()
	}

	fillArrWG.Wait()
	fmt.Println("filled random key val pairs chan with size:", lInputSize)

	for idx := range lKeyValPairs {
		keyVal :=<- lKeyValChan
		lKeyValPairs[idx] = keyVal
	}
}


//=================================== 32 bit


func TestCMapSmallConcurrentOps32(t *testing.T) {
	cMap := cmap.NewCMap[uint32]()
	
	t.Run("test insert", func(t *testing.T) {
		var insertWG sync.WaitGroup

		for _, val := range sKeyValPairs {
			insertWG.Add(1)
			go func (val KeyVal) {
				defer insertWG.Done()
	
				cMap.Put(val.Key, val.Value)
			}(val)
		}

		insertWG.Wait()
	})

	t.Run("test retrieve", func(t *testing.T) {
		var retrieveWG sync.WaitGroup

		for _, val := range sKeyValPairs {
			retrieveWG.Add(1)
			go func (val KeyVal) {
				defer retrieveWG.Done()
	
				value := cMap.Get(val.Key)
				// t.Logf("actual: %s, expected: %s", value, val.Value)
				if ! bytes.Equal(value, val.Value) {
					t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", value, val.Value)
				}
			}(val)
		}
	
		retrieveWG.Wait()
	})

	t.Run("test delete", func(t *testing.T) {
		var delWG sync.WaitGroup

		for _, val := range sKeyValPairs {
			delWG.Add(1)
			go func (val KeyVal) {
				defer delWG.Done()
	
				cMap.Delete(val.Key)
			}(val)
		}

		delWG.Wait()
	})

	t.Log("done")
}

func TestCMapLargeConcurrentOps32(t *testing.T) {
	cMap := cmap.NewCMap[uint32]()
	
	t.Run("test insert", func(t *testing.T) {
		var insertWG sync.WaitGroup

		for _, val := range lKeyValPairs {
			insertWG.Add(1)
			go func (val KeyVal) {
				defer insertWG.Done()
				
				cMap.Put(val.Key, val.Value)
			}(val)
		}

		insertWG.Wait()
	})

	t.Run("test retrieve", func(t *testing.T) {
		var retrieveWG sync.WaitGroup

		for _, val := range lKeyValPairs {
			retrieveWG.Add(1)
			go func (val KeyVal) {
				defer retrieveWG.Done()
	
				value := cMap.Get(val.Key)
				// t.Logf("actual: %s, expected: %s", value, val.Value)
				if ! bytes.Equal(value, val.Value) {
					t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", value, val.Value)
				}
			}(val)
		}
	
		retrieveWG.Wait()
	})

	t.Run("test delete", func(t *testing.T) {
		var delWG sync.WaitGroup

		for _, val := range lKeyValPairs {
			delWG.Add(1)
			go func (val KeyVal) {
				defer delWG.Done()
	
				cMap.Delete(val.Key)
			}(val)
		}

		delWG.Wait()
	})

	t.Log("done")
}


//=================================== 64 bit


func TestCMapSmallConcurrentOps64(t *testing.T) {
	cMap := cmap.NewCMap[uint64]()

	t.Run("test insert", func(t *testing.T) {
		var insertWG sync.WaitGroup

		for _, val := range sKeyValPairs {
			insertWG.Add(1)
			go func (val KeyVal) {
				defer insertWG.Done()
	
				cMap.Put(val.Key, val.Value)
			}(val)
		}

		insertWG.Wait()
	})

	t.Run("test retrieve", func(t *testing.T) {
		var retrieveWG sync.WaitGroup

		for _, val := range sKeyValPairs {
			retrieveWG.Add(1)
			go func (val KeyVal) {
				defer retrieveWG.Done()
	
				value := cMap.Get(val.Key)
				// t.Logf("actual: %s, expected: %s", value, val.Value)
				if ! bytes.Equal(value, val.Value) {
					t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", value, val.Value)
				}
			}(val)
		}
	
		retrieveWG.Wait()
	})

	t.Run("test delete", func(t *testing.T) {
		var delWG sync.WaitGroup

		for _, val := range sKeyValPairs {
			delWG.Add(1)
			go func (val KeyVal) {
				defer delWG.Done()
	
				cMap.Delete(val.Key)
			}(val)
		}

		delWG.Wait()
	})

	t.Log("done")
}

func TestCMapLargeConcurrentOps64(t *testing.T) {
	cMap := cmap.NewCMap[uint64]()

	t.Run("test insert", func(t *testing.T) {
		var insertWG sync.WaitGroup

		for _, val := range lKeyValPairs {
			insertWG.Add(1)
			go func (val KeyVal) {
				defer insertWG.Done()
				
				cMap.Put(val.Key, val.Value)
			}(val)
		}

		insertWG.Wait()
	})

	t.Run("test retrieve", func(t *testing.T) {
		var retrieveWG sync.WaitGroup

		for _, val := range lKeyValPairs {
			retrieveWG.Add(1)
			go func (val KeyVal) {
				defer retrieveWG.Done()
	
				value := cMap.Get(val.Key)
				// t.Logf("actual: %s, expected: %s", value, val.Value)
				if ! bytes.Equal(value, val.Value) {
					t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", value, val.Value)
				}
			}(val)
		}
	
		retrieveWG.Wait()
	})

	t.Run("test delete", func(t *testing.T) {
		var delWG sync.WaitGroup

		for _, val := range lKeyValPairs {
			delWG.Add(1)
			go func (val KeyVal) {
				defer delWG.Done()
	
				cMap.Delete(val.Key)
			}(val)
		}

		delWG.Wait()
	})

	t.Log("done")
}


//=================================== helper

func GenerateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	return randomBytes, nil
}