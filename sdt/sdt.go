package sdt

import (
	"math"
)

var MapSdt map[string]*Sdt

const (
	MAX_FLOAT         = 3.40282e+038
	COMPRESS_MIN      = 10
	COMPRESS_MAX      = 100
	COMPRESS_OUT_DATA = 50
)

//Sdt 旋转门初始化类
type Sdt struct {
	ListPoint []Point
	E         float64
	upGate    float64
	downGate  float64
	nowUp     float64
	nowDown   float64

	currentData    Point
	lastReadData   Point
	lastStoredData Point

	lastStoredT int
	currentT    int

	listOutPoint   []Point
	listFirstPoint []Point

	id                string
	isFirstPoint      bool
	compressIndex     int
	secondCompressSdt *Sdt

	compressMin        int
	compressMax        int
	compressOutDataMax int
}

type Point struct {
	V float64
	T string
}

//NewSdt 构造函数 初始化SDT类
func NewSdt(id string) *Sdt {
	var result = new(Sdt)
	result.compressIndex = 0
	result.id = id
	result.upGate = MAX_FLOAT
	result.downGate = -MAX_FLOAT
	result.currentT = 0
	result.lastStoredT = 0
	result.isFirstPoint = false
	result.secondCompressSdt = nil
	result.compressMin = COMPRESS_MIN
	result.compressMax = COMPRESS_MAX
	result.compressOutDataMax = COMPRESS_OUT_DATA
	return result
}

//NewSdtIndex 构造函数 初始化SDT类
func NewSdtIndex(id string, compressIndex int) *Sdt {
	var result = new(Sdt)
	result.compressIndex = compressIndex
	result.id = id
	result.upGate = MAX_FLOAT
	result.downGate = -MAX_FLOAT
	result.currentT = 0
	result.lastStoredT = 0
	result.isFirstPoint = false
	result.secondCompressSdt = nil
	result.compressMin = COMPRESS_MIN
	result.compressMax = COMPRESS_MAX
	result.compressOutDataMax = COMPRESS_OUT_DATA
	return result
}

//CalculateE 设置阀值
func (s *Sdt) CalculateE(minv, avgv, maxv float64) {
	max := maxv - avgv
	min := avgv - minv
	minE := s.sizeComparison(max, min, "min")
	s.E = s.round(minE/10.0, 10)
	if s.E < 0.00001 {
		s.E = 0.1
	}
	s.isFirstPoint = true
}

//SetE 设置阀值
func (s *Sdt) SetE(e float64) {
	s.E = e
	s.isFirstPoint = true
}

//InputData 添加需要压缩的点
func (s *Sdt) InputData(t string, v float64) {
	if len(s.listFirstPoint) < s.compressMin {
		p := Point{v, t}
		s.listFirstPoint = append(s.listFirstPoint, p)
	}
	if s.isFirstPoint {
		s.lastStoredData.T = t
		s.lastStoredData.V = v
		s.lastReadData = s.lastStoredData
		s.currentT = 0
		s.lastStoredT = 0
		s.isFirstPoint = false
	}
	s.compress(t, v)
}

func (s *Sdt) compress(t string, v float64) {
	s.currentT++
	s.nowUp = (v - s.lastStoredData.V - s.E) / float64(s.currentT-s.lastStoredT)
	if s.nowUp > s.upGate {
		s.upGate = s.nowUp
	}
	s.nowDown = (v - s.lastStoredData.V + s.E) / float64(s.currentT-s.lastStoredT)
	if s.nowDown > s.downGate {
		s.downGate = s.nowDown
	}
	if s.upGate >= s.downGate {
		s.listOutPoint = append(s.listOutPoint, s.lastReadData) // 保存前一个点
		s.lastStoredT = s.currentT - 1                          // 修改最近保存数据时间点
		s.lastStoredData = s.lastReadData

		//初始化两扇门为当前点与上个点的斜率
		s.upGate = (v - s.lastStoredData.V - s.E)
		s.downGate = (v - s.lastStoredData.V + s.E)
	}
	s.lastReadData.T = t
	s.lastReadData.V = v
}

//OutputData 返回压缩后的的结构
func (s *Sdt) OutputData() *[]Point {
	s.listOutPoint = append(s.listOutPoint, s.lastReadData) // 保存最后一个点
	if len(s.listOutPoint) < s.compressMin {
		return &s.listFirstPoint
	} else {
		if len(s.listOutPoint) > s.compressMax {
			if s.compressIndex > 10 { // 限定10次
				step := len(s.listOutPoint) / s.compressOutDataMax
				// for k, _ := range s.listOutPoint {
				for i := 0; i < len(s.listOutPoint); i++ {
					if i%step != 0 {
						s.listOutPoint = append(s.listOutPoint[0:i], s.listOutPoint[i+1:]...)
					}
				}

			} else {
				s.secondCompressSdt = NewSdtIndex(s.id, s.compressIndex+1)
				s.secondCompressSdt.SetE(s.E * float64(len(s.listOutPoint)) / float64(s.compressOutDataMax))
				for _, v := range s.listOutPoint {
					s.secondCompressSdt.InputData(v.T, v.V)
				}
				return s.secondCompressSdt.OutputData()
			}
		}
	}
	return &s.listOutPoint
}

//float大小比较替换
func (s *Sdt) sizeComparison(max, min float64, reuval string) float64 {
	if max < min {
		temp := min
		min = max
		max = temp
	}
	if reuval == "min" {
		return min
	} else if reuval == "max" {
		return max
	}
	return 0
}

//float精确位数
func (s *Sdt) round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}
