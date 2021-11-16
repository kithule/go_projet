package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

//Result deines the result of a Computation
type Result struct {
	product  int
	row, col int
}

//Compute take a Computation as input and writes the Result on to a channel
func Compute(row, col int, rowSlice, colSlice []int, out chan<- Result) {
	n := len(rowSlice)
	product := 0
	for i := 0; i < n; i++ {
		product += rowSlice[i] * colSlice[i]
	}
	r := Result{row: row, col: col, product: product}
	out <- r
	wg.Done()
}
func transpose(a [][]int) [][]int {
	r, c := len(a), len(a[0])
	aTranspose := make([][]int, c)

	for i := 0; i < c; i++ {
		aTranspose[i] = make([]int, r)
		for j := 0; j < r; j++ {

			aTranspose[i][j] = a[j][i]
		}
	}
	return aTranspose
}

//MultiplyMatricesGoRoutine multiples 2 matrices using go routines & channels

func MultiplyMatricesGoRoutine(a, b [][]int) (product [][]int) {
	r1, c2 := len(a), len(b[0])

	product = make([][]int, r1)
	out := make(chan Result, r1*c2)

	bTranspose := transpose(b)
	goroutnum := 0

	if c2 != r1 {
		panic("b column size != a row size")
	}

	for i := 0; i < r1; i++ {
		product[i] = make([]int, c2)
		for j := 0; j < c2; j++ {
			goroutnum++
			wg.Add(1)
			go Compute(i, j, a[i], bTranspose[j], out)

		}
	}
	fmt.Printf(" %v Goroutines started\n", goroutnum)

	//collection
	for i := 0; i < r1; i++ {
		for j := 0; j < c2; j++ {
			result := <-out
			product[result.row][result.col] = result.product
		}
	}
	return product
}
func MultiplyMatrices(a, b [][]int) (product [][]int) {

	r1, c2, r2 := len(a), len(b[0]), len(b)
	product = make([][]int, r1)

	if c2 != r1 {
		panic("a column size != b row size")
	}
	for i := 0; i < r1; i++ {
		product[i] = make([]int, c2)
		for j := 0; j < c2; j++ {
			for k := 0; k < r2; k++ {
				product[i][j] += a[i][k] * b[k][j]
			}

		}
	}

	return product
}

func generateRandomMatrix(r, c int) [][]int {
	rg := rand.New(rand.NewSource(time.Now().UnixNano()))

	var matrix [][]int = make([][]int, r)

	for i := 0; i < r; i++ {
		matrix[i] = make([]int, c)
		for j := 0; j < c; j++ {
			matrix[i][j] = rg.Intn(100)
		}
	}

	return matrix
}
func main() {
	a, b := generateRandomMatrix(1000, 1000), generateRandomMatrix(1000, 1000)
	//c, d := generateRandomMatrix(1000, 1000), generateRandomMatrix(1000, 1000)

	showMatrixElements(a)
	showMatrixElements(b)
	showMatrixElements(MultiplyMatricesGoRoutine(a, b))

	now := time.Now()
	MultiplyMatricesGoRoutine(a, b)
	//MultiplyMatricesGoRoutine(c, d)

	fmt.Println(time.Now().Sub(now))

	now2 := time.Now()
	MultiplyMatrices(a, b)
	fmt.Println(time.Now().Sub(now2))

}
func showMatrixElements(m [][]int) {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[0]); j++ {
			fmt.Printf("%d\t", m[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}
