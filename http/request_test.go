package http

import (
   "net/url"
   "os"
   "testing"
)

func Test_Request(t *testing.T) {
   req := Get(&url.URL{
      Scheme: "http",
      Host: "httpbin.org",
      Path: "/get",
   })
   Default_Client.Log_Level = 2
   res, err := Default_Client.Do(req)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   os.Stdout.ReadFrom(res.Body)
}
