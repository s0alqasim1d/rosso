package json

import (
   "bytes"
   "encoding/json"
   "io"
)

// github.com/golang/go/blob/go1.20.3/src/encoding/xml/xml.go
func unmarshal(text, sep []byte, v any, before bool) error {
   _, text, found := bytes.Cut(text, sep)
   if !found {
      return io.EOF
   }
   if before {
      text = append(sep, text...)
   }
   dec := NewDecoder(bytes.NewReader(text))
   for {
      _, err := dec.Token()
      if err != nil {
         text = text[:dec.InputOffset()]
         return json.Unmarshal(text, v)
      }
   }
}

var (
   MarshalIndent = json.MarshalIndent
   NewDecoder = json.NewDecoder
)

func Cut(text, sep []byte, v any) error {
   return unmarshal(text, sep, v, false)
}

func Cut_Before(text, sep []byte, v any) error {
   return unmarshal(text, sep, v, true)
}
