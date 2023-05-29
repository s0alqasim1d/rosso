package hls

import (
   "strconv"
   "strings"
)

// github.com/golang/go/blob/go1.20.4/src/internal/types/testdata/check/slices.go
func Filter[T any](s []T, f func(T) bool) []T {
   var values []T
   for _, value := range s {
      if f(value) {
         values = append(values, value)
      }
   }
   return values
}

// github.com/golang/exp/blob/2e198f4/slices/slices.go
func Index_Func[T any](s []T, f func(T) bool) int {
   for i, value := range s {
      if f(value) {
         return i
      }
   }
   return -1
}

type Media struct {
   Group_ID string
   Type string
   Name string
   Characteristics string
   Raw_URI string
}

func (m Media) String() string {
   var b strings.Builder
   b.WriteString("group ID: ")
   b.WriteString(m.Group_ID)
   b.WriteString("\ntype: ")
   b.WriteString(m.Type)
   b.WriteString("\nname: ")
   b.WriteString(m.Name)
   if m.Characteristics != "" {
      b.WriteString("\ncharacteristics: ")
      b.WriteString(m.Characteristics)
   }
   return b.String()
}

type Stream struct {
   Bandwidth int64
   Raw_URI string
   Audio string
   Codecs string
   Resolution string
}

func (m Stream) String() string {
   var b []byte
   if m.Resolution != "" {
      b = append(b, "resolution: "...)
      b = append(b, m.Resolution...)
      b = append(b, '\n')
   }
   b = append(b, "bandwidth: "...)
   b = strconv.AppendInt(b, m.Bandwidth, 10)
   if m.Codecs != "" {
      b = append(b, "\ncodecs: "...)
      b = append(b, m.Codecs...)
   }
   if m.Audio != "" {
      b = append(b, "\naudio: "...)
      b = append(b, m.Audio...)
   }
   return string(b)
}
