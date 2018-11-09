package sdt

import (
	"math"
)

//MapSdt sdt map
var MapSdt = make(map[string]*Sdt)

//MapSdtClass 生成sdt类
func MapSdtClass(id string) *Sdt {
	var sdt *Sdt
	if v, ok := MapSdt[id]; ok {
		//存在
		sdt = v
	} else {
		//不存在创建
		MapSdt[id] = NewSdt(id)
		sdt = MapSdt[id]
	}
	return sdt
}

const (
	MAX_FLOAT = 3.40282e+038
	// COMPRESS_MIN      = 100
	// COMPRESS_OUT_DATA = 200
	// COMPRESS_NUMBER   = 10
	COMPRESS_MIN      = 5
	COMPRESS_OUT_DATA = 12
	COMPRESS_NUMBER   = 10
)

//Sdt 旋转门初始化类
type Sdt struct {
	// ListPoint []Point
	E        float64 // 压缩阈值
	upGate   float64
	downGate float64
	nowUp    float64 // 当前数据的上斜率
	nowDown  float64 // 当前数据的下斜率

	currentData    Point // 当前数据
	lastReadData   Point // 当前数据的前一个数据
	lastStoredData Point // 最近保存的点

	lastStoredT int64 // 最近保存数据的时间刻度
	currentT    int64

	id                string // 数据点ID
	isFirstPoint      bool
	compressIndex     int
	secondCompressSdt *Sdt

	compressMin        int
	compressMax        int
	compressOutDataMax int

	listOutPoint []Point

	listFirstPoint []Point

	firstE float64 //第一次阀值
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
	result.compressOutDataMax = COMPRESS_OUT_DATA
	return result
}

//NewSdtIndex 构造函数 初始化SDT类
func NewSdtIndex(id string, compressIndex int, firstE float64) *Sdt {
	var result = new(Sdt)
	result.firstE = firstE
	result.compressIndex = compressIndex
	result.id = id
	result.upGate = MAX_FLOAT
	result.downGate = -MAX_FLOAT
	result.currentT = 0
	result.lastStoredT = 0
	result.isFirstPoint = false
	result.secondCompressSdt = nil
	result.compressMin = COMPRESS_MIN
	result.compressOutDataMax = COMPRESS_OUT_DATA
	return result
}

//CalculateE 设置阀值
func (s *Sdt) CalculateE(minv, avgv, maxv float64) {
	minE := s.sizeComparison((maxv - avgv), (avgv - minv), "min")
	s.E = s.round(minE/10.0, 10)
	if s.E < 0.00001 {
		s.E = 0.1
	}
	s.isFirstPoint = true
	s.firstE = s.E
}

//SetE 设置阀值
func (s *Sdt) SetE(e float64) {
	s.E = e
	s.isFirstPoint = true
}

//InputData 添加需要压缩的点
func (s *Sdt) InputData(t string, v float64) {
	if len(s.listFirstPoint) < s.compressMin {
		s.listFirstPoint = append(s.listFirstPoint, Point{v, t})
	}
	if s.isFirstPoint {
		s.lastStoredData = Point{v, t}
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
	if s.nowUp > s.upGate { //上升趋势，越大门开得越大
		s.upGate = s.nowUp
	}
	s.nowDown = (v - s.lastStoredData.V + s.E) / float64(s.currentT-s.lastStoredT)
	if s.nowDown < s.downGate { //下降趋势，越小门开得越大
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
	if len(s.listFirstPoint) < s.compressMin {
		return &s.listFirstPoint
	} else {
		if len(s.listOutPoint) > s.compressOutDataMax {
			if s.compressIndex > COMPRESS_NUMBER { // 限定10次
				step := len(s.listOutPoint) / s.compressOutDataMax
				for k, v := range s.listOutPoint {
					if k == 0 {
						s.listOutPoint = []Point{}
					}
					if k%step == 0 {
						s.listOutPoint = append(s.listOutPoint, v)
					}
				}

			} else {
				s.secondCompressSdt = NewSdtIndex(s.id, s.compressIndex+1, s.firstE)
				s.secondCompressSdt.SetE(s.firstE * float64(s.secondCompressSdt.compressIndex))
				// fmt.Println(s.secondCompressSdt.E, s.secondCompressSdt.compressIndex)
				// s.secondCompressSdt.SetE(s.E * float64(len(s.listOutPoint)) / float64(s.compressOutDataMax))
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
