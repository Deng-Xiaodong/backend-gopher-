package binaryfusefilters

import (
	"binaryfuseMM/util"
	"log"
	"math"
	"math/bits"
	"math/rand"
	"strconv"
)

var (
	k_list   map[string][]byte
	enc_list [][]byte
	EMM      [][]byte
)

const (
	MaxIterations = 1000
	K_e           = 0x234
)

type BinaryFuseMM struct {
	Seed          uint64
	SegmentLength uint32
	SegmentCount  uint32
	Capacity      uint32
}
type KV struct {
	key   string
	value string
	count int
}
type StackItem struct {
	k     int
	index uint32
}

// 将关键字处理为uint64
func keys2uint64(kv_list []KV, level int) []uint64 {

	size := len(kv_list)
	hashes := make([]uint64, size)
	for i := 0; i < size; i++ {
		hashes[i] = Tri_GGM_Path_1(util.GetSha256([]byte(kv_list[i].key)), kv_list[i].count, level)
	}
	return hashes
}

// tk只与关键字有关，且作为GGM的根部
func Tri_GGM_Path_1(tk []byte, count, level int) uint64 {

	kv_bytes := util.Tri_GGM_Path(tk, level, util.Tts(count, 3, level))
	return uint64(util.BytesToInt64(kv_bytes))

}

func calculateSegmentLength(arity uint32, size uint32) uint32 {

	if arity == 3 {
		return uint32(1) << int(math.Floor(math.Log(float64(size))/math.Log(3.33)+2.25))
	} else if arity == 4 {
		return uint32(1) << int(math.Floor(math.Log(float64(size))/math.Log(2.91)-0.5))
	} else {
		return 65536
	}
}

func calculateSizeFactor(arity uint32, size uint32) float64 {
	if arity == 3 {
		return math.Max(1.125, 0.875+0.25*math.Log(1000000)/math.Log(float64(size)))
	} else if arity == 4 {
		return math.Max(1.075, 0.77+0.305*math.Log(600000)/math.Log(float64(size)))
	} else {
		return 2.0
	}
}

