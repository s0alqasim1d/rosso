package dash

import (
   "encoding/xml"
   "io"
   "sort"
   "strconv"
   "strings"
)

type Adaptation struct {
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Lang string `xml:"lang,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Representation []Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
}

type Content_Protection struct {
   PSSH string `xml:"pssh"`
   Scheme_ID_URI string `xml:"schemeIdUri,attr"`
}

type Presentation struct {
   Period struct {
      Adaptation_Set []Adaptation `xml:"AdaptationSet"`
   }
}

func (p Presentation) Representation() []Representation {
   var reps []Representation
   for i, ada := range p.Period.Adaptation_Set {
      for _, rep := range ada.Representation {
         rep.Adaptation = &p.Period.Adaptation_Set[i]
         if rep.Codecs == "" {
            rep.Codecs = ada.Codecs
         }
         if rep.Content_Protection == nil {
            rep.Content_Protection = ada.Content_Protection
         }
         if rep.MIME_Type == "" {
            rep.MIME_Type = ada.MIME_Type
         }
         if rep.Segment_Template == nil {
            rep.Segment_Template = ada.Segment_Template
         }
         reps = append(reps, rep)
      }
   }
   return reps
}

type Segment struct {
   D int `xml:"d,attr"` // duration
   R int `xml:"r,attr"` // repeat
   T int `xml:"t,attr"` // time
}

type Segment_Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   Segment_Timeline struct {
      S []Segment
   } `xml:"SegmentTimeline"`
   Start_Number int `xml:"startNumber,attr"`
}

func Audio(r Representation) bool {
   return r.MIME_Type == "audio/mp4"
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

func Video(r Representation) bool {
   return r.MIME_Type == "video/mp4"
}

func New_Presentation(r io.Reader) (*Presentation, error) {
   pre := new(Presentation)
   err := xml.NewDecoder(r).Decode(pre)
   if err != nil {
      return nil, err
   }
   return pre, nil
}

func (s Segment) replace(ref, old string) string {
   return strings.Replace(ref, old, strconv.Itoa(s.T), 1)
}
