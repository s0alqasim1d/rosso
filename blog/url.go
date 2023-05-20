package url

import "net/url"

func set_path(u *url.URL, p string) error {
   var err error
   u.Path, err = url.PathUnescape(p)
   if err != nil {
      return err
   }
   u.RawPath = p
   return nil
}