// returns random number, modifies the seed
func splitmix64(seed *uint64) uint64 {
	*seed = *seed + 0x9E3779B97F4A7C15
	z := *seed
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

// 用于对数据集的预处理（根据键的哈希值进行哈希排序；原论文说是根据它们所在的segment从小到打排序）
func mixsplit(key, seed uint64) uint64 {
	h := key + seed
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	h *= 0xc4ceb9fe1a85ec53
	h ^= h >> 33
	return h
}

func InitBinaryFuseMM(size uint32, arity uint32) *BinaryFuseMM {
	f := &BinaryFuseMM{}
	//arity := uint32(4)
	f.SegmentLength = calculateSegmentLength(arity, size)
	if f.SegmentLength > 262144 {
		f.SegmentLength = 262144
	}

	sizeFactor := calculateSizeFactor(arity, size)
	capacity := uint32(0)
	if size > 1 {
		capacity = uint32(math.Round(float64(size) * sizeFactor))
	}
	f.SegmentCount = uint32(math.Round(float64(capacity / f.SegmentLength)))
	capacity = f.SegmentLength * (f.SegmentCount + arity - 1)
	f.Capacity = capacity
	return f
}

// 确定性函数

func Get3HashFromHash(hash uint64, segmentLength, segmentCount uint32) (uint32, uint32, uint32) {
	//h0、 h1、h2 occupy three distinct and consecutive segments.

	//ho的范围[0,segmentC*segmentL]
	hi, _ := bits.Mul64(hash, uint64(segmentLength*segmentCount))
	h0 := uint32(hi)
	h1 := h0 + segmentLength
	h2 := h1 + segmentLength
	//h1 h2 会在当前segment置换
	h1 ^= uint32(hash>>18) & (segmentLength - 1)
	h2 ^= uint32(hash) & (segmentLength - 1)
	return h0, h1, h2

}
func Get4HashFromHash(hash uint64, segmentLength, segmentCount uint32) (uint32, uint32, uint32, uint32) {
	//h0、 h1、h2 occupy three distinct and consecutive segments.

	hi, _ := bits.Mul64(hash, uint64(segmentLength*segmentCount))
	h0 := uint32(hi)
	h1 := h0 + segmentLength
	h2 := h1 + segmentLength
	h3 := h2 + segmentLength
	h1 ^= uint32(hash>>21) & (segmentLength - 1)
	h2 ^= uint32(hash>>42) & (segmentLength - 1)
	h3 ^= uint32(hash) & (segmentLength - 1)
	return h0, h1, h2, h3

}

func SetupBinaryFuseMM(kv_list []KV, level int, arity uint32) *BinaryFuseMM {

	//参数初始化
	size := uint32(len(kv_list))
	//initKeyHash(int(size))
	f := InitBinaryFuseMM(size, arity)

	rngcounter := uint64(1)
	f.Seed = splitmix64(&rngcounter)
	capacity := f.Capacity

	//生成密钥和密文（指纹）
	enc_list = make([][]byte, size)
	k_list = make(map[string][]byte)

	for i, kv := range kv_list {

		var K []byte
		if vk, ok := k_list[kv.key]; ok {
			K = vk
		} else {
			K = util.GetSha256([]byte(strconv.Itoa(K_e) + kv.key))
			k_list[kv.key] = K
		}
		enc_list[i] = util.AesEncrypt(K, []byte(kv.value))
	}
	crypedLength := len(enc_list[0])
	//容器初始化
	reverseOrder := make([]int, size)
	alone := make([]uint32, capacity)
	t2count := make([]uint8, capacity)
	t2 := make([]int, capacity)
	var stack []*StackItem
	EMM = make([][]byte, capacity)

	iterations := 0
	for true {
		iterations += 1
		if iterations > MaxIterations {
			// The probability of this happening is lower than the
			// the cosmic-ray probability (i.e., a cosmic ray corrupts your system),
			// but if it happens, we just fill the fingerprint with ones which
			// will flag all possible keys as 'possible', ensuring a correct result.
			log.Fatalln("超过最大尝试次数")
		}

		//最小大于切片树的对数值
		blockBits := 1
		for (1 << blockBits) < f.SegmentCount {
			blockBits += 1
		}
		startPos := make([]uint, 1<<blockBits)
		for i, _ := range startPos {
			// important: we do not want i * size to overflow!!!
			startPos[i] = uint((uint64(i) * uint64(size)) >> blockBits)
		}

		keyshash := keys2uint64(kv_list, level)
		//根据key的第一个哈希值（h0）所在的segment位置排序
		for i, key := range keyshash {
			hash := mixsplit(key, f.Seed)
			segment_index := hash >> (64 - blockBits)
			for startPos[segment_index] >= uint(size) || reverseOrder[startPos[segment_index]] != 0 {
				segment_index++
				segment_index &= (1 << blockBits) - 1
			}
			reverseOrder[startPos[segment_index]] = i + 1
			startPos[segment_index] += 1
		}
		//println("初始化生成的key_1的index如下：\n")

		for i := 0; i < len(reverseOrder); i++ {

			k := reverseOrder[i] - 1
			//if k == -1 {
			//	continue
			//}
			//hash := keyshash[k]
			hash := Tri_GGM_Path_1(util.GetSha256([]byte(kv_list[k].key)), kv_list[k].count, level)
			index1, index2, index3 := Get3HashFromHash(hash, f.SegmentLength, f.SegmentCount)
			//index1, index2, index3, index4 := Get4HashFromHash(hash, f.SegmentLength, f.SegmentCount)

			t2count[index1]++
			t2[index1] ^= k
			t2count[index2]++
			t2[index2] ^= k
			t2count[index3]++
			t2[index3] ^= k
			//t2count[index4]++
			//t2[index4] ^= k
		}

		// End of key addition

		Qsize := 0
		// Add sets with one key to the queue.
		for i := uint32(0); i < capacity; i++ {
			alone[Qsize] = i //alone记录了哪些index为单一
			if t2count[i] == 1 {
				Qsize++
			}
		}

		stacksize := uint32(0)

		for Qsize > 0 {
			Qsize--
			index := alone[Qsize]
			if t2count[index] == 1 {
				k := t2[index]
				//hash := keyshash[k]
				hash := Tri_GGM_Path_1(util.GetSha256([]byte(kv_list[k].key)), kv_list[k].count, level)
				//删除该关键字的三个位置，并将(key,index)添加到stack
				stack = append(stack, &StackItem{k: k, index: index})
				stacksize++

				index1, index2, index3 := Get3HashFromHash(hash, f.SegmentLength, f.SegmentCount)
				//index1, index2, index3, index4 := Get4HashFromHash(hash, f.SegmentLength, f.SegmentCount)

				t2[index1] ^= k
				t2count[index1]--
				if t2count[index1] == 1 {
					alone[Qsize] = index1
					Qsize++
				}

				t2[index2] ^= k
				t2count[index2]--
				if t2count[index2] == 1 {
					alone[Qsize] = index2
					Qsize++
				}

				t2[index3] ^= k
				t2count[index3]--
				if t2count[index3] == 1 {
					alone[Qsize] = index3
					Qsize++
				}

				//t2[index4] ^= k
				//t2count[index4]--
				//if t2count[index4] == 1 {
				//	alone[Qsize] = index4
				//	Qsize++
				//}
			}
		}

		if stacksize == size {
			break
		}

		for i := uint32(0); i < size; i++ {
			reverseOrder[i] = 0
		}
		for i := uint32(0); i < capacity; i++ {
			t2count[i] = 0
			t2[i] = 0
		}
		rngcounter++
		f.Seed = splitmix64(&rngcounter)
	}

	for i := (int)(size - 1); i >= 0; i-- {

		// the hash of the key we insert next
		itme := stack[i]

		hash := Tri_GGM_Path_1(util.GetSha256([]byte(kv_list[itme.k].key)), kv_list[itme.k].count, level)
		index := itme.index

		xor := enc_list[itme.k]
		index1, index2, index3 := Get3HashFromHash(hash, f.SegmentLength, f.SegmentCount)
		//index1, index2, index3, index4 := Get4HashFromHash(hash, f.SegmentLength, f.SegmentCount)
		switch index {
		case index1:
			if EMM[index2] == nil {
				EMM[index2] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
				//EMM[index2] = make([]byte, crypedLength)
			}
			if EMM[index3] == nil {
				EMM[index3] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
				//EMM[index3] = make([]byte, crypedLength)
			}
			//if EMM[index4] == nil {
			//	//EMM[index3] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
			//	EMM[index4] = make([]byte, crypedLength)
			//}
			xor = util.Xor(xor, util.Xor(EMM[index2], EMM[index3]))
			//xor = util.Xor(xor, util.Xor(util.Xor(EMM[index2], EMM[index3]), EMM[index4]))

		case index2:
			if EMM[index1] == nil {
				EMM[index1] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
				//EMM[index1] = make([]byte, crypedLength)
			}
			if EMM[index3] == nil {
				EMM[index3] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
				//EMM[index3] = make([]byte, crypedLength)
			}
			//if EMM[index4] == nil {
			//	//EMM[index3] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
			//	EMM[index4] = make([]byte, crypedLength)
			//}
			xor = util.Xor(xor, util.Xor(EMM[index1], EMM[index3]))
			//xor = util.Xor(xor, util.Xor(util.Xor(EMM[index1], EMM[index3]), EMM[index4]))
		case index3:
			if EMM[index1] == nil {
				EMM[index1] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
				//EMM[index1] = make([]byte, crypedLength)
			}
			if EMM[index2] == nil {
				EMM[index2] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
				//EMM[index2] = make([]byte, crypedLength)
			}
			//if EMM[index4] == nil {
			//	//EMM[index3] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
			//	EMM[index4] = make([]byte, crypedLength)
			//}
			xor = util.Xor(xor, util.Xor(EMM[index1], EMM[index2]))
			//xor = util.Xor(xor, util.Xor(util.Xor(EMM[index1], EMM[index2]), EMM[index4]))
			//case index4:
			//	if EMM[index1] == nil {
			//		//EMM[index1] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
			//		EMM[index1] = make([]byte, crypedLength)
			//	}
			//	if EMM[index2] == nil {
			//		//EMM[index2] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
			//		EMM[index2] = make([]byte, crypedLength)
			//	}
			//	if EMM[index3] == nil {
			//		//EMM[index3] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
			//		EMM[index3] = make([]byte, crypedLength)
			//	}
			//	//xor = util.Xor(xor, util.Xor(EMM[index1], EMM[index2]))
			//	xor = util.Xor(xor, util.Xor(util.Xor(EMM[index1], EMM[index2]), EMM[index3]))
		}
		EMM[index] = xor
	}

	for i, bytes := range EMM {
		if bytes == nil {
			EMM[i] = util.GetSha256(util.Int64ToBytes(rand.Int63()))[:crypedLength]
			//EMM[i] = make([]byte, crypedLength)
		}
	}
	return f
}
