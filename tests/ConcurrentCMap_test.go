package cmaptests

import "bytes"
import "crypto/rand"
import "sync"
import "testing"

import "github.com/sirgallo/cmap"


type KeyVal struct {
	Key []byte
	Value []byte
}


//=================================== 32 bit

func TestMapRandomSmallConcurrentOperations32(t *testing.T) {
	cMap := cmap.NewCMap[uint32]()

	inputSize := 100000
	keyValPairs := make([]KeyVal, inputSize)

	for idx := range keyValPairs {
		randomBytes, _ := GenerateRandomBytes(32)
		keyValPairs[idx] = KeyVal{ Key: randomBytes, Value: randomBytes }
	}

	t.Log("seeded keyValPairs array:", inputSize)

	t.Log("inserting values -->")
	var insertWG sync.WaitGroup

	for _, val := range keyValPairs {
		insertWG.Add(1)
		go func (val KeyVal) {
			defer insertWG.Done()

			cMap.Put(val.Key, val.Value)
		}(val)
	}

	insertWG.Wait()

	t.Log("retrieving values -->")
	var retrieveWG sync.WaitGroup

	for _, val := range keyValPairs {
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

	t.Log("done")
}

func TestMapRandomLargeConcurrentOperations32(t *testing.T) {
	cMap := cmap.NewCMap[uint32]()

	inputSize := 10000000

	keyValPairs := make([]KeyVal, inputSize)
	keyValChan := make(chan KeyVal, inputSize)
	
	var fillArrWG sync.WaitGroup

	for range keyValPairs {
		fillArrWG.Add(1)
		go func () {
			defer fillArrWG.Done()

			randomBytes, _ := GenerateRandomBytes(32)
			keyValChan <- KeyVal{ Key: randomBytes, Value: randomBytes }
		}()
	}

	fillArrWG.Wait()
	t.Log("filled random key val pairs chan with size:", inputSize)

	for idx := range keyValPairs {
		keyVal :=<- keyValChan
		keyValPairs[idx] = keyVal
	}

	t.Log("seeded keyValPairs array:", inputSize)

	t.Log("inserting values -->")
	var insertWG sync.WaitGroup

	for _, val := range keyValPairs {
		insertWG.Add(1)
		go func (val KeyVal) {
			defer insertWG.Done()
			
			cMap.Put(val.Key, val.Value)
		}(val)
	}

	insertWG.Wait()

	t.Log("retrieving values -->")
	var retrieveWG sync.WaitGroup

	for _, val := range keyValPairs {
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

	t.Log("done")
}


//=================================== 64 bit

func TestMapRandomSmallConcurrentOperations64(t *testing.T) {
	cMap := cmap.NewCMap[uint64]()

	inputSize := 100000
	keyValPairs := make([]KeyVal, inputSize)

	for idx := range keyValPairs {
		randomBytes, _ := GenerateRandomBytes(32)
		keyValPairs[idx] = KeyVal{ Key: randomBytes, Value: randomBytes }
	}

	t.Log("seeded keyValPairs array:", inputSize)

	t.Log("inserting values -->")
	var insertWG sync.WaitGroup

	for _, val := range keyValPairs {
		insertWG.Add(1)
		go func (val KeyVal) {
			defer insertWG.Done()

			cMap.Put(val.Key, val.Value)
		}(val)
	}

	insertWG.Wait()

	t.Log("retrieving values -->")
	var retrieveWG sync.WaitGroup

	for _, val := range keyValPairs {
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

	t.Log("done")
}

func TestMapRandomLargeConcurrentOperations64(t *testing.T) {
	cMap := cmap.NewCMap[uint64]()

	inputSize := 10000000

	keyValPairs := make([]KeyVal, inputSize)
	keyValChan := make(chan KeyVal, inputSize)
	
	var fillArrWG sync.WaitGroup

	for range keyValPairs {
		fillArrWG.Add(1)
		go func () {
			defer fillArrWG.Done()

			randomBytes, _ := GenerateRandomBytes(32)
			keyValChan <- KeyVal{ Key: randomBytes, Value: randomBytes }
		}()
	}

	fillArrWG.Wait()
	t.Log("filled random key val pairs chan with size:", inputSize)

	for idx := range keyValPairs {
		keyVal :=<- keyValChan
		keyValPairs[idx] = keyVal
	}

	t.Log("seeded keyValPairs array:", inputSize)

	t.Log("inserting values -->")
	var insertWG sync.WaitGroup

	for _, val := range keyValPairs {
		insertWG.Add(1)
		go func (val KeyVal) {
			defer insertWG.Done()
			
			cMap.Put(val.Key, []byte(val.Value))
		}(val)
	}

	insertWG.Wait()

	t.Log("retrieving values -->")
	var retrieveWG sync.WaitGroup

	for _, val := range keyValPairs {
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