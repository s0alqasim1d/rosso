package dash

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
