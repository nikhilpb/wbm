package main

import (
	"fmt"
	"github.com/nikhilpb/wbm/block"
)

var N int = 10
var blk []block.Block
var xi_bar, xa_bar []float64
var xab_aux, xib_aux [][]float64
var n_threads = N
var rho float64 = 1.0

func main() {
	// initialize w
	w := make([][]float64, N)
	for i := 0; i < N; i++ {
		w[i] = make([]float64, N)
		for j := 0; j < N; j++ {
			if i == j {
				w[i][j] = float64(i)
			}
		}
	}

	blk = make([]block.Block, N)
	for i := 0; i < N; i++ {
		block.Init(&blk[i], w[i], N, i)
	}

	xa_bar = make([]float64, N)
	xi_bar = make([]float64, N)

	xab_aux = make([][]float64, n_threads)
	xib_aux = make([][]float64, n_threads)
	for i := 0; i < N; i++{
		xab_aux[i] = make([]float64, N)
		xib_aux[i] = make([]float64, N)
	}

	admm_serial()
	
	for i := 0; i < N; i++{
		fmt.Printf("%f %f\n", xa_bar[i], xi_bar[i])
	}
	fmt.Println("done")
}

func average(){
	for i := 0; i < N; i++{
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

func project_all(low int, high int, tid int){
  for i := low; i < high; i++ {
    for j := 0; j < N; j++ {
      xab_aux[tid][j] = xa_bar[j] - (1.0 + blk[i].Ya[j])/rho
      xib_aux[tid][j] = xi_bar[j] - (1.0 + blk[i].Yi[j])/rho
    }
    blk[i].Project(xab_aux[tid], xib_aux[tid])
  } 
}

func dual_update(low int, high int){
  for i := low; i < high; i++{ 
    for j := 0; j < N; j++ {
      blk[i].Yi[j] += rho * (blk[i].Xi[j] - xi_bar[j])
      blk[i].Ya[j] += rho * (blk[i].Xa[j] - xa_bar[j])
    }
  }
}

func admm_serial(){
  for t := 0; t < 1000; t++{
    dual_update(0, N);
    project_all(0, N, 0); 
    average();
  }
}

