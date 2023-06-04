// these values can be under AdaptationSet or Representation:
// - ContentProtection
// - SegmentTemplate
// - mimeType
package dash

import (
   "encoding/xml"
   "fmt"
   "io"
   "strings"
)

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
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Height int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Width int64 `xml:"width,attr"`
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
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

func (r Representation) Text(a Adaptation_Set) string {
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
   if a.Role != nil {
      s = append(s, "role: " + a.Role.Value)
   }
   if a.Lang != "" {
      s = append(s, "language: " + a.Lang)
   }
   return strings.Join(s, "\n")
}

func Adaptation_Sets(r io.Reader) ([]Adaptation_Set, error) {
   var s struct {
      Period struct {
         Adaptation_Set []Adaptation_Set `xml:"AdaptationSet"`
      }
   }
   err := xml.NewDecoder(r).Decode(&s)
   if err != nil {
      return nil, err
   }
   return s.Period.Adaptation_Set, nil
}

func (r Representation) replace_ID(ref string) string {
   return strings.Replace(ref, "$RepresentationID$", r.ID, 1)
}

func (r Representation) Initialization() string {
   return r.replace_ID(r.Segment_Template.Initialization)
}

func (r Representation) Media() []string {
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

func (r Representation) Widevine() *Content_Protection {
   for _, c := range r.Content_Protection {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return &c
      }
   }
   return nil
}

func Audio(r Representation) bool {
   return r.MIME_Type == "audio/mp4"
}

func Video(r Representation) bool {
   return r.MIME_Type == "video/mp4"
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

