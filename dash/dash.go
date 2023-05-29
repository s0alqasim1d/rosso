package dash

import (
   "encoding/xml"
   "io"
   "sort"
   "strconv"
   "strings"
)

func Audio(r Representation) bool {
   return r.MIME_Type == "audio/mp4"
}

func Video(r Representation) bool {
   return r.MIME_Type == "video/mp4"
}

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

// github.com/golang/go/blob/go1.20.4/src/strings/strings.go
func Last_Index_Func[T any](s []T, f func(T) bool) int {
   i := len(s) - 1
   for i >= 0 {
      if f(s[i]) {
         return i
      }
      i--
   }
   return -1
}

// github.com/golang/exp/blob/2e198f4/slices/sort.go
func Sort_Func[T any](s []T, less func(a, b T) bool) {
   sort.Slice(s, func(i, j int) bool {
      return less(s[i], s[j])
   })
}

func New_Presentation(r io.Reader) (*Presentation, error) {
   pre := new(Presentation)
   err := xml.NewDecoder(r).Decode(pre)
   if err != nil {
      return nil, err
   }
   return pre, nil
}

func (r Representation) Ext() string {
   switch {
   case Audio(r):
      return ".m4a"
   case Video(r):
      return ".m4v"
   }
   return ""
}

func (r Representation) Initialization() string {
   return r.replace_ID(r.Segment_Template.Initialization)
}

func (r Representation) Media() []string {
   var start int
   if r.Segment_Template.Start_Number != nil {
      start = *r.Segment_Template.Start_Number
   }
   var refs []string
   for _, seg := range r.Segment_Template.Segment_Timeline.S {
      for seg.T = start; seg.R >= 0; seg.R-- {
         ref := r.replace_ID(r.Segment_Template.Media)
         if r.Segment_Template.Start_Number != nil {
            ref = seg.replace(ref, "$Number$")
            seg.T++
            start++
         } else {
            ref = seg.replace(ref, "$Time$")
            seg.T += seg.D
            start += seg.D
         }
         refs = append(refs, ref)
      }
   }
   return refs
}

func (r Representation) Role() string {
   if r.Adaptation.Role != nil {
      return r.Adaptation.Role.Value
   }
   return ""
}

func (r Representation) String() string {
   var b []byte
   b = append(b, "ID: "...)
   b = append(b, r.ID...)
   if r.Width >= 1 {
      b = append(b, "\n\twidth: "...)
      b = strconv.AppendInt(b, r.Width, 10)
   }
   if r.Height >= 1 {
      b = append(b, "\n\theight: "...)
      b = strconv.AppendInt(b, r.Height, 10)
   }
   if r.Bandwidth >= 1 {
      b = append(b, "\n\tbandwidth: "...)
      b = strconv.AppendInt(b, r.Bandwidth, 10)
   }
   if r.Codecs != "" {
      b = append(b, "\n\tcodecs: "...)
      b = append(b, r.Codecs...)
   }
   b = append(b, "\n\ttype: "...)
   b = append(b, r.MIME_Type...)
   if r.Adaptation.Role != nil {
      b = append(b, "\n\trole: "...)
      b = append(b, r.Adaptation.Role.Value...)
   }
   if r.Adaptation.Lang != "" {
      b = append(b, "\n\tlanguage: "...)
      b = append(b, r.Adaptation.Lang...)
   }
   return string(b)
}

func (r Representation) Widevine() *Content_Protection {
   for _, c := range r.Content_Protection {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return &c
      }
   }
   return nil
}

func (r Representation) replace_ID(ref string) string {
   return strings.Replace(ref, "$RepresentationID$", r.ID, 1)
}

func (s Segment) replace(ref, old string) string {
   return strings.Replace(ref, old, strconv.Itoa(s.T), 1)
}
