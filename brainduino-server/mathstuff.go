package main

import "math/cmplx"

func abs(cin []complex128) []float64 {
	fout := make([]float64, len(cin))
	for idx, v := range cin {
		fout[idx] = cmplx.Abs(v)
	}
	return fout
}
