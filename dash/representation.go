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
      b = append(b, "\nwidth: "...)
      b = strconv.AppendInt(b, r.Width, 10)
   }
   if r.Height >= 1 {
      b = append(b, "\nheight: "...)
      b = strconv.AppendInt(b, r.Height, 10)
   }
   if r.Bandwidth >= 1 {
      b = append(b, "\nbandwidth: "...)
      b = strconv.AppendInt(b, r.Bandwidth, 10)
   }
   if r.Codecs != "" {
      b = append(b, "\ncodecs: "...)
      b = append(b, r.Codecs...)
   }
   b = append(b, "\ntype: "...)
   b = append(b, r.MIME_Type...)
   if r.Adaptation.Role != nil {
      b = append(b, "\nrole: "...)
      b = append(b, r.Adaptation.Role.Value...)
   }
   if r.Adaptation.Lang != "" {
      b = append(b, "\nlanguage: "...)
      b = append(b, r.Adaptation.Lang...)
   }
   return string(b)
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
   Adaptation *Adaptation // this is to get to Role
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

func (r Representation) Ext() string {
   switch {
   case Audio(r):
      return ".m4a"
   case Video(r):
      return ".m4v"
   }
   return ""
}

func (r Representation) Initialization() string {
   return r.replace_ID(r.Segment_Template.Initialization)
}

func (r Representation) Role() string {
   if r.Adaptation.Role != nil {
      return r.Adaptation.Role.Value
   }
   return ""
}

func (r Representation) Widevine() *Content_Protection {
   for _, c := range r.Content_Protection {
      if c.Scheme_ID_URI == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return &c
      }
   }
   return nil
}

func (r Representation) replace_ID(ref string) string {
   return strings.Replace(ref, "$RepresentationID$", r.ID, 1)
}
