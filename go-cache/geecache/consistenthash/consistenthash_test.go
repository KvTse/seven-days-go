package consistenthash

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMap_Get(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// 添加节点名称为 6/4/2 的三个真实节点
	// then keys is 6/16/26/4/14/24/2/12/22 -> 2/4/6/12/14/16/22/24/26
	// 27循环回来了
	hash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		node := hash.Get(k)
		fmt.Printf("get node %s\n", node)
		if node != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	// then keys change to  2/4/6/8/12/14/16/18/22/24/26/28
	// 27 hit virtual node 28 the true node is 8
	hash.Add("8")
	for k, v := range testCases {
		node := hash.Get(k)
		fmt.Printf("get node %s\n", node)
		if node != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
