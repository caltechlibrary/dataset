package dataset

import (
	"math/rand"
)

// ShuffleStrings shuffles an array of strings in place.
// NOTE: You need to initialize your random number generator, e.g.
//    random := rand.New(rand.NewSource(time.Now().UnixNano())
// This maybe obsolete after go v1.10.x
func ShuffleStrings(a []string, random *rand.Rand) {
	for i := len(a) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
