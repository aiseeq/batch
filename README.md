Batch processing library
---

Simple library for delayed batch processing

Example
---

```go
package main

import (
    "fmt"
    "github.com/aiseeq/batch"
    "time"
)

type myRow struct{ field1, field2 string }

func main() {
    b := batch.New(func(rows chan interface{}) {
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
```