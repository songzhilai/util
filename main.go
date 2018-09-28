package main

import (
	"fmt"
	"math/rand"
	"time"
)

const CNT int = 10000000

type Point struct {
	v float32
	t int64
}

var listPoint []Point

func main() {

	var f [CNT]float32
	var last_f float32 = 0

	for i := 0; i < CNT; i++ {
		f[i] = 1.1 + float32(rand.Intn(100))
	}

	t1 := time.Now()
	for i := 0; i < CNT; i++ {
		if f[i]-last_f > 30 {
			if f[i]/last_f > 2.01 {
				var pt Point
				pt.v = f[i]
				pt.t = 12
				// /listPoint = append(listPoint, pt)
				//fmt.Println(pt.v)
			}
		}

		last_f = f[i]
	}

	t2 := time.Now()
	fmt.Println((t2.UnixNano() - t1.UnixNano()) / 1e6)
	fmt.Println(len(listPoint))

}
