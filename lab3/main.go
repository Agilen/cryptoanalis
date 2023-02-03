package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"time"

	"github.com/ALTree/bigfloat"
)

var zero = big.NewInt(0)
var one = big.NewInt(1)

type CrtEntry struct {
	C, N *big.Int
}

func main() {

	C, N, MITH := read()
	var crtentry []CrtEntry
	for i, _ := range C {
		crtentry = append(crtentry, CrtEntry{C[i], N[i]})
	}
	st := time.Now()
	c := CRT(C, N)
	fmt.Printf("%x\n\n", SolveCrtMany(crtentry))
	fmt.Printf("%x\n", c)
	x0 := new(big.Int).Div(c, big.NewInt(5))
	m := root(c, x0.Sqrt(x0.Sqrt(x0)))
	fmt.Println("\n\n")
	fmt.Printf("%x\n", m)
	fmt.Println(time.Since(st).Seconds())
	////////////////////////////
	st = time.Now()
	t, s := meetInTheMiddle(MITH[0], MITH[1], 65537, 20)
	ts := new(big.Int).Mul(big.NewInt(int64(t)), big.NewInt(int64(s)))
	fmt.Printf("%x\n", ts)
	fmt.Println(time.Since(st).Seconds())
}

func CRT(C []*big.Int, N []*big.Int) *big.Int {
	if len(C) != len(N) {
		return big.NewInt(0)
	}
	module := big.NewInt(1)
	for _, v := range N {
		module.Mul(module, v)
	}

	ans := big.NewInt(0)
	var modules []*big.Int

	for i := range N {
		modules = append(modules, new(big.Int).Div(module, N[i]))
	}

	for i, _ := range C {
		u := new(big.Int)
		v := new(big.Int)
		new(big.Int).GCD(u, v, N[i], modules[i])
		buf1 := big.NewInt(1).Mul(v, modules[i])
		buf1.Mul(buf1, C[i])
		ans.Add(ans, buf1)
	}

	return ans.Mod(ans, module)
}

func read() ([]*big.Int, []*big.Int, []*big.Int) {
	file, _ := ioutil.ReadFile("SE.txt")
	file2, _ := ioutil.ReadFile("MITH.txt")
	v := strings.Split(string(file), "\n")
	v2 := strings.Split(string(file2), "\n")

	var C []*big.Int
	var N []*big.Int
	var MITH []*big.Int

	for i := 0; i < len(v)-1; i += 2 {
		c := new(big.Int)
		c.SetString(v[i], 16)
		C = append(C, c)
		n := new(big.Int)
		n.SetString(v[i+1], 16)
		N = append(N, n)
	}
	mith1, _ := new(big.Int).SetString(v2[0], 16)
	mith2, _ := new(big.Int).SetString(v2[1], 16)
	MITH = append(MITH, mith1, mith2)
	return C, N, MITH
}

func SolveCrtMany(eqs []CrtEntry) *big.Int {
	if len(eqs) == 0 {
		panic("cannot have 0 entries to solve")
	}
	if len(eqs) == 1 {
		return new(big.Int).Mod(eqs[0].C, eqs[0].N)
	}
	eqs2 := make([]CrtEntry, len(eqs))
	copy(eqs2, eqs)
	return solveCrtManyIntern(eqs2)
}

func solveCrtManyIntern(eqs []CrtEntry) *big.Int {
	f := eqs[0]
	s := eqs[1]
	x := SolveCrt(f.C, f.N, s.C, s.N)
	if len(eqs) == 2 {
		return x
	}
	eqs[1] = CrtEntry{x, new(big.Int).Mul(f.N, s.N)}
	return solveCrtManyIntern(eqs[1:])
}

func SolveCrt(a, m, b, n *big.Int) *big.Int {
	gcd := new(big.Int)
	s := new(big.Int)
	t := new(big.Int)
	gcd.GCD(s, t, m, n)

	// let eqn = bsm, eqn2 = ant
	eqn := new(big.Int)
	eqn2 := new(big.Int)
	eqn.Mul(b, s)
	eqn.Mul(eqn, m)
	eqn2.Mul(a, n)
	eqn2.Mul(eqn2, t)

	// now, let eqn = bsm + ant, eqn2 = m * n
	eqn.Add(eqn, eqn2)
	eqn2.Mul(m, n)
	return eqn.Mod(eqn, eqn2)
}

func root(Aint *big.Int, x0int *big.Int) *big.Int {
	A := new(big.Float).SetInt(Aint)
	x0 := new(big.Float).SetInt(x0int)
	e := big.NewFloat(5)
	buf1, buf2, buf3 := new(big.Float), new(big.Float), new(big.Float)
	xk := x0
	xkp := new(big.Float).Copy(xk)
	I := 0
	for true {
		buf1 = bigfloat.Pow(xk, e)
		if buf3.Sub(buf1, A).Cmp(big.NewFloat(1)) == -1 || buf3.Sub(buf1, A).Cmp(big.NewFloat(1)) == 0 {
			buf3 = bigfloat.Pow(xkp, e)
			buf3.Sub(buf3, A)
			buf1.Sub(buf1, A)
			if buf1.Cmp(buf3) == 1 {
				return round(xkp)
			}
			return round(xk)
		}
		buf2 = bigfloat.Pow(xk, big.NewFloat(4))
		buf2.Mul(big.NewFloat(5), buf2)
		buf2.Quo(buf1, buf2)
		xkp.Copy(xk)
		xk.Sub(xk, buf2)
		I++
	}

	return round(xkp)
}

func round(f *big.Float) *big.Int {
	fInt := new(big.Int)
	f.Int(fInt)
	buf := new(big.Float)
	if f.Sub(f, buf.SetInt(fInt)).Cmp(big.NewFloat(0.5)) == -1 {
		return fInt
	}

	return fInt.Add(fInt, big.NewInt(1))
}

func meetInTheMiddle(c, n *big.Int, e, l int) (int, int) {
	r := (1 << (l >> 1)) - 1
	var array []*big.Int
	for i := 0; i < r; i++ {
		array = append(array, new(big.Int).Exp(big.NewInt(int64(i+1)), big.NewInt(int64(e)), n))
	}
	for i := 0; i < r; i++ {
		v := new(big.Int).ModInverse(array[i], n)
		v.Mul(v, c)
		v.Mod(v, n)

		for j := 0; j < r; j++ {
			if v.Cmp(array[j]) == 0 {
				return i + 1, j + 1
			}
		}
	}

	return 0, 0
}

func ffx(pow float64, a float64, i float64) (float64, float64, float64) {
	a *= pow
	pow -= 1
	return pow, a, i * (i + 1)
}
