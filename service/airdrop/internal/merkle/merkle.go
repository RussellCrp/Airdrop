package merkle

import (
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/sha3"
)

// HashKeccak256 returns keccak256 hash of data.
func HashKeccak256(data []byte) [32]byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}

// FromHex32 decodes 0x-prefixed hex string into [32]byte.
func FromHex32(s string) ([32]byte, error) {
	var out [32]byte
	if len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X") {
		s = s[2:]
	}
	bz, err := hex.DecodeString(s)
	if err != nil {
		return out, err
	}
	if len(bz) != 32 {
		return out, errors.New("invalid length")
	}
	copy(out[:], bz)
	return out, nil
}

func ToHex32(b [32]byte) string {
	return "0x" + hex.EncodeToString(b[:])
}

type Tree struct {
	Levels [][][32]byte // level 0 is leaves
}

// Build constructs a Merkle tree from leaves.
func Build(leaves [][32]byte) *Tree {
	if len(leaves) == 0 {
		return &Tree{Levels: [][][32]byte{}}
	}
	levels := make([][][32]byte, 0)
	cur := make([][32]byte, len(leaves))
	copy(cur, leaves)
	levels = append(levels, cur)

	for len(cur) > 1 {
		var next [][32]byte
		for i := 0; i < len(cur); i += 2 {
			if i+1 == len(cur) {
				// duplicate last
				next = append(next, HashPair(cur[i], cur[i]))
			} else {
				next = append(next, HashPair(cur[i], cur[i+1]))
			}
		}
		cur = next
		levels = append(levels, cur)
	}
	return &Tree{Levels: levels}
}

func HashPair(a, b [32]byte) [32]byte {
	buf := make([]byte, 64)
	copy(buf[:32], a[:])
	copy(buf[32:], b[:])
	return HashKeccak256(buf)
}

func (t *Tree) Root() ([32]byte, bool) {
	if len(t.Levels) == 0 {
		return [32]byte{}, false
	}
	top := t.Levels[len(t.Levels)-1]
	if len(top) == 0 {
		return [32]byte{}, false
	}
	return top[0], true
}

// Proof returns Merkle proof for leaf index.
func (t *Tree) Proof(index int) [][32]byte {
	if len(t.Levels) == 0 {
		return nil
	}
	proof := make([][32]byte, 0)
	curIndex := index
	for level := 0; level < len(t.Levels)-1; level++ {
		nodes := t.Levels[level]
		if curIndex%2 == 0 {
			// right sibling
			sib := curIndex + 1
			if sib >= len(nodes) {
				sib = curIndex
			}
			proof = append(proof, nodes[sib])
		} else {
			// left sibling
			sib := curIndex - 1
			proof = append(proof, nodes[sib])
		}
		curIndex /= 2
	}
	return proof
}
