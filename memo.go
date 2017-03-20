package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const data = 10000
const precision = 64
const inputs = 6
const slack = 1000

var memory map[hash]float64

type vector []float64
type hash uint64
type planes []*vector

func (v *vector) lsh(p planes) hash {
	var h hash
	for i, r := range p {
		if v.dot(r) >= 0 {
			h ^= 1 << uint(i)
		}
	}

	return h
}

func (v *vector) dot(w *vector) float64 {
	if len(*v) != len(*w) {
		panic("too stupid to math")
	}

	var result float64

	for i, c := range *v {
		result += c * (*w)[i]
	}

	return result
}

func newPlanes(n int, dim int) planes {
	p := make(planes, n)

	for i := range p {
		v := make(vector, dim)
		for c := range v {
			v[c] = rand.Float64()*2 - 1
		}
		p[i] = &v
	}

	return p
}

func slowSum(input vector) float64 {
	result := 0.0
	for _, i := range input {
		time.Sleep(time.Nanosecond * slack)
		result += i
	}

	return result
}

func cache(fn func(vector) float64, p planes, input vector) float64 {
	h := input.lsh(p)
	result, ok := memory[h]

	if !ok {
		result = fn(input)
		memory[h] = result
	}

	return result
}

func main() {
	memory = make(map[hash]float64)
	seed := time.Now().Unix()
	rng := rand.New(rand.NewSource(seed))

	fmt.Print("Calc\t")
	sum := 0.0
	start := time.Now()
	for n := 0; n < data; n++ {
		v := make(vector, inputs)
		for i := range v {
			v[i] = rng.Float64()*2 - 1
		}
		sum += slowSum(v)
	}
	fmt.Printf("= %6f\tin %v\n", sum, time.Since(start))

	rng.Seed(seed)
	fmt.Print("Cache\t")
	cacheSum := 0.0
	p := newPlanes(precision, inputs)

	start = time.Now()
	for n := 0; n < data; n++ {
		v := make(vector, inputs)
		for i := range v {
			v[i] = rng.Float64()*2 - 1
		}
		cacheSum += cache(slowSum, p, v)
	}
	fmt.Printf("= %6f\tin %v\n", cacheSum, time.Since(start))

	hits := 100 * float64(data-len(memory)) / data
	diff := 100 * math.Abs((sum-cacheSum)/sum)

	fmt.Printf("\t%3.2f%% diff\n\t%3.2f%% hits\n", diff, hits)
}
