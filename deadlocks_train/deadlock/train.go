package deadlock

import (
	. "github.com/cutajarj/multithreadingingo/deadlocks_train/common"
	"time"
)

// MoveTrain distance代表火車需要前進的距離,crossings代表將會經過過的交岔路口列表
func MoveTrain(train *Train, distance int, crossings []*Crossing) {
	//在火車到達目的地前持續移動火車
	for train.Front < distance {
		//每次移動1
		train.Front += 1
		//檢查我們是否在一個路口上
		for _, crossing := range crossings {
			//剛好在路口
			if train.Front == crossing.Position {
				//取得鎖
				crossing.Intersection.Mutex.Lock()
				//表明該路口被哪一台火車占用
				crossing.Intersection.LockedBy = train.Id
			}
			//車尾
			back := train.Front - train.TrainLength
			//判斷是否通過交岔點
			if back == crossing.Position {
				crossing.Intersection.LockedBy = -1
				crossing.Intersection.Mutex.Unlock()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
