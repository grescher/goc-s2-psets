// This package contains the encoding/decoding methods for a slice of User type.
// Schema:
//
//		 Name               uint8(length) + [length]byte
//		 ActiveIndex | Age  uint64: 63-bit bool (active field) | 62-0 bits uint (age field)
//		 Mass               float64
//		 Books              uint8(all books length) + [length]byte
//	                     all books come as a single comma-separated string
package user

import (
	"encoding/binary"
	"io"
	"math"
	"strings"
)

const (
	ActiveMask    uint64 = 1 << 63
	AgeMask       uint64 = math.MaxUint64 ^ ActiveMask
	maxNumOfUsers        = 8
	kgPerOz              = 0.0283495
	kgPerQq              = 100.0
)

func Encode(w io.Writer, users []User) (err error) {
	for _, u := range users {
		// Encoding of the Name field.
		nameLen := uint8(len(u.Name))
		if err = binary.Write(w, binary.BigEndian, nameLen); err != nil {
			break
		}
		if err = binary.Write(w, binary.BigEndian, []byte(u.Name)); err != nil {
			break
		}

		// Encoding of the Active and Age fields.
		var activeAndAge uint64
		if u.ActiveIndex > 0 {
			activeAndAge = ActiveMask
		}
		activeAndAge |= uint64(u.Age)
		if err = binary.Write(w, binary.BigEndian, activeAndAge); err != nil {
			break
		}

		// Encoding of the Mass field.
		if err = binary.Write(w, binary.BigEndian, u.Mass); err != nil {
			break
		}

		// Encoding of the Books field.
		books := strings.Join(u.Books, ",")
		booksLen := uint8(len(books))
		if err = binary.Write(w, binary.BigEndian, booksLen); err != nil {
			break
		}
		if err = binary.Write(w, binary.BigEndian, []byte(books)); err != nil {
			break
		}
	}
	if err != nil {
		return err
	}

	return nil
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
