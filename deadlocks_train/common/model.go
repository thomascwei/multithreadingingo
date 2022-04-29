package common

import "sync"

type Train struct {
	Id int
	// 火車長度
	TrainLength int
	// 車頭位置
	Front int
}

// Intersection 火車與交岔點的互動
type Intersection struct {
	Id       int
	Mutex    sync.Mutex
	LockedBy int
}

// Crossing 交岔點
type Crossing struct {
	// 某條軌道起點到交岔點的距離,也就是交岔點的位置
	Position     int
	Intersection *Intersection
}
