package dash

import (
   "encoding/xml"
   "io"
   "strconv"
   "strings"
)

func (r Representation) String() string {
   var s []string
   if r.Width >= 1 {
      s = append(s, "width: " + strconv.Itoa(r.Width))
   }
   if r.Height >= 1 {
      s = append(s, "height: " + strconv.Itoa(r.Height))
   }
   if r.Bandwidth >= 1 {
      s = append(s, "bandwidth: " + strconv.Itoa(r.Bandwidth))
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
         rep := rep
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
         if rep.Segment_Template != nil {
            rep.Segment_Template.Representation = &rep
         }
         reps = append(reps, rep)
      }
   }
   return reps, nil
}

func (s Segment_Template) Get_Media() []string {
   var refs []string
   for _, seg := range s.Segment_Timeline.S {
      seg.T = s.Start_Number
      for seg.R >= 0 {
         {
            ref := s.Media
            replace(&ref, "$Number$", strconv.Itoa(seg.T))
            replace(&ref, "$RepresentationID$", s.Representation.ID)
            refs = append(refs, ref)
         }
         s.Start_Number++
         seg.R--
         seg.T++
      }
   }
   return refs
}

func Audio(r Representation) bool {
   return *r.MIME_Type == "audio/mp4"
}

func Not[E any](fn func(E) bool) func(E) bool {
   return func(value E) bool {
      return !fn(value)
   }
}

func Video(r Representation) bool {
   return *r.MIME_Type == "video/mp4"
}

func replace(s *string, in, out string) {
   *s = strings.Replace(*s, in, out, 1)
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

func (r Representation) Widevine() string {
   for _, c := range r.Content_Protection {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return c.PSSH
      }
   }
   return ""
}

func (s Segment_Template) Get_Initialization() string {
   replace(&s.Initialization, "$RepresentationID$", s.Representation.ID)
   return s.Initialization
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

type Representation struct {
   // roku.com
   Bandwidth int `xml:"bandwidth,attr"`
   // roku.com
   Codecs string `xml:"codecs,attr"`
   // roku.com
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   // roku.com
   Height int `xml:"height,attr"`
   // roku.com
   ID string `xml:"id,attr"`
   // paramountplus.com
   MIME_Type *string `xml:"mimeType,attr"`
   // roku.com
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   // roku.com
   Width int `xml:"width,attr"`
   // roku.com
   Adaptation_Set *Adaptation_Set
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
   Representation *Representation 
   Segment_Timeline struct {
      S []Segment
   } `xml:"SegmentTimeline"`
   Start_Number int `xml:"startNumber,attr"`
}
