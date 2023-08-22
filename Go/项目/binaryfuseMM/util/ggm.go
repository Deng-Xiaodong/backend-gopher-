package util

import (
	"crypto/sha256"
)

var m = make(map[int64][]byte)

func Tri_GGM_Path(root []byte, level int, path []int) []byte {
	current_node := root
	for i := 0; i < level; i++ {
		temp := path[i]
		switch temp {
		case 0:
			current_node = current_node[1:10]
			key := hash64(BytesToInt64(current_node), int64(temp))
			if _, ok := m[key]; ok {
				current_node = m[key]
			} else {
				bytes := sha256.Sum256(current_node)
				current_node = append([]byte(nil), bytes[:]...)
				m[key] = current_node
			}
		case 1:
			current_node = current_node[11:20]
			key := hash64(BytesToInt64(current_node), int64(temp))
			if _, ok := m[key]; ok {
				current_node = m[key]
			} else {
				bytes := sha256.Sum256(current_node)
				current_node = append([]byte(nil), bytes[:]...)
				m[key] = current_node
			}
		default:
			current_node = current_node[21:30]
			key := hash64(BytesToInt64(current_node), int64(temp))
			if _, ok := m[key]; ok {
				current_node = m[key]
			} else {
				bytes := sha256.Sum256(current_node)
				current_node = append([]byte(nil), bytes[:]...)
				m[key] = current_node

			}
		}
	}
	return current_node
}
func Tts(inNum int, index int, level int) []int {
	result := make([]int, level)
	for i := 0; i < level; i++ {
		result[i] = inNum % index
		inNum /= index
	}
	return result
}
