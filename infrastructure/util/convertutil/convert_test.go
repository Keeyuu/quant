package convertutil

import (
	"testing"
)

func TestName(t *testing.T) {
	contract1 := ContractInfo{
		Name:  "kk",
		Phone: "188374617348",
	}
	contract2 := ContractInfo{
		Name:  "ls",
		Phone: "18337423f48",
	}
	address1 := AddressInfo{
		Detail:    "广东",
		PostCode:  "510000",
		Emergency: contract1,
		Contract:  []ContractInfo{contract1, contract2},
	}
	address2 := AddressInfo{
		Detail:   "江苏",
		PostCode: "510042",
	}

	tf := TestField{
		Name:       "Joe",
		Age:        30,
		Address:    address1,
		ProductIds: []string{"13123218980321", "1478fjaf134fjlda"},
		SpuName: map[string]string{
			"cn-CN": "中国", "vi-VN": "越南",
		},
		ExtraNeedings: []AddressInfo{address1, address2},
	}
	m, err := GetMongoModelFieldValue(tf)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("map:%+v", m)
}
