package dolang

import (
	"math/rand"
	"time"
)

const INT_MAX = int(^uint(0) >> 1)

func init() {
	defineVMInsts()
	// for rand
	rand.Seed(time.Now().Unix())
}
