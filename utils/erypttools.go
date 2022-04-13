package utils

import (
	"hash"
)

type (
	PasswordKit struct {
		encoder hash.Hash
	}
)

func NewPasswordKit(ha hash.Hash) *PasswordKit {
	return &PasswordKit{ha}
}
func (pwd *PasswordKit) Create(salt []byte, val string) []byte {
	saltlen := len(salt)
	if saltlen > 255 {
		panic("length of solt is greate than 255")
	}
	bu := []byte{byte(saltlen)}

	pwd.encoder.Write(salt)
	pwd.encoder.Write([]byte(val))

	return append(bu, pwd.encoder.Sum(salt)...)

}
func (pwd *PasswordKit) Verify(enc []byte, val string) bool {
	cn := int(enc[0])
	salt := enc[1 : cn+1]
	epw := pwd.Create(salt, val)
	if len(epw) == len(enc) {
		for ix, bb := range epw {
			if enc[ix] != bb {
				return false
			}
		}
		return true
	} else {
		return false
	}

}
