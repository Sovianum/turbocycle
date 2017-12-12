package laws

func getRRel(hRel, lRel float64) float64 {
	return hRel*(1+lRel) + (1-hRel)*(1-lRel)
}

func getHRel(rRel, lRel float64) float64 {
	return (rRel + lRel - 1) / (2 * lRel)
}
