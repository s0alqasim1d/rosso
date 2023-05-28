package dash

import "sort"

func filter[T any](s []T, fn func(T) bool) []T {
   var values []T
   for _, value := range s {
      if fn(value) {
         values = append(values, value)
      }
   }
   return values
}

func Audio(s []Representation) []Representation {
   return filter(s, func(r Representation) bool {
      return r.MIME_Type == "audio/mp4"
   })
}

func Index_Bandwidth(s []Representation, min int64) int {
   sort.Slice(s, func(i, j int) bool {
      return s[i].Bandwidth < s[j].Bandwidth
   })
   return sort.Search(len(s), func(i int) bool {
      return s[i].Bandwidth >= min
   })
}

func Video(s []Representation) []Representation {
   return filter(s, func(r Representation) bool {
      return r.MIME_Type == "video/mp4"
   })
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
