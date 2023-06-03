package dash

import (
   "fmt"
   "strings"
)

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

func (r Represent) Initialization() string {
   return r.replace_ID(r.Segment_Template.Initialization)
}

func (r Represent) Role() string {
   if r.Adaptation.Role != nil {
      return r.Adaptation.Role.Value
   }
   return ""
}

func (r Represent) Widevine() *Content_Protection {
   for _, c := range r.Content_Protection {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return &c
      }
   }
   return nil
}

func (r Represent) replace_ID(ref string) string {
   return strings.Replace(ref, "$RepresentationID$", r.ID, 1)
}

type Represent struct {
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Height int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Width int64 `xml:"width,attr"`
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
   Adaptation *Adaptation // this is to get to Role
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
