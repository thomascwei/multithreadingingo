package hierarchy

import (
	. "github.com/cutajarj/multithreadingingo/deadlocks_train/common"
	"sort"
	"time"
)

//依車身長度對同時會通過的交岔點上鎖
//reserveStart, reserveEnd為車頭車尾位置
func lockIntersectionsInDistance(id, reserveStart, reserveEnd int, crossings []*Crossing) {
	//待上鎖清單
	var intersectionsToLock []*Intersection
	for _, crossing := range crossings {
		//若交岔點位在車頭與車尾之間且未被自己上鎖則加入待上鎖清單
		if reserveEnd >= crossing.Position && reserveStart <= crossing.Position && crossing.Intersection.LockedBy != id {
			intersectionsToLock = append(intersectionsToLock, crossing.Intersection)
		}
	}
	//依id排序, 不會產生死鎖的關鍵
	sort.Slice(intersectionsToLock, func(i, j int) bool {
		return intersectionsToLock[i].Id < intersectionsToLock[j].Id
	})
	//依序上鎖
	for _, it := range intersectionsToLock {
		it.Mutex.Lock()
		it.LockedBy = id
		time.Sleep(10 * time.Millisecond)
	}
}

func MoveTrain(train *Train, distance int, crossings []*Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				//把會同時佔用的交岔口一次性全部上鎖
				lockIntersectionsInDistance(train.Id, crossing.Position, crossing.Position+train.TrainLength, crossings)
			}
			back := train.Front - train.TrainLength
			if back == crossing.Position {
				crossing.Intersection.LockedBy = -1
				crossing.Intersection.Mutex.Unlock()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
