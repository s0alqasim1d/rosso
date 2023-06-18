package dash

import (
   "2a.pages.dev/mech"
   "2a.pages.dev/rosso/dash"
)

func (d decrypt) media(
   client_ID []byte
   n mech.Namer
   p widevine.Poster
   private_key []byte
) error {
   name, err := mech.Name(n)
   if err != nil {
      return err
   }
   file, err := os.Create(name + d.r.Ext())
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := file.ReadFrom(d.b); err != nil {
      return err
   }
   media := d.r.Segment_Template.Get_Media()
   pro := http.Progress_Chunks(file, len(media))
   pssh, err := base64.StdEncoding.DecodeString(d.r.Widevine())
   if err != nil {
      return err
   }
   mod, err := widevine.New_Module(private_key, client_ID, pssh)
   if err != nil {
      return err
   }
   keys, err := mod.Post(p)
   if err != nil {
      return err
   }
   key := keys.Content().Key
   for _, ref := range media {
      req, err := http.NewRequest("GET", ref, nil)
      if err != nil {
         return err
      }
      req.URL = d.u.ResolveReference(req.URL)
      res, err := new(http.Client).Do(req)
      if err != nil {
         return err
      }
      pro.Add_Chunk(res.ContentLength)
      if err := d.d.Segment(res.Body, pro, keys.Content().Key); err != nil {
         return err
      }
      if err := res.Body.Close(); err != nil {
         return err
      }
   }
   return nil
}
