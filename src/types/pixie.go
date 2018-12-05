package types

type VerifyResult struct {
	ID  string
	Num int
}

type VerifyResultList []VerifyResult

func (r VerifyResultList) Len() int {
	return len(r)
}
func (r VerifyResultList) Less(i, j int) bool {
	return r[i].Num > r[j].Num
}
func (r VerifyResultList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
