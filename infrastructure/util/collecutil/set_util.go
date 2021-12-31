package collecutil

import mapset "github.com/deckarep/golang-set"

func Convert2StringSlice(set mapset.Set) (res []string) {
	res = make([]string, 0, set.Cardinality())
	interfaceSlice := set.ToSlice()
	for i := 0; i < len(interfaceSlice); i++ {
		res = append(res, interfaceSlice[i].(string))
	}
	return
}
