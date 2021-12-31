package randutil

import "math/rand"

// charSet
var AlphaNumLowerCase string

func init() {
	initAlphaNum()
}

func initAlphaNum() {
	s := make([]byte, 36)
	i := 0
	for c:='0';c<='9';c++ {
		s[i] = byte(c)
		i++
	}
	for c:='a';c<='z';c++ {
		s[i] = byte(c)
		i++
	}
	AlphaNumLowerCase = string(s)
}

func RandStr(charSet string, length int) string {
	var s []byte
	charSetSize := len(charSet)
	if len(charSet) > 0 {
		s = make([]byte, length)
		for i:=0; i<length; i++ {
			s[i] = charSet[rand.Int63n(int64(charSetSize))]
		}
	}
	return string(s)
}