package dash

import (
   "encoding/xml"
   "fmt"
   "io"
   "strings"
)

func (s Segment_Template) Replace() []string {
   var refs []string
   for _, seg := range s.Segment_Timeline.S {
      seg.T = s.Start_Number
      for seg.R >= 0 {
         {
            ref := strings.Replace(s.Media, "$Number$", fmt.Sprint(seg.T), 1)
            refs = append(refs, ref)
         }
         s.Start_Number++
         seg.R--
         seg.T++
      }
   }
   return refs
}

// amcplus.com
type Adaptation_Set struct {
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Lang string `xml:"lang,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   Representation []Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
}

// roku.com
type Content_Protection struct {
   PSSH string `xml:"pssh"`
   Scheme_ID_URI string `xml:"schemeIdUri,attr"`
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

// roku.com
type Segment struct {
   D int `xml:"d,attr"` // duration
   R int `xml:"r,attr"` // repeat
   T int `xml:"t,attr"` // time
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
         if rep.Content_Protection == nil {
            rep.Content_Protection = ada.Content_Protection
         }
         if rep.MIME_Type == nil {
            rep.MIME_Type = &ada.MIME_Type
         }
         if rep.Segment_Template == nil {
            rep.Segment_Template = ada.Segment_Template
         }
         rep.replace(&rep.Segment_Template.Initialization)
         rep.replace(&rep.Segment_Template.Media)
         reps = append(reps, rep)
      }
   }
   return reps, nil
}

func (r Representation) replace(s *string) {
   *s = strings.Replace(*s, "$RepresentationID$", r.ID, 1)
}

type Representation struct {
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
   MIME_Type *string `xml:"mimeType,attr"`
   // roku.com
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   // roku.com
   Width int64 `xml:"width,attr"`
   // roku.com
   Adaptation_Set *Adaptation_Set
}

func Audio(r Representation) bool {
   return *r.MIME_Type == "audio/mp4"
}

func Video(r Representation) bool {
   return *r.MIME_Type == "video/mp4"
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
   s = append(s, "type: " + *r.MIME_Type)
   if r.Adaptation_Set.Role != nil {
      s = append(s, "role: " + r.Adaptation_Set.Role.Value)
   }
   if r.Adaptation_Set.Lang != "" {
      s = append(s, "language: " + r.Adaptation_Set.Lang)
   }
   return strings.Join(s, "\n")
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

func (r Representation) Widevine() string {
   for _, c := range r.Content_Protection {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return c.PSSH
      }
   }
   return ""
}
