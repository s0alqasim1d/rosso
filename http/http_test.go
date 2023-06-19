package http

import (
   "net/http"
   "net/url"
   "os"
   "testing"
)

func Test_URL(t *testing.T) {
   req, err := Get_Parse("http://httpbin.org/get")
   if err != nil {
      t.Fatal(err)
   }
   res, err := new(http.Transport).RoundTrip(req.Request)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   os.Stdout.ReadFrom(res.Body)
}

func Test_Body(t *testing.T) {
   req := Post(&url.URL{
      Scheme: "http",
      Host: "httpbin.org",
      Path: "/post",
   })
   req.Body_String("hello=world")
   res, err := new(http.Transport).RoundTrip(req.Request)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   os.Stdout.ReadFrom(res.Body)
}

func do() error {
   c := Default_Client
   c.CheckRedirect = nil
   c.Log_Level = 9
   c.Transport = new(http.Transport)
   c.Status = 201
   req := Get(&url.URL{
      Scheme: "http",
      Host: "httpbin.org",
      Path: "/status/201",
   })
   res, err := c.Do(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

func Test_Client(t *testing.T) {
   err := do()
   if err != nil {
      t.Fatal(err)
   }
   if Default_Client.CheckRedirect == nil {
      t.Fatal("CheckRedirect")
   }
   if Default_Client.Log_Level == 9 {
      t.Fatal("Log_Level")
   }
   if Default_Client.Status == 201 {
      t.Fatal("Status")
   }
   if Default_Client.Transport != nil {
      t.Fatal("Transport")
   }
}
