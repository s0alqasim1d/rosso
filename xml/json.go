package xml

import (
   "bytes"
   "encoding/xml"
   "io"
)

func decoder(data []byte) *xml.Decoder {
   dec := xml.NewDecoder(bytes.NewReader(data))
   dec.AutoClose = xml.HTMLAutoClose
   dec.Strict = false
   return dec
}

type Scanner struct {
   Data []byte
   Sep []byte
}

func (s Scanner) Decode(val any) error {
   data := append(s.Sep, s.Data...)
   dec := decoder(data)
   for {
      _, err := dec.Token()
      if err != nil {
         high := dec.InputOffset()
         return decoder(data[:high]).Decode(val)
      }
   }
}

func (s *Scanner) Scan() bool {
   var found bool
   _, s.Data, found = bytes.Cut(s.Data, s.Sep)
   return found
}

func Indent(dst io.Writer, src io.Reader, prefix, indent string) error {
   decode := xml.NewDecoder(src)
   encode := xml.NewEncoder(dst)
   encode.Indent(prefix, indent)
   for {
      token, err := decode.Token()
      if err == io.EOF {
         return encode.Flush()
      }
      if err != nil {
         return err
      }
      data, ok := token.(xml.CharData)
      if ok {
         token = xml.CharData(bytes.TrimSpace(data))
      }
      if err := encode.EncodeToken(token); err != nil {
         return err
      }
   }
}
