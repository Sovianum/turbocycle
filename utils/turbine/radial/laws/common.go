package laws

func getRRel(hRel, lRel float64) float64 {
	return hRel*(1+lRel) + (1-hRel)*(1-lRel)
}
