package merkle

import (
	"encoding/hex"
	"math/big"
	"strings"
)

// EncodeLeaf encodes (index, address, amount) like Solidity abi.encodePacked(uint256,address,uint256)
// and returns keccak256 hash.
func EncodeLeaf(index uint64, addressHex string, amountStr string) [32]byte {
	// uint256 index (32 bytes big-endian)
	idxBytes := uintToBytes(index)
	// address (20 bytes)
	addrBytes := addressToBytes(addressHex)
	// uint256 amount (32 bytes big-endian)
	amtBytes := decimalToUint256Bytes(amountStr)

	data := make([]byte, 0, 32+20+32)
	data = append(data, idxBytes...)
	data = append(data, addrBytes...)
	data = append(data, amtBytes...)
	return HashKeccak256(data)
}

func uintToBytes(v uint64) []byte {
	b := make([]byte, 32)
	for i := 0; i < 8; i++ {
		b[31-i] = byte(v & 0xff)
		v >>= 8
	}
	return b
}

func addressToBytes(addr string) []byte {
	addr = strings.TrimPrefix(strings.ToLower(addr), "0x")
	if len(addr) != 40 {
		// return 20 zero bytes
		return make([]byte, 20)
	}
	bz, err := hex.DecodeString(addr)
	if err != nil || len(bz) != 20 {
		return make([]byte, 20)
	}
	return bz
}

func decimalToUint256Bytes(amount string) []byte {
	n := new(big.Int)
	_, ok := n.SetString(amount, 10)
	if !ok {
		n = big.NewInt(0)
	}
	bz := n.Bytes()
	out := make([]byte, 32)
	copy(out[32-len(bz):], bz)
	return out
}
