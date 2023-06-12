package strconv

import (
   "encoding/hex"
   "errors"
   "unicode"
   "unicode/utf8"
)

const escape_character = '~'

var error_escape = errors.New("invalid printable escape")

func Encode(src []byte) string {
   var dst []byte
   for len(src) >= 1 {
      r, size := decode_rune(src)
      s := src[:size]
      if r == utf8.RuneError && size == 1 {
         var d [2]byte
         hex.Encode(d[:], s)
         dst = append(dst, escape_character)
         dst = append(dst, d[:]...)
      } else {
         dst = append(dst, s...)
      }
      src = src[size:]
   }
   return string(dst)
}

func decode(src string) ([]byte, error) {
   var dst []byte
   for len(src) >= 1 {
      s := src[0]
      if s == escape_character {
         if len(src) <= 2 {
            return nil, error_escape
         }
         d, err := hex.DecodeString(src[1:3])
         if err != nil {
            return nil, err
         }
         dst = append(dst, d...)
         src = src[3:]
      } else {
         dst = append(dst, s)
         src = src[1:]
      }
   }
   return dst, nil
}

// I originally used:
// mimesniff.spec.whatwg.org#binary-data-byte
// but it fails with:
// fileformat.info/info/unicode/char/1b
func Binary_Data(r rune) bool {
   // this needs to be first because newline is a control character
   if unicode.IsSpace(r) {
      return false
   }
   // we could also use unicode.IsGraphic or unicode.IsPrint, but they call
   // unicode.In
   return unicode.IsControl(r)
}

func decode_rune(p []byte) (rune, int) {
   r, size := utf8.DecodeRune(p)
   if r == escape_character || Binary_Data(r) {
      return utf8.RuneError, 1
   }
   return r, size
}
