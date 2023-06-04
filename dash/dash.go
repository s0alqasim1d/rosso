package dash

import (
   "encoding/xml"
   "fmt"
   "io"
   "strings"
)

// amcplus.com
type Adaptation_Set struct {
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Lang string `xml:"lang,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Representation []Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
}

type Representation struct {
   // roku.com
   Adaptation_Set *Adaptation_Set
   // roku.com
   Bandwidth int64 `xml:"bandwidth,attr"`
   // roku.com
   Codecs string `xml:"codecs,attr"`
   // roku.com
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   // roku.com
   Height int64 `xml:"height,attr"`
   // roku.com
   ID string `xml:"id,attr"`
   // paramountplus.com
   MIME_Type string `xml:"mimeType,attr"`
   // roku.com
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   // roku.com
   Width int64 `xml:"width,attr"`
}

// roku.com
type Content_Protection struct {
   PSSH string `xml:"pssh"`
   Scheme_ID_URI string `xml:"schemeIdUri,attr"`
}

// roku.com
type Segment struct {
   D int `xml:"d,attr"` // duration
   R int `xml:"r,attr"` // repeat
   T int `xml:"t,attr"` // time
}

// roku.com
type Segment_Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   Segment_Timeline struct {
      S []Segment
   } `xml:"SegmentTimeline"`
   Start_Number int `xml:"startNumber,attr"`
}

func (r Representation) Widevine() *Content_Protection {
   for _, c := range r.get_content_protection() {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return &c
      }
   }
   return nil
}

func (r Representation) get_content_protection() []Content_Protection {
   if r.Content_Protection != nil {
      return r.Content_Protection
   }
   return r.Adaptation_Set.Content_Protection
}

func Audio(r Representation) bool {
   return r.get_MIME_type() == "audio/mp4"
}

func Video(r Representation) bool {
   return r.get_MIME_type() == "video/mp4"
}

func (r Representation) String() string {
   var s []string
   if r.Width >= 1 {
      s = append(s, fmt.Sprint("width: ", r.Width))
   }
   if r.Height >= 1 {
      s = append(s, fmt.Sprint("height: ", r.Height))
   }
   if r.Bandwidth >= 1 {
      s = append(s, fmt.Sprint("bandwidth: ", r.Bandwidth))
   }
   if r.Codecs != "" {
      s = append(s, "codecs: " + r.Codecs)
   }
   s = append(s, "type: " + r.get_MIME_type())
   if r.Adaptation_Set.Role != nil {
      s = append(s, "role: " + r.Adaptation_Set.Role.Value)
   }
   if r.Adaptation_Set.Lang != "" {
      s = append(s, "language: " + r.Adaptation_Set.Lang)
   }
   return strings.Join(s, "\n")
}

func (r Representation) get_MIME_type() string {
   if r.MIME_Type != "" {
      return r.MIME_Type
   }
   return r.Adaptation_Set.MIME_Type
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

func Not[E any](fn func(E) bool) func(E) bool {
   return func(value E) bool {
      return !fn(value)
   }
}

func (r Representation) replace_ID(ref string) string {
   return strings.Replace(ref, "$RepresentationID$", r.ID, 1)
}

func (s Segment) replace(ref, old string) string {
   return strings.Replace(ref, old, fmt.Sprint(s.T), 1)
}

func Representations(r io.Reader) ([]Representation, error) {
   var s struct {
      Period struct {
         Adaptation_Set []Adaptation_Set `xml:"AdaptationSet"`
      }
   }
   err := xml.NewDecoder(r).Decode(&s)
   if err != nil {
      return nil, err
   }
   var reps []Representation
   for _, ada := range s.Period.Adaptation_Set {
      ada := ada
      for _, rep := range ada.Representation {
         rep.Adaptation_Set = &ada
         reps = append(reps, rep)
      }
   }
   return reps, nil
}

func (r Representation) Initialization() string {
   return r.replace_ID(r.get_segment_template().Initialization)
}

func (r Representation) get_segment_template() *Segment_Template {
   if r.Segment_Template != nil {
      return r.Segment_Template
   }
   return r.Adaptation_Set.Segment_Template
}

func (r Representation) Media() []string {
   var refs []string
   tem := r.get_segment_template()
   for _, seg := range tem.Segment_Timeline.S {
      seg.T = tem.Start_Number
      for seg.R >= 0 {
         {
            ref := r.replace_ID(tem.Media)
            ref = seg.replace(ref, "$Number$")
            refs = append(refs, ref)
         }
         tem.Start_Number++
         seg.T++
         seg.R--
      }
   }
   return refs
}
