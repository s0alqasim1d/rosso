package xml

import (
   "bytes"
   "encoding/xml"
   "io"
)

// github.com/golang/go/blob/go1.20.3/src/encoding/xml/xml.go
func decode(text, sep []byte, v any, before bool) error {
   _, text, found := bytes.Cut(text, sep)
   if !found {
      return io.EOF
   }
   if before {
      text = append(sep, text...)
   }
   dec := new_decoder(text)
   for {
      _, err := dec.Token()
      if err != nil {
         text = text[:dec.InputOffset()]
         return new_decoder(text).Decode(v)
      }
   }
}

func Cut(text, sep []byte, v any) error {
   return decode(text, sep, v, false)
}

func Cut_Before(text, sep []byte, v any) error {
   return decode(text, sep, v, true)
}

func Indent(dst io.Writer, src io.Reader, prefix, indent string) error {
   dec := xml.NewDecoder(src)
   enc := xml.NewEncoder(dst)
   enc.Indent(prefix, indent)
   for {
      token, err := dec.Token()
      if err == io.EOF {
         return enc.Flush()
      }
      if err != nil {
         return err
      }
      text, ok := token.(xml.CharData)
      if ok {
         token = xml.CharData(bytes.TrimSpace(text))
      }
      if err := enc.EncodeToken(token); err != nil {
         return err
      }
   }
}

func new_decoder(text []byte) *xml.Decoder {
   dec := xml.NewDecoder(bytes.NewReader(text))
   dec.AutoClose = xml.HTMLAutoClose
   dec.Strict = false
   return dec
}
