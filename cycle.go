// Create a (not-very) pseudorandom cycle using modular exponentiation.

package powercycle

import (
	"fmt"
	"github.com/cznic/mathutil"
	"math/big"
	"math/rand"
)

type Cycle struct {
	trivial bool    // used for n <= 4
	n       big.Int // this is a cycle on {0, ..., n-1}
	p       big.Int // p is a prime with n+1 <= p and with (p-1)/2 prime
	g       big.Int // primitive root modulo p
}

func (cyc Cycle) String() string {
	if cyc.trivial {
		return fmt.Sprintf("<Cycle n=%d (trivial)>", cyc.n.Uint64())
	}
	return fmt.Sprintf("<Cycle n=%d p=%d g=%d>",
		cyc.n.Uint64(), cyc.p.Uint64(), cyc.g.Uint64())
}

// Applies the cycle to the element x.
func (cyc *Cycle) Apply(x uint64) uint64 {
	if cyc.trivial {
		return (x + 1) % cyc.n.Uint64()
	}
	var y big.Int
	y.SetUint64(x + 1)
	for {
		y.Mul(&y, &cyc.g)
		y.Mod(&y, &cyc.p)
		if y.Cmp(&cyc.n) <= 0 {
			return y.Uint64() - 1
		}
	}
}

// Create a new (not-so-)pseudorandom Cycle on {0, ..., n-1}.
// See Cycle.Apply to use it.
func New(n uint64) (cyc Cycle) {
	cyc.n.SetUint64(n)
	if n <= 4 {
		cyc.trivial = true
		return
	}
	p := n + 1
	if p%2 == 0 {
		p += 1
	}
	for !mathutil.IsPrimeUint64((p-1)/2) || !mathutil.IsPrimeUint64(p) {
		p += 2
	}
	cyc.p.SetUint64(p)
	for {
		cyc.g.SetInt64(rand.Int63n(int64(p-2)) + 2)

		// check whether cyc.g is a generator of F_{cyc.p}
		var tmp big.Int
		tmp.Mul(&cyc.g, &cyc.g)
		tmp.Mod(&tmp, &cyc.p)
		if tmp.Uint64() == 1 {
			continue
		}

		tmp.SetUint64((p - 1) / 2)
		tmp.Exp(&cyc.g, &tmp, &cyc.p)
		if tmp.Uint64() == 1 {
			continue
		}
		break
	}
	return
}
