package arbitrator

import (
	. "github.com/cutajarj/multithreadingingo/deadlocks_train/common"
	"sync"
	"time"
)

var (
	controller = sync.Mutex{}
	//cond用於協調想要訪問共享資源的那些goroutines。
	//當共享資源的狀態發生改變時，它可以被用來通知被互斥鎖阻塞的goroutine。
	cond = sync.NewCond(&controller)
)

//檢查必須同時鎖定的交岔點是否可同時獲得
func allFree(intersectionsToLock []*Intersection) bool {
	for _, it := range intersectionsToLock {
		if it.LockedBy >= 0 {
			return false
		}
	}
	return true
}

//依車身長度對同時會通過的交岔點上鎖
func lockIntersectionsInDistance(id, reserveStart int, reserveEnd int, crossings []*Crossing) {
	var intersectionsToLock []*Intersection
	//整理出待上鎖slice
	for _, crossing := range crossings {
		if reserveEnd >= crossing.Position &&
			reserveStart <= crossing.Position &&
			crossing.Intersection.LockedBy != id {
			intersectionsToLock = append(intersectionsToLock, crossing.Intersection)
		}
	}
	//controller為全局變量
	controller.Lock()
	//如果無法一次性全部上鎖就會等待cond通知再回來
	for !allFree(intersectionsToLock) {
		//暫時先跳出 等接到cond通知會繼續
		cond.Wait()
	}
	//以改id取代真正的上鎖
	for _, it := range intersectionsToLock {
		it.LockedBy = id
		time.Sleep(10 * time.Millisecond)
	}
	controller.Unlock()
}

func MoveTrain(train *Train, distance int, crossings []*Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				lockIntersectionsInDistance(train.Id, crossing.Position, crossing.Position+train.TrainLength, crossings)
			}
			back := train.Front - train.TrainLength
			if back == crossing.Position {
				//通過交岔路口後將id改為-1代表未被鎖定,並通知所有wait的goroutine
				controller.Lock()
				crossing.Intersection.LockedBy = -1
				cond.Broadcast()
				controller.Unlock()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
