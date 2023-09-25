package cmaptests

import "testing"
import "sync/atomic"

import "github.com/sirgallo/cmap"


//=================================== 32 bit

func TestMapOperations32(t *testing.T) {
	cMap := cmap.NewCMap[string, uint32]()

	cMap.Insert("hello", "world")
	cMap.Insert("new", "wow!")
	cMap.Insert("again", "test!")
	cMap.Insert("woah", "random entry")
	cMap.Insert("key", "Saturday!")
	cMap.Insert("sup", "6")
	cMap.Insert("final", "the!")
	cMap.Insert("6", "wow!")
	cMap.Insert("asdfasdf", "add 10")
	cMap.Insert("asdfasdf", "123123") // note same key, will update value
	cMap.Insert("asd", "queue!")
	cMap.Insert("fasdf", "interesting")
	cMap.Insert("yup", "random again!")
	cMap.Insert("asdf", "hello")
	cMap.Insert("asdffasd", "uh oh!")
	cMap.Insert("fasdfasdfasdfasdf", "error message")
	cMap.Insert("fasdfasdf", "info!")
	cMap.Insert("woah", "done")

	rootBitMap := (*cmap.CMapNode[string, uint32])(atomic.LoadPointer(&cMap.Root)).BitMap

	t.Logf("cMap after inserts")
	cMap.PrintChildren()

	expectedBitMap := uint32(542198999)
	t.Logf("actual root bitmap: %d, expected root bitmap: %d\n", rootBitMap, expectedBitMap)
	t.Logf("actual root bitmap: %032b, expected root bitmap: %032b\n", rootBitMap, expectedBitMap)
	if expectedBitMap != rootBitMap {
		t.Errorf("actual bitmap does not match expected bitmap: actual(%032b), expected(%032b)\n", rootBitMap, expectedBitMap)
	}

	t.Log("retrieve values")

	val1 := cMap.Retrieve("hello")
	expVal1 :=  "world"
	t.Logf("actual: %s, expected: %s", val1, expVal1)
	if val1 != expVal1 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)\n", val1, expVal1)
	}

	val2 := cMap.Retrieve("new")
	expVal2 :=  "wow!"
	t.Logf("actual: %s, expected: %s", val2, expVal2)
	if val2 != expVal2 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)\n", val2, expVal2)
	}

	val3 := cMap.Retrieve("asdf")
	expVal3 := "hello"
	t.Logf("actual: %s, expected: %s", val3, expVal3)
	if val3 != expVal3 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)", val3, expVal3)
	}

	val4 := cMap.Retrieve("asdfasdf")
	expVal4 := "123123"
	t.Logf("actual: %s, expected: %s", val4, expVal4)
	if val4 != expVal4 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)", val4, expVal4)
	}

	cMap.Delete("hello")
	cMap.Delete("yup")
	cMap.Delete("asdf")
	cMap.Delete("asdfasdf")
	cMap.Delete("new")
	cMap.Delete("6")

	rootBitMapAfterDelete := (*cmap.CMapNode[string, uint32])(atomic.LoadPointer(&cMap.Root)).BitMap
	t.Logf("bitmap of root after deletes: %032b\n", rootBitMapAfterDelete)
	t.Logf("bitmap of root after deletes: %d\n", rootBitMapAfterDelete)

	t.Log("hamt after deletes")
	cMap.PrintChildren()

	expectedRootBitmapAfterDelete := uint32(536956102)
	t.Log("actual bitmap:", rootBitMapAfterDelete, "expected bitmap:", expectedRootBitmapAfterDelete)
	if expectedRootBitmapAfterDelete != rootBitMapAfterDelete {
		t.Errorf("actual bitmap does not match expected bitmap: actual(%032b), expected(%032b)\n", rootBitMapAfterDelete, expectedRootBitmapAfterDelete)
	}
}


//=================================== 64 bit

func TestMapOperations64(t *testing.T) {
	cMap := cmap.NewCMap[string, uint64]()

	cMap.Insert("hello", "world")
	cMap.Insert("new", "wow!")
	cMap.Insert("again", "test!")
	cMap.Insert("woah", "random entry")
	cMap.Insert("key", "Saturday!")
	cMap.Insert("sup", "6")
	cMap.Insert("final", "the!")
	cMap.Insert("6", "wow!")
	cMap.Insert("asdfasdf", "add 10")
	cMap.Insert("asdfasdf", "123123") // note same key, will update value
	cMap.Insert("asd", "queue!")
	cMap.Insert("fasdf", "interesting")
	cMap.Insert("yup", "random again!")
	cMap.Insert("asdf", "hello")
	cMap.Insert("asdffasd", "uh oh!")
	cMap.Insert("fasdfasdfasdfasdf", "error message")
	cMap.Insert("fasdfasdf", "info!")
	cMap.Insert("woah", "done")

	rootBitMap := (*cmap.CMapNode[string, uint64])(atomic.LoadPointer(&cMap.Root)).BitMap

	t.Logf("cMap after inserts")
	cMap.PrintChildren()

	expectedBitMap := uint64(18084858599620633)
	t.Logf("actual root bitmap: %d, expected root bitmap: %d\n", rootBitMap, expectedBitMap)
	t.Logf("actual root bitmap: %032b, expected root bitmap: %032b\n", rootBitMap, expectedBitMap)
	if expectedBitMap != rootBitMap {
		t.Errorf("actual bitmap does not match expected bitmap: actual(%032b), expected(%032b)\n", rootBitMap, expectedBitMap)
	}

	t.Log("retrieve values")

	val1 := cMap.Retrieve("hello")
	expVal1 :=  "world"
	t.Logf("actual: %s, expected: %s", val1, expVal1)
	if val1 != expVal1 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)\n", val1, expVal1)
	}

	val2 := cMap.Retrieve("new")
	expVal2 :=  "wow!"
	t.Logf("actual: %s, expected: %s", val2, expVal2)
	if val2 != expVal2 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)\n", val2, expVal2)
	}

	val3 := cMap.Retrieve("asdf")
	expVal3 := "hello"
	t.Logf("actual: %s, expected: %s", val3, expVal3)
	if val3 != expVal3 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)", val3, expVal3)
	}

	val4 := cMap.Retrieve("asdfasdf")
	expVal4 := "123123"
	t.Logf("actual: %s, expected: %s", val4, expVal4)
	if val4 != expVal4 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)", val4, expVal4)
	}

	cMap.Delete("hello")
	cMap.Delete("yup")
	cMap.Delete("asdf")
	cMap.Delete("asdfasdf")
	cMap.Delete("new")
	cMap.Delete("6")

	rootBitMapAfterDelete := (*cmap.CMapNode[string, uint64])(atomic.LoadPointer(&cMap.Root)).BitMap
	t.Logf("bitmap of root after deletes: %032b\n", rootBitMapAfterDelete)
	t.Logf("bitmap of root after deletes: %d\n", rootBitMapAfterDelete)

	t.Log("hamt after deletes")
	cMap.PrintChildren()

	expectedRootBitmapAfterDelete := uint64(18014472667152401)
	t.Log("actual bitmap:", rootBitMapAfterDelete, "expected bitmap:", expectedRootBitmapAfterDelete)
	if expectedRootBitmapAfterDelete != rootBitMapAfterDelete {
		t.Errorf("actual bitmap does not match expected bitmap: actual(%032b), expected(%032b)\n", rootBitMapAfterDelete, expectedRootBitmapAfterDelete)
	}
}