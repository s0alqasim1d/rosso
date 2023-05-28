package dash

import "sort"

func Audio(r Representation) bool {
   return r.MIME_Type == "audio/mp4"
}

func Video(r Representation) bool {
   return r.MIME_Type == "video/mp4"
}

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Height int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   Width int64 `xml:"width,attr"`
}

// github.com/golang/go/blob/go1.20.4/src/internal/types/testdata/check/slices.go
func Filter[T any](s []T, fn func(T) bool) []T {
   var values []T
   for _, value := range s {
      if fn(value) {
         values = append(values, value)
      }
   }
   return values
}

// github.com/golang/exp/blob/2e198f4/slices/slices.go
func Index_Func[E any](s []E, fn func(E) bool) int {
   for index, value := range s {
      if fn(value) {
         return index
      }
   }
   return -1
}

// github.com/golang/exp/blob/2e198f4/slices/sort.go
func Sort_Func[E any](s []E, less func(a, b E) bool) {
   sort.Slice(s, func(i, j int) bool {
      return less(s[i], s[j])
   })
}
