package beignet

import "encoding/binary"

// buildArm64Bootstrap returns a small aarch64 stub which:
// - sets x0 = base + payloadOffset
// - sets x1 = payloadSize
// - sets x2 = base + symbolOffset
// - jumps to base + loaderEntryOffsetAbs using br
//
// The loader expects payload pointer/size in x0/x1 and the entry symbol pointer in x2.
func buildArm64Bootstrap(payloadOffset, payloadSize, symbolOffset, loaderEntryOffsetAbs uint64) []byte {
	var out []byte

	// adr x9, #0
	out = appendU32LE(out, encADR(9, 0))

	// x0 = base + payloadOffset
	out = appendMovImm64X(out, 0, payloadOffset)
	out = appendU32LE(out, encADDRegX(0, 0, 9))

	// x1 = payloadSize
	out = appendMovImm64X(out, 1, payloadSize)

	// x2 = base + symbolOffset
	out = appendMovImm64X(out, 2, symbolOffset)
	out = appendU32LE(out, encADDRegX(2, 2, 9))

	// x16 = base + loaderEntryOffsetAbs
	out = appendMovImm64X(out, 16, loaderEntryOffsetAbs)
	out = appendU32LE(out, encADDRegX(16, 16, 9))

	// br x16
	out = appendU32LE(out, encBR(16))

	return out
}

func appendMovImm64X(out []byte, rd uint8, imm uint64) []byte {
	out = appendU32LE(out, encMOVZX(rd, uint16(imm&0xffff), 0))
	out = appendU32LE(out, encMOVKX(rd, uint16((imm>>16)&0xffff), 16))
	out = appendU32LE(out, encMOVKX(rd, uint16((imm>>32)&0xffff), 32))
	out = appendU32LE(out, encMOVKX(rd, uint16((imm>>48)&0xffff), 48))
	return out
}

func appendU32LE(out []byte, v uint32) []byte {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	return append(out, b[:]...)
}

// ADR (immediate).
// Base opcode 0x10000000. imm is a signed 21-bit byte offset.
func encADR(rd uint8, imm int32) uint32 {
	immlo := uint32(imm) & 0x3
	immhi := (uint32(imm) >> 2) & 0x7ffff
	return 0x10000000 | (immlo << 29) | (immhi << 5) | uint32(rd)
}

// MOVZ Xd, imm16, LSL shift
func encMOVZX(rd uint8, imm16 uint16, shift uint8) uint32 {
	hw := uint32(shift / 16)
	return 0xD2800000 | (hw << 21) | (uint32(imm16) << 5) | uint32(rd)
}

// MOVK Xd, imm16, LSL shift
func encMOVKX(rd uint8, imm16 uint16, shift uint8) uint32 {
	hw := uint32(shift / 16)
	return 0xF2800000 | (hw << 21) | (uint32(imm16) << 5) | uint32(rd)
}

// ADD Xd, Xn, Xm (shifted register), shift=LSL #0
func encADDRegX(rd, rn, rm uint8) uint32 {
	return 0x8B000000 | (uint32(rm) << 16) | (uint32(rn) << 5) | uint32(rd)
}

// BR Xn
func encBR(rn uint8) uint32 {
	return 0xD61F0000 | (uint32(rn) << 5)
}
