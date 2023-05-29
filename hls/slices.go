package hls

import (
   "strconv"
   "strings"
)

type Medium struct {
   Group_ID string
   Name string
   Raw_URI string
   Type string
   Characteristics string
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
   b = append(b, "bandwidth: "...)
   b = strconv.AppendInt(b, m.Bandwidth, 10)
   if m.Resolution != "" {
      b = append(b, "\n\tresolution: "...)
      b = append(b, m.Resolution...)
   }
   if m.Codecs != "" {
      b = append(b, "\n\tcodecs: "...)
      b = append(b, m.Codecs...)
   }
   if m.Audio != "" {
      b = append(b, "\n\taudio: "...)
      b = append(b, m.Audio...)
   }
   return string(b)
}

func (m Medium) String() string {
   var b strings.Builder
   b.WriteString("group ID: ")
   b.WriteString(m.Group_ID)
   b.WriteString("\n\ttype: ")
   b.WriteString(m.Type)
   b.WriteString("\n\tname: ")
   b.WriteString(m.Name)
   if m.Characteristics != "" {
      b.WriteString("\n\tcharacteristics: ")
      b.WriteString(m.Characteristics)
   }
   return b.String()
}

type Media []Medium

type Streams []Stream

func filter[T Mixed](slice []T, callback func(T) bool) []T {
   var carry []T
   for _, item := range slice {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func index[T Mixed](slice []T, callback func(T, T) bool) int {
   carry := -1
   for i, item := range slice {
      if carry == -1 || callback(slice[carry], item) {
         carry = i
      }
   }
   return carry
}

func (m Media) Filter(f func(Medium) bool) Media {
   return filter(m, f)
}

func (m Streams) Filter(f func(Stream) bool) Streams {
   return filter(m, f)
}

func (m Media) Index(f func(a, b Medium) bool) int {
   return index(m, f)
}

func (m Streams) Index(f func(a, b Stream) bool) int {
   return index(m, f)
}

func (m Streams) Bandwidth(v int64) int {
   distance := func(a Stream) int64 {
      if a.Bandwidth > v {
         return a.Bandwidth - v
      }
      return v - a.Bandwidth
   }
   return m.Index(func(carry, item Stream) bool {
      return distance(item) < distance(carry)
   })
}
