package dash

type Adaptation struct {
   Codecs string `xml:"codecs,attr"`
   Content_Protection []Content_Protection `xml:"ContentProtection"`
   Lang string `xml:"lang,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Representation []Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
   Segment_Template *Segment_Template `xml:"SegmentTemplate"`
}

type Content_Protection struct {
   PSSH string `xml:"pssh"`
   Scheme_ID_URI string `xml:"schemeIdUri,attr"`
}

type Presentation struct {
   Period struct {
      Adaptation_Set []Adaptation `xml:"AdaptationSet"`
   }
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
