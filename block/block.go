package block

import (
	"sort"
)

var tolerance float64 = 0.000001

type Block struct {
	W        []float64 //.Weights
	Xa       []float64 // primal variable, ad
	Xi       []float64 // primal variable, impression
	Ya       []float64 // dual variable advertiser
	Yi       []float64 // dual variable impression
	Ind      int       // the Index of the corresponding ad
	N				 int
}

type Pair struct {
	Ind 	int
	Value float64
}

type Pairs []Pair

func (s Pairs) Len() int      { return len(s) }
func (s Pairs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Pairs) Less(i, j int) bool { return s[i].Value < s[j].Value }

func (b Block) Project(Xab []float64, Xib []float64){
	ucs := make(Pairs, b.N)
	var res float64
	count := 0
	for i := 0; i < b.N; i++ {
    res = b.W[i] - Xab[b.Ind] - Xib[i]
    if (res > 0){
      ucs[count].Ind = i
      ucs[count].Value = res
      count++
    }
  }
 	sum := 0.0
  aic := 0
  if count > 0{
    sort.Sort(ucs[0:count])
    for i := count - 1; i > -1; i--{
      if (ucs[i].Value < (sum /(float64(aic) + 1))){
        break
      }
      sum += ucs[i].Value
      aic++
    }
  }
  for i := 0; i < b.N; i++ {
    b.Xa[i] = Xab[i]
    b.Xi[i] = Xib[i]
  }
  for i := count - 1; i > count - 1 -aic; i--{
    b.Xi[ucs[i].Ind] += ucs[i].Value - sum/(float64(aic) + 1)
  }
  b.Xa[b.Ind] += (1.0 / (float64(aic) + 1.0)) * sum
}

func Init(b *Block, W []float64, size int, Ind int){
	b.N = size
  b.W        = make([]float64, b.N)
  b.Xa       = make([]float64, b.N)
	b.Xi       = make([]float64, b.N)
  b.Ya       = make([]float64, b.N)
  b.Yi       = make([]float64, b.N)
  b.Ind = Ind
  for i := 0; i < b.N; i++{
    b.W[i]    = W[i]
    b.Xa[i]   = 0.0
    b.Xi[i]   = 0.0
    b.Ya[i]   = 0.0
    b.Yi[i]   = 0.0
  }
}
