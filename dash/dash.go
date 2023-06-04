package dash

import (
   "encoding/xml"
   "fmt"
   "io"
   "strings"
)

func (r Represent) replace_ID(ref string) string {
   return strings.Replace(ref, "$RepresentationID$", r.ID, 1)
}

func (r Represent) Initialization() string {
   return r.replace_ID(r.Segment_Template.Initialization)
}

func (r Represent) Media() []string {
   var refs []string
   for _, seg := range r.Segment_Template.Segment_Timeline.S {
      seg.T = r.Segment_Template.Start_Number
      for seg.R >= 0 {
         {
            ref := r.replace_ID(r.Segment_Template.Media)
            ref = seg.replace(ref, "$Number$")
            refs = append(refs, ref)
         }
         r.Segment_Template.Start_Number++
         seg.T++
         seg.R--
      }
   }
   return refs
}

func (s Segment) replace(ref, old string) string {
   return strings.Replace(ref, old, fmt.Sprint(s.T), 1)
}

func (r Represent) String() string {
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
   s = append(s, "type: " + r.MIME_Type)
   if r.Adaptation.Role != nil {
      s = append(s, "role: " + r.Adaptation.Role.Value)
   }
   if r.Adaptation.Lang != "" {
      s = append(s, "language: " + r.Adaptation.Lang)
   }
   return strings.Join(s, "\n")
}

func (r Represent) Widevine() *Content_Protection {
   for _, c := range r.Content_Protection {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return &c
      }
   }
   return nil
}

func Audio(r Represent) bool {
   return r.MIME_Type == "audio/mp4"
}

func Video(r Represent) bool {
   return r.MIME_Type == "video/mp4"
}

func (r Represent) Ext() string {
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

func New_Media(r io.Reader) (*Media, error) {
   med := new(Media)
   err := xml.NewDecoder(r).Decode(med)
   if err != nil {
      return nil, err
   }
   return med, nil
}

func (m Media) Represents() []Represent {
   var reps []Represent
   for i, ada := range m.Period.Adaptation_Set {
      for _, rep := range ada.Representation {
         rep.Adaptation = &m.Period.Adaptation_Set[i]
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

type Media struct {
   Period struct {
      Adaptation_Set []Adaptation `xml:"AdaptationSet"`
   }
}

type Content_Protection struct {
   PSSH string `xml:"pssh"`
   Scheme_ID_URI string `xml:"schemeIdUri,attr"`
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

// SegmentTemplate is sometimes under AdaptationSet and
// sometimes under Representation
type Adaptation struct {
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Lang string `xml:"lang,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Representation []Represent
   Role *struct {
      Value string `xml:"value,attr"`
   }
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
}

// SegmentTemplate is sometimes under AdaptationSet and
// sometimes under Representation
type Represent struct {
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Height int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Width int64 `xml:"width,attr"`
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   Adaptation *Adaptation // this is to get to Lang and Role
}
