package dash

import (
   "2a.pages.dev/mech"
   "2a.pages.dev/rosso/dash"
)

type decrypt struct {
   b *bytes.Buffer
   d mp4.Decrypt
   r dash.Representer
   u *url.URL
}

func initialization(r dash.Representer, u *url.URL) (*decrypt, error) {
   var dec decrypt
   dec.b = new(bytes.Buffer)
   dec.d = make(mp4.Decrypt)
   dec.r = r
   dec.u = u
   req, err := http.NewRequest(
      "GET", r.Segment_Template.Get_Initialization(), nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL = dec.u.ResolveReference(req.URL)
   res, err := new(http.Client).Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if err := dec.d.Init(res.Body, dec.b); err != nil {
      return nil, err
   }
   return &dec, nil
}
