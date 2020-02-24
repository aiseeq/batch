package batch

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	b := New(func(chan interface{}) {}, 1, time.Second)
	if b == nil {
		t.Error("new failed")
	}
}

func TestBatch_Add(t *testing.T) {
	type myRow struct{ field1, field2 string }
	var insertCounter int

	b := New(func(rows chan interface{}) {
		for row := range rows {
			if _, ok := row.(myRow); !ok {
				t.Error("interface conversion failed")
				return
			} else {
				insertCounter++
			}
		}
	}, 2, time.Millisecond)

	b.Add(myRow{"Because", "we add 3 rows"})
	b.Add(myRow{"And the limit is", "two"})
	ok := b.Add(myRow{"This row", "will be discarded"})

	b.Wait()
	if insertCounter != 2 {
		t.Errorf("Wrong inserted rows count. Expected: 2, got: %d", insertCounter)
	}
	if ok {
		t.Errorf("No overflow")
	}
}

func TestRaces(t *testing.T) {
	var insertCounter, okCounter int

	b := New(func(rows chan interface{}) {
		for row := range rows {
			if row.(string) != "ok" {
				t.Error("interface conversion failed")
				return
			} else {
				insertCounter++
			}
		}
		return
	}, 1000, 0*time.Second)
	for x := 0; x < 100000; x++ {
		if b.Add("ok") {
			okCounter++
		}
	}
	b.Wait()
	if insertCounter != okCounter {
		t.Error("Wrong inserted rows count", insertCounter)
	}
}

func BenchmarkAdd(bm *testing.B) {
	bm.StopTimer()

	type myRow struct{ field1, field2 string }

	b := New(func(rows chan interface{}) {
		for row := range rows {
			_ = row.(*myRow)
		}
	}, 10000, time.Second)
	row := myRow{}

	bm.StartTimer()
	for x := 0; x < bm.N; x++ {
		b.Add(&row)
	}
}

func ExampleBatch() {
	type myRow struct{ field1, field2 string }

	b := New(func(rows chan interface{}) {
		// Here you describe elements processing
		for row := range rows {
			if mr, ok := row.(myRow); !ok {
				fmt.Println("Error: Interface conversion failed!")
				return
			} else {
				fmt.Println(mr.field1 + ", " + mr.field2)
			}
		}
	},
		// Maximum batch size, maximum queue length. Determines how many elements could be stored before processing
		// If there will be more than that, new elements will be discarded
		// You should select reasonable number keeping in mind, that queue will reside in memory
		1000000,
		// flush() callback period
		time.Second)
	b.Add(myRow{"Row 1, Field 1", "Row 1, Field 2"})
	b.Add(myRow{"Row 2, Field 1", "Row 2, Field 2"})
	b.Wait() // Call it before application termination
}
