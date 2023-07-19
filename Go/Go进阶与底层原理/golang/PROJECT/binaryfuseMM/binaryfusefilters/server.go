package binaryfusefilters

import (
	"binaryfuseMM/util"
)

type server struct {
	maxVolume     uint32
	level         int
	segmentLength uint32
	segmentCount  uint32
	EMM           [][]byte
}

func NewServer(maxV, segmentL, segmentC uint32, level int) *server {
	return &server{
		maxVolume:     maxV,
		segmentLength: segmentL,
		segmentCount:  segmentC,
		level:         level,
		EMM:           EMM,
	}
}

func (s *server) Query(tk []byte) [][]byte {
	resp := make([][]byte, s.maxVolume)
	//println("服务器根据token生成的index如下：\n")
	for i := uint32(0); i < s.maxVolume; i++ {
		hash := Tri_GGM_Path_1(tk, int(i), s.level)
		h0, h1, h2 := Get3HashFromHash(hash, s.segmentLength, s.segmentCount)
		//log.Printf("%d:{%d,%d,%d}\n", i, h0, h1, h2)
		resp[i] = util.Xor(s.EMM[h0], util.Xor(s.EMM[h1], s.EMM[h2]))
	}
	return resp
}
func (s *server) Query4(tk []byte) [][]byte {
	resp := make([][]byte, s.maxVolume)
	//println("服务器根据token生成的index如下：\n")
	for i := uint32(0); i < s.maxVolume; i++ {
		hash := Tri_GGM_Path_1(tk, int(i), s.level)
		h0, h1, h2, h3 := Get4HashFromHash(hash, s.segmentLength, s.segmentCount)
		//log.Printf("%d:{%d,%d,%d}\n", i, h0, h1, h2)
		resp[i] = util.Xor(s.EMM[h0], util.Xor(util.Xor(s.EMM[h1], s.EMM[h2]), s.EMM[h3]))
	}
	return resp
}
