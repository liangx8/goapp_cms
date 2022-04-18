package utils

import (
	"hash"
	"log"
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

	pwd.encoder.Reset()
	pwd.encoder.Write(salt)
	pwd.encoder.Write([]byte(val))
	return append([]byte{byte(saltlen)}, pwd.encoder.Sum(salt)...)
}
func (pwd *PasswordKit) Verify(enc []byte, val string) bool {
	cn := uint(enc[0])
	salt := make([]byte, cn)
	copy(salt, enc[1:cn+1])
	epw := pwd.Create(salt, val)
	log.Printf("%p:%x", epw, epw)
	log.Printf("%p:%x", enc, enc)
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
