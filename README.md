# util
golang 
sdt(旋转门压缩算法)

```
package main

import (
	"fmt"

	"github.com/songzhilai/util/util/sdt"
)

func main() {

	sdt := sdt.NewSdt("1")
	sdt.CalculateE(1, 500.5, 1000)
	for i := 1; i <= 1000; i++ {
		// if i >= 10 {
		// 	sdt.InputData(fmt.Sprintf("%d.1", i), float64(i))
		// } else {
		sdt.InputData(fmt.Sprintf("%d.1", i), float64(i))
		// }
	}
	fmt.Println(sdt)
	fmt.Println(sdt.OutputData())
}
```
