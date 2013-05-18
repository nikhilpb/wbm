package main

import (
	"fmt"
	"github.com/nikhilpb/wbm/block"
	"math"
	"os"
	"runtime"
	"strconv"
)

var N int = 10
var blk []block.Block
var xi_bar, xa_bar []float64
var xab_aux, xib_aux [][]float64
var rho float64 = 1.0
var step_count = 500

func main() {
	p_flag := false

	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) > 1 {
		n, _ := strconv.Atoi(os.Args[1])
		N = n
	} else {
		N = 10
	}

	if len(os.Args) > 2 {
		p_flag = true
	}

	// initialize w
	w := make([][]float64, N)
	for i := 0; i < N; i++ {
		w[i] = make([]float64, N)
		for j := 0; j < N; j++ {
			if i == j {
				w[i][j] = float64(i)
			} else {
				w[i][j] = float64(i) - math.Abs(float64(i-j))
			}
		}
	}

	blk = make([]block.Block, N)
	for i := 0; i < N; i++ {
		block.Init(&blk[i], w[i], N, i)
	}

	xa_bar = make([]float64, N)
	xi_bar = make([]float64, N)

	xab_aux = make([][]float64, N)
	xib_aux = make([][]float64, N)
	for i := 0; i < N; i++ {
		xab_aux[i] = make([]float64, N)
		xib_aux[i] = make([]float64, N)
	}

	if p_flag {
		fmt.Printf("parallel mode: ")
		fmt.Printf("maximum %d CPUs used\n", runtime.NumCPU())
		admm_parallel()
	} else {
		fmt.Println("serial mode")
		admm_serial()
	}

	obj := 0.0
	for i := 0; i < N; i++ {
		obj += xa_bar[i] + xi_bar[i]
	}
	target_obj := float64(N * (N - 1.0) / 2.0)
	err := math.Abs((obj-target_obj)/target_obj)
	fmt.Printf("Objective : %f, Target: %f, Error: %f\n", obj, target_obj, err)
}

func average() {
	for i := 0; i < N; i++ {
		xa_bar[i] = 0.0
		xi_bar[i] = 0.0
		for j := 0; j < N; j++ {
			xa_bar[i] += blk[j].Xa[i]
			xi_bar[i] += blk[j].Xi[i]
		}
		xi_bar[i] = xi_bar[i] / float64(N)
		xa_bar[i] = xa_bar[i] / float64(N)
	}
}

func project_all(low int, high int, tid int) {
	for i := low; i < high; i++ {
		for j := 0; j < N; j++ {
			xab_aux[tid][j] = xa_bar[j] - (1.0+blk[i].Ya[j])/rho
			xib_aux[tid][j] = xi_bar[j] - (1.0+blk[i].Yi[j])/rho
		}
		blk[i].Project(xab_aux[tid], xib_aux[tid])
	}
}

func dual_update(low int, high int) {
	for i := low; i < high; i++ {
		for j := 0; j < N; j++ {
			blk[i].Yi[j] += rho * (blk[i].Xi[j] - xi_bar[j])
			blk[i].Ya[j] += rho * (blk[i].Xa[j] - xa_bar[j])
		}
	}
}

func admm_parallel_help(i int, ch chan int) {
	for j := 0; j < N; j++ {
		blk[i].Yi[j] += rho * (blk[i].Xi[j] - xi_bar[j])
		blk[i].Ya[j] += rho * (blk[i].Xa[j] - xa_bar[j])
		xab_aux[i][j] = xa_bar[j] - (1.0+blk[i].Ya[j])/rho
		xib_aux[i][j] = xi_bar[j] - (1.0+blk[i].Yi[j])/rho
	}
	blk[i].Project(xab_aux[i], xib_aux[i])
	ch <- 1
}

func admm_parallel() {
	ch := make(chan int, N)
	for t := 0; t < step_count; t++ {
		for i := 0; i < N; i++ {
			go admm_parallel_help(i, ch)
		}
		for i := 0; i < N; i++ {
			<-ch
		}
		average()
	}
}

func admm_serial() {
	for t := 0; t < step_count; t++ {
		dual_update(0, N)
		project_all(0, N, 0)
		average()
	}
}
