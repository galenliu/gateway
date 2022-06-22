package dnssd

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func btoh(i byte) byte {
	if i > 9 {
		return 0x61 + (i - 10)
	}
	return 0x30 + i
}

// randomHex returns a random hex string.
func randHex() string {
	var b [16]byte
	// Read might block
	// > crypto/rand: blocked for 60 seconds waiting to read random data from the kernel
	// > https://github.com/golang/go/commit/1961d8d72a53e780effa18bfa8dbe4e4282df0b2
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	var out [32]byte
	for i := 0; i < len(b); i++ {
		out[i*2] = btoh((b[i] >> 4) & 0xF)
		out[i*2+1] = btoh(b[i] & 0xF)
	}
	return string(out[:])
}

// mac48Address returns a MAC-48-like address from the argument string
func mac48Address(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	result := h.Sum(nil)

	var c []string
	c = append(c, toHex(result[0]))
	c = append(c, toHex(result[1]))
	c = append(c, toHex(result[2]))
	c = append(c, toHex(result[3]))
	c = append(c, toHex(result[4]))
	c = append(c, toHex(result[5]))

	// setup id needs the mac address in upper case
	return strings.ToUpper(strings.Join(c, ":"))
}

func toHex(b byte) string {
	return hex.EncodeToString([]byte{b})
}
