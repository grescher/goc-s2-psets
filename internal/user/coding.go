package user

import (
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

func Encode(in []User) (out []byte) {
	return out
}

func Decode(r io.Reader) (out []User, err error) {
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
		user.Name = string(name)
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
