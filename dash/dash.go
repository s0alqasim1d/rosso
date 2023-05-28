package dash

import (
   "strconv"
   "strings"
)

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
   b = append(b, "\n\tMIME type: "...)
   b = append(b, r.MIME_Type...)
   if r.Adaptation.Role != nil {
      b = append(b, "\n\trole: "...)
      b = append(b, r.Adaptation.Role.Value...)
   }
   if r.Adaptation.Lang != "" {
      b = append(b, "\n\tlang: "...)
      b = append(b, r.Adaptation.Lang...)
   }
   return string(b)
}

type Adaptation struct {
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Lang string `xml:"lang,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Role *struct {
      Value string `xml:"value,attr"`
   }
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   Representation []Representation
}

func (r Representation) Widevine() *Content_Protection {
   for _, c := range r.Content_Protection {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return &c
      }
   }
   return nil
}

type Content_Protection struct {
   PSSH string `xml:"pssh"`
   Scheme_ID_URI string `xml:"schemeIdUri,attr"`
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

func (r Representation) Initialization() string {
   return r.replace_ID(r.Segment_Template.Initialization)
}

type Segment_Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   Segment_Timeline struct {
      S []Segment
   } `xml:"SegmentTimeline"`
   Start_Number *int `xml:"startNumber,attr"`
}

func (r Representation) Ext() string {
   switch r.MIME_Type {
   case "video/mp4":
      return ".m4v"
   case "audio/mp4":
      return ".m4a"
   }
   return ""
}

func (r Representation) Role() string {
   if r.Adaptation.Role == nil {
      return ""
   }
   return r.Adaptation.Role.Value
}

type Presentation struct {
   Period struct {
      Adaptation_Set []Adaptation `xml:"AdaptationSet"`
   }
}

type Segment struct {
   D int `xml:"d,attr"` // duration
   R int `xml:"r,attr"` // repeat
   T int `xml:"t,attr"` // time
}

func (r Representation) replace_ID(ref string) string {
   return strings.Replace(ref, "$RepresentationID$", r.ID, 1)
}

func (s Segment) replace(ref, old string) string {
   return strings.Replace(ref, old, strconv.Itoa(s.T), 1)
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
