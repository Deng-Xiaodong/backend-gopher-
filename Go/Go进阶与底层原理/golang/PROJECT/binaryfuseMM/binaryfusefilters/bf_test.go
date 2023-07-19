package binaryfusefilters

import (
	"binaryfuseMM/util"
	"fmt"
	"log"
	"math"
	"testing"
)

var data = []KV{{"key_0", "value_0", 0}, {"key_0", "value_1", 1}, {"key_0", "value_10", 2},
	{"key_1", "value_0", 0}, {"key_1", "value_9", 1}, {"key_0", "value_20", 3},
	{"key_2", "value_0", 0}, {"key_2", "value_0", 1}, {"key_0", "value_30", 4},
	{"key_3", "value_0", 0}, {"key_3", "value_0", 1}, {"key_0", "value_100", 5},
	{"key_4", "value_0", 0}, {"key_4", "value_0", 1}, {"key_0", "value_200", 6},
	{"key_5", "value_0", 0}, {"key_5", "value_0", 1}, {"key_0", "value_300", 7},
	{"key_6", "value_0", 0}, {"key_6", "value_0", 1}, {"key_0", "value_1000", 8},
	{"key_7", "value_0", 0}, {"key_7", "value_0", 1}, {"key_0", "value_2000", 9},
	{"key_8", "value_0", 0}, {"key_8", "value_0", 1}, {"key_0", "value_3000", 10},
	{"key_9", "value_0", 0}, {"key_9", "value_0", 1}, {"key_0", "value_110", 11},
	{"key_10", "value_0", 0}, {"key_10", "value_0", 1}, {"key_0", "value_120", 12},
	{"key_11", "value_0", 0}, {"key_11", "value_0", 1}, {"key_0", "value_130", 13},
	{"key_12", "value_0", 0}, {"key_12", "value_0", 1}, {"key_0", "value_140", 14},
	{"key_13", "value_0", 0}, {"key_13", "value_0", 1}, {"key_0", "value_150", 15},
	{"key_14", "value_0", 0}, {"key_14", "value_0", 1}, {"key_0", "value_160", 16},
	{"key_15", "value_0", 0}, {"key_15", "value_0", 1}, {"key_0", "value_170", 17},
	{"key_16", "value_0", 0}, {"key_16", "value_0", 1}, {"key_15", "value_0", 2},
	{"key_17", "value_0", 0}, {"key_17", "value_0", 1}, {"key_16", "value_0", 2},
	{"key_18", "value_0", 0}, {"key_18", "value_0", 1}, {"key_17", "value_0", 2},
	{"key_19", "value_0", 0}, {"key_19", "value_0", 1}, {"key_18", "value_0", 2},
	{"key_20", "value_0", 0}, {"key_20", "value_0", 1}, {"key_19", "value_0", 2},
	{"key_21", "value_0", 0}, {"key_21", "value_0", 1}, {"key_20", "value_0", 2},
	{"key_22", "value_0", 0}, {"key_22", "value_0", 1}, {"key_21", "value_0", 2},
	{"key_23", "value_0", 0}, {"key_23", "value_0", 1}, {"key_22", "value_0", 2},
	{"key_24", "value_0", 0}, {"key_24", "value_0", 1}, {"key_23", "value_0", 2},
	{"key_25", "value_0", 0}, {"key_25", "value_0", 1}, {"key_24", "value_0", 2}}

func TestMappingStep(t *testing.T) {
	//不能存在重复的<key,count>
	//处理数据集的时候可能存在超过两个空白KV，两个<"",0>导致失败
	//data := make([]KV, 1e7)
	//for i := 0; i < len(data); i++ {
	//	data[i] = KV{key: fmt.Sprintf("key_%d", i), value: fmt.Sprintf("value_%d", i), count: 0}
	//}

	//for i, j := len(data)-len(data1), 0; i < 1e3; i, j = i+1, j+1 {
	//	data[i] = data1[j]
	//}

	//maxVolum关系到level，level关系到tree-base constraint PRF
	maxVolum := uint32(18)
	level := int(math.Ceil(math.Log(float64(maxVolum)) / math.Log(3.0))) //GGM Tree level

	f := SetupBinaryFuseMM(data, level, uint32(3))
	fmt.Printf("服务器存储消耗为%f\n", float64(f.Capacity)/float64(len(data)))
	s := NewServer(maxVolum, f.SegmentLength, f.SegmentCount, level)
	queryKey := "key_0"
	resp := s.Query(util.GetSha256([]byte(queryKey)))

	println(queryKey + "搜素结果如下：")
	k := k_list[queryKey]
	for i, bytes := range resp {
		decrypt := string(util.AesDecrypt(k, bytes))
		fmt.Printf("row_%d:%s\n", i, decrypt)
	}

}
func TestSingel(t *testing.T) {
	for i := 0; i < 5; i++ {
		bytes := util.GetSha256([]byte("key_1"))
		fmt.Printf("%d:%v\n", i, bytes)
		hash := Tri_GGM_Path_1(bytes, 0, 3)
		fmt.Printf("%d:%d\n", i, hash)
		h0, h1, h2 := Get3HashFromHash(hash, 214, 26)
		log.Printf("%d:{%d,%d,%d}\n", i, h0, h1, h2)
	}

}

func TestA(t *testing.T) {
	a := make([]int, 3, 6)
	b := make([]int, 3, 6)
	copy(b, a)
	b[0] = 99
	fmt.Printf("ptr_a:%p\n%v\nptr_b:%p\n%v\n", a, a, b, b)

}
