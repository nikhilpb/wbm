package block

import (
	"math"
)

var tolerance float64 = 0.000001

type Block struct {
	W        []float64 //.Weights
	Xa       []float64 // primal variable, ad
	Xi       []float64 // primal variable, impression
	Lmbd     []float64 // dual variables
	Lmbd_sum float64   // sum of dual variables
	Ya       []float64 // dual variable advertiser
	Yi       []float64 // dual variable impression
	Ind      int       // the Index of the corresponding ad
	N				 int
}

func (b Block) Project(Xab []float64, Xib []float64) {
	res := make([]float64, b.N)
	sat := true
	// check.Which constraints are satisfied
	for i := 0; i < b.N; i++ {
		res[i] = b.W[i] - Xab[b.Ind] - Xib[i]
		if res[i] < tolerance {
			b.Lmbd[i] = 0.0
		} else {
			sat = false
		}
	}
	// if all constraints are satisfied set x = xb and quit
	if sat {
		for i := 0; i < b.N; i++ {
			b.Xi[i] = Xib[i]
			b.Xa[i] = Xab[i]
		}
		return
	}

	// main loop of the projection function
	var ln, lo float64
	var ic int
	for true {
		// co-ordinate search
		sat = true
		for i := 0; i < b.N; i++ {
			if b.Lmbd[i]+b.Lmbd_sum-res[i] < -tolerance {
				ic = i
				sat = false
				break
			} else if (b.Lmbd[i] > tolerance) && (math.Abs(b.Lmbd[i]+b.Lmbd_sum-res[i]) > tolerance) {
				ic = i
				sat = false
				break
			}
		}

		// if no direction eXists, update x and break
		if sat {
			for i := 0; i < b.N; i++ {
				b.Xi[i] = Xib[i] + b.Lmbd[i]
				b.Xa[i] = Xab[i]
			}
			b.Xa[b.Ind] = Xab[b.Ind] + b.Lmbd_sum
			break
		}

		// update lambda and sum of lambdas
		lo = b.Lmbd[ic]
		ln = (res[ic] - b.Lmbd_sum + b.Lmbd[ic]) / 2.0
		if ln >= 0.0 {
			b.Lmbd[ic] = ln
		} else {
			b.Lmbd[ic] = 0.0
		}
		b.Lmbd_sum += b.Lmbd[ic] - lo
	}
	return
}

func Init(b *Block, W []float64, size int, Ind int){
	b.N = size
  b.W        = make([]float64, b.N)
  b.Xa       = make([]float64, b.N)
	b.Xi       = make([]float64, b.N)
  b.Lmbd     = make([]float64, b.N)
  b.Lmbd_sum = 0.0
  b.Ya       = make([]float64, b.N)
  b.Yi       = make([]float64, b.N)
  b.Ind = Ind
  for i := 0; i < b.N; i++{
    b.W[i]    = W[i]
    b.Xa[i]   = 0.0
    b.Xi[i]   = 0.0
    b.Lmbd[i] = 0.0
    b.Ya[i]   = 0.0
    b.Yi[i]   = 0.0
  }
}
