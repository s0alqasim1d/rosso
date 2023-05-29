package hls

import (
   "2a.pages.dev/rosso/slices"
   "strconv"
   "strings"
)

type Media []Medium

func (m Media) Filter(f func(Medium) bool) Media {
   return slices.Filter(m, f)
}

func (m Media) Index(f func(Medium) bool) int {
   return slices.Index_Func(m, f)
}

type Medium struct {
   Group_ID string
   Type string
   Name string
   Characteristics string
   Raw_URI string
}

func (m Medium) String() string {
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

type Streams []Stream

func (s Streams) Filter(f func(Stream) bool) Streams {
   return slices.Filter(s, f)
}

func (s Streams) Index(f func(Stream) bool) int {
   return slices.Index_Func(s, f)
}

func (s Streams) Last_Index(f func(Stream) bool) int {
   return slices.Last_Index_Func(s, f)
}

func (s Streams) Sort(f func(a, b Stream) bool) {
   slices.Sort_Func(s, f)
}
