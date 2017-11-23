// Create a (not-very) pseudorandom cycle using modular exponentiation.
package powercycle

import (
	"fmt"
	"github.com/cznic/mathutil"
	"math/big"
	"math/rand"
)

// Cyclic permutation based on modular exponentation
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
	return newCycle(n, 2)
}

// Create a new Cycle on {0, ..., n-1} based on F_p such that
// r is a divisor of p - 1.
func newCycle(n uint64, r uint32) (cyc Cycle) {
	if r%2 == 1 {
		r *= 2
	}
	var cr uint64 = uint64(r)
	cyc.n.SetUint64(n)
	if n <= 4 {
		cyc.trivial = true
		return
	}
	p := n + 1

	// We are looking for the smallest prime above n of the form r * q + 1,
	// where q is some other prime.
	if (p-1)%cr != 0 {
		p += cr - ((p - 1) % cr)
	}
	for !mathutil.IsPrimeUint64((p-1)/cr) || !mathutil.IsPrimeUint64(p) {
		p += cr
	}
	cyc.p.SetUint64(p)

	// g is a primitive root modulo p if and only if
	//    g ^ (phi(p) / d)  !=  1
	// for every prime divisor d of phi(p) = p - 1.
	// As phi(p) = p - 1 = r * q, we need to factorize r.
	rFactors := mathutil.FactorInt(r)

	for {
		cyc.g.SetInt64(rand.Int63n(int64(p-2)) + 2)

		// check whether cyc.g is a generator of F_{cyc.p}
		var tmp big.Int
		tmp.SetUint64(cr)
		tmp.Exp(&cyc.g, &tmp, &cyc.p)
		if tmp.Uint64() == 1 {
			continue
		}

		ok := true
		for _, factor := range rFactors {
			tmp.SetUint64((p - 1) / uint64(factor.Prime))
			tmp.Exp(&cyc.g, &tmp, &cyc.p)
			if tmp.Uint64() == 1 {
				ok = false
				break
			}
		}

		if ok {
			break
		}
	}
	return
}

// Creates a permutation on n elements which consists of (at most) m cycles of
// approximately the same size. Also returns for each of these cycles
// an element in it.
func NewSplit(n uint64, m uint32) (per Cycles, xs []uint64) {
	var oldG, tmp big.Int
	var i uint32
	r := m
	if r == 1 {
		r = 2
	}
	per.cyc = newCycle(n, r)
	if per.cyc.trivial {
		return per, []uint64{0}
	}
	oldG.Set(&per.cyc.g)
	per.m.SetUint64(uint64(m))
	per.cyc.g.Exp(&oldG, &per.m, &per.cyc.p)
	xs = make([]uint64, 0, m)
	tmp.SetUint64(1)
	for i = 0; i < m; i++ {
		var x big.Int
		x.Set(&tmp)
		for {
			x.Mul(&x, &per.cyc.g)
			x.Mod(&x, &per.cyc.p)
			if x.Cmp(&per.cyc.n) <= 0 {
				xs = append(xs, x.Uint64()-1)
				break
			}
			if x.Cmp(&tmp) == 0 {
				break
			}
		}
		tmp.Mul(&tmp, &oldG)
		tmp.Mod(&tmp, &per.cyc.p)
	}
	return
}

// Permutation consisting of m approximately similar sized Cycles
type Cycles struct {
	cyc Cycle
	m   big.Int // number of cycles in the permutation
}

// Applies the permutation to the element x.
func (per *Cycles) Apply(x uint64) uint64 {
	return per.cyc.Apply(x)
}

func (per Cycles) String() string {
	if per.cyc.trivial {
		return fmt.Sprintf("<Cycles n=%d m=1 (trivial)>", per.cyc.n.Uint64())
	}
	return fmt.Sprintf("<Cycles n=%d m=%d p=%d g=%d>",
		per.cyc.n.Uint64(),
		per.m.Uint64(),
		per.cyc.p.Uint64(), per.cyc.g.Uint64())
}
