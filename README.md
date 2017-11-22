powercycle
==========

Go package to generate a (not-so\*) pseudo-random cycle.

Example
-------

```go
package main

import (
	"fmt"
	"github.com/bwesterb/powercycle"
)

func main() {
	var x uint64
	cycle := powercycle.New(10)
	for i := 0; i < 10; i++ {
		fmt.Println(x)
		x = cycle.Apply(x)
	}
}
```

might output

```
0
6
4
1
2
9
3
5
8
7
```

Not-so pseudorandom
-------------------

For efficiency, the cycles generated are of a very particular form.  So do
not use this package if you want to have a real pseudo-random cycle.

How it works
------------
To generate a cycle of size n, we find a prime p > n + 1 such that
furthermore (p - 1)/2 is also a prime.  This makes it easy to find a
generator g modulo p.  The action of the cycle is usually given by

    x ---> (((x + 1) * g) % p) - 1

As often p > n + 1, this might yield a number bigger than n.
In that case the action is repeated.
