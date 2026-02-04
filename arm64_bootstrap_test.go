package beignet

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"runtime"
	"testing"
)

func TestBuildArm64Bootstrap_StrictCompatibility(t *testing.T) {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("darwin/arm64 only")
	}

	cases := [][4]uint64{
		{0, 0, 0, 0},
		{1, 2, 3, 4},
		{0xffff, 0x10000, 0x12345678, 0xdeadbeef},
		{0x1122334455667788, 0x8877665544332211, 0x0, 0xffffffffffffffff},
		{0x1000, 0x2000, 0x3000, 0x4000},
	}

	r := rand.New(rand.NewSource(1337))
	for i := 0; i < 32; i++ {
		cases = append(cases, [4]uint64{r.Uint64(), r.Uint64(), r.Uint64(), r.Uint64()})
	}

	for i, tc := range cases {
		got, err := buildArm64Bootstrap(tc[0], tc[1], tc[2], tc[3])
		if err != nil {
			t.Fatalf("case %d: buildArm64Bootstrap error: %v", i, err)
		}
		want := buildArm64BootstrapHardcoded(tc[0], tc[1], tc[2], tc[3])

		if len(got) != arm64BootstrapLen {
			t.Fatalf("case %d: bootstrap length got=%d want=%d", i, len(got), arm64BootstrapLen)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("case %d: bootstrap bytes mismatch", i)
		}
	}
}

func buildArm64BootstrapHardcoded(payloadOffset, payloadSize, symbolOffset, loaderEntryOffsetAbs uint64) []byte {
	var out []byte

	// adr x9, #0
	out = appendU32LEHardcoded(out, encADRHardcoded(9, 0))

	// x0 = base + payloadOffset
	out = appendMovImm64XHardcoded(out, 0, payloadOffset)
	out = appendU32LEHardcoded(out, encADDRegXHardcoded(0, 0, 9))

	// x1 = payloadSize
	out = appendMovImm64XHardcoded(out, 1, payloadSize)

	// x2 = base + symbolOffset
	out = appendMovImm64XHardcoded(out, 2, symbolOffset)
	out = appendU32LEHardcoded(out, encADDRegXHardcoded(2, 2, 9))

	// x16 = base + loaderEntryOffsetAbs
	out = appendMovImm64XHardcoded(out, 16, loaderEntryOffsetAbs)
	out = appendU32LEHardcoded(out, encADDRegXHardcoded(16, 16, 9))

	// br x16
	out = appendU32LEHardcoded(out, encBRHardcoded(16))

	return out
}

func appendMovImm64XHardcoded(out []byte, rd uint8, imm uint64) []byte {
	out = appendU32LEHardcoded(out, encMOVZXHardcoded(rd, uint16(imm&0xffff), 0))
	out = appendU32LEHardcoded(out, encMOVKXHardcoded(rd, uint16((imm>>16)&0xffff), 16))
	out = appendU32LEHardcoded(out, encMOVKXHardcoded(rd, uint16((imm>>32)&0xffff), 32))
	out = appendU32LEHardcoded(out, encMOVKXHardcoded(rd, uint16((imm>>48)&0xffff), 48))
	return out
}

func appendU32LEHardcoded(out []byte, v uint32) []byte {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	return append(out, b[:]...)
}

// ADR (immediate).
// Base opcode 0x10000000. imm is a signed 21-bit byte offset.
func encADRHardcoded(rd uint8, imm int32) uint32 {
	immlo := uint32(imm) & 0x3
	immhi := (uint32(imm) >> 2) & 0x7ffff
	return 0x10000000 | (immlo << 29) | (immhi << 5) | uint32(rd)
}

// MOVZ Xd, imm16, LSL shift
func encMOVZXHardcoded(rd uint8, imm16 uint16, shift uint8) uint32 {
	hw := uint32(shift / 16)
	return 0xD2800000 | (hw << 21) | (uint32(imm16) << 5) | uint32(rd)
}

// MOVK Xd, imm16, LSL shift
func encMOVKXHardcoded(rd uint8, imm16 uint16, shift uint8) uint32 {
	hw := uint32(shift / 16)
	return 0xF2800000 | (hw << 21) | (uint32(imm16) << 5) | uint32(rd)
}

// ADD Xd, Xn, Xm (shifted register), shift=LSL #0
func encADDRegXHardcoded(rd, rn, rm uint8) uint32 {
	return 0x8B000000 | (uint32(rm) << 16) | (uint32(rn) << 5) | uint32(rd)
}

// BR Xn
func encBRHardcoded(rn uint8) uint32 {
	return 0xD61F0000 | (uint32(rn) << 5)
}

