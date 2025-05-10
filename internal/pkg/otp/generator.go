package otp

import (
	"encoding/binary"
	"fmt"
	"math/rand"
)

func GenerateOTP() string {
	var buf [4]byte // 4 bytes = 32 bits
	_, err := rand.Read(buf[:])
	if err != nil {
		panic("OTP generation failed: " + err.Error())
	}
	n := binary.LittleEndian.Uint32(buf[:]) % 1000000
	return fmt.Sprintf("%06d", n)
}
