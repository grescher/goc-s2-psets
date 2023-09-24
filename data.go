package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"strings"
)

const (
	ActiveMask    = 1 << 63
	AgeMask       = math.MaxUint64 ^ ActiveMask
	maxNumOfUsers = 8
	kgPerOz       = 0.0283495
	kgPerQq       = 100.0
)

type User struct {
	Name        string   // uint8(length) + [length]byte
	Age         uint8    // 1 bit bool (active field) + 63 bit uint (age field)
	ActiveIndex uint8    // (see above)
	Mass        float64  // regular float64
	Books       []string // uint8(all books length) + [length]byte, all books come as a single comma-separated string
}

func Users(r io.Reader) (out []User, err error) {
	for err != io.EOF && len(out) < maxNumOfUsers {
		var nameLength uint8
		if err = binary.Read(r, binary.BigEndian, &nameLength); err != nil {
			break
		}
		name := make([]byte, nameLength)
		if err = binary.Read(r, binary.BigEndian, &name); err != nil {
			break
		}
		var activeAndAge uint64
		if err = binary.Read(r, binary.BigEndian, &activeAndAge); err != nil {
			break
		}
		var active uint8
		if activeAndAge&ActiveMask > 0 {
			active = 1
		}
		var mass float64
		if err = binary.Read(r, binary.BigEndian, &mass); err != nil {
			break
		}
		var booksLen uint8
		if err = binary.Read(r, binary.BigEndian, &booksLen); err != nil {
			break
		}
		books := make([]byte, booksLen)
		if err = binary.Read(r, binary.BigEndian, &books); err != nil {
			break
		}

		var user User
		user.Name = strings.TrimSpace(string(name))
		user.ActiveIndex = active << len(out)
		user.Age = uint8(activeAndAge & AgeMask)
		user.Mass = verifyMass(mass)
		user.Books = strings.Split(string(books), ",")
		out = append(out, user)
	}
	if len(out) != maxNumOfUsers && err != io.EOF {
		return nil, err
	}
	return out, nil
}

func verifyMass(m float64) float64 {
	switch {
	case m > 0.0009 && m < 1: // quintals to kg
		m = m * kgPerQq
	case m > 620:
		m = m * kgPerOz // ounces to kg
	}
	return m
}

func Reader() io.Reader {
	return bytes.NewBuffer(
		[]byte{
			0x8, 0x4a, 0x6f, 0x68, 0x6e, 0x20, 0x44, 0x6f, 0x65, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1e, 0x40, 0x54, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x11, 0x48, 0x61, 0x72, 0x72, 0x79, 0x20, 0x50, 0x6f, 0x74, 0x74, 0x65, 0x72, 0x2c, 0x31, 0x39, 0x38, 0x34, 0x8, 0x4a, 0x61, 0x6b, 0x65, 0x20, 0x44, 0x6f, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0x40, 0x4e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x20, 0x4a, 0x61, 0x6e, 0x65, 0x20, 0x44, 0x6f, 0x65, 0x20, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x96, 0x3f, 0xe8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1c, 0x48, 0x61, 0x72, 0x72, 0x79, 0x20, 0x50, 0x6f, 0x74, 0x74, 0x65, 0x72, 0x2c, 0x47, 0x61, 0x6d, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x54, 0x68, 0x72, 0x6f, 0x6e, 0x65, 0x73, 0x1, 0x9, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5a, 0x40, 0xbf, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc, 0x48, 0x61,
			0x72, 0x72, 0x79, 0x20, 0x50, 0x6f, 0x74, 0x74, 0x65, 0x72, 0xd9, 0x56, 0x6d, 0x30, 0x77, 0x65, 0x45, 0x35, 0x47, 0x56, 0x58, 0x68, 0x55, 0x62, 0x6c, 0x4a, 0x57, 0x56, 0x30, 0x64, 0x53, 0x55, 0x46, 0x5a, 0x74, 0x4d, 0x57, 0x39, 0x57, 0x56, 0x6c, 0x6c, 0x33, 0x57, 0x6b, 0x52,
			0x53, 0x56, 0x31, 0x5a, 0x74, 0x65, 0x44, 0x42, 0x61, 0x52, 0x56, 0x59, 0x77, 0x56, 0x6d, 0x73, 0x78, 0x57, 0x47, 0x56, 0x47, 0x57, 0x6c, 0x5a, 0x69, 0x56, 0x45, 0x5a, 0x49, 0x57, 0x56, 0x64, 0x34, 0x53, 0x32, 0x4d, 0x78, 0x54, 0x6e, 0x4e, 0x69, 0x52, 0x30, 0x5a, 0x58, 0x56,
			0x6d, 0x35, 0x43, 0x62, 0x31, 0x5a, 0x73, 0x56, 0x6d, 0x46, 0x57, 0x4d, 0x56, 0x70, 0x57, 0x54, 0x56, 0x56, 0x57, 0x61, 0x47, 0x56, 0x71, 0x51, 0x54, 0x6b, 0x3d, 0xa, 0x56, 0x6d, 0x30, 0x77, 0x65, 0x45, 0x35, 0x47, 0x56, 0x58, 0x68, 0x55, 0x62, 0x6c, 0x4a, 0x57, 0x56, 0x30, 0x64, 0x53, 0x55, 0x46, 0x5a, 0x74, 0x4d, 0x57, 0x39, 0x57, 0x56, 0x6c, 0x6c, 0x33, 0x57, 0x6b, 0x52, 0x53, 0x56, 0x31, 0x5a, 0x74, 0x65, 0x44, 0x42, 0x61, 0x52, 0x56, 0x59, 0x77, 0x56, 0x6d, 0x73, 0x78, 0x57, 0x47, 0x56, 0x47, 0x57, 0x6c, 0x5a, 0x69, 0x56, 0x45, 0x5a, 0x49, 0x57, 0x56, 0x64, 0x34, 0x53, 0x32, 0x4d, 0x78, 0x54, 0x6e, 0x4e, 0x69, 0x52, 0x30, 0x5a, 0x58, 0x56, 0x6d, 0x35, 0x43, 0x62, 0x31, 0x5a, 0x73, 0x56, 0x6d, 0x46, 0x57, 0x4d, 0x56, 0x70, 0x57, 0x54, 0x56, 0x56, 0x57, 0x61, 0x47, 0x56, 0x71, 0x51, 0x54, 0x6b, 0x3d, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x54, 0x68, 0x65, 0x20, 0x48, 0x75, 0x6e, 0x67, 0x65, 0x72, 0x20, 0x47, 0x61, 0x6d, 0x65, 0x73, 0x8, 0x0, 0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b, 0x4d, 0x6f, 0x62, 0x79, 0x20, 0x44, 0x69, 0x63, 0x6b, 0x2c, 0x49, 0x74, 0x2c, 0x54, 0x68, 0x65, 0x20, 0x47, 0x72, 0x65, 0x65, 0x6e, 0x20, 0x4d, 0x69, 0x6c, 0x65,
		},
	)
}
