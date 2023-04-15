package http

import (
   "net/http"
   "testing"
)

func do() error {
   c := Default_Client
   c.CheckRedirect = nil
   c.Log_Level = 9
   c.Transport = new(http.Transport)
   c.Status = 201
   req := Get()
   req.URL.Scheme = "http"
   req.URL.Host = "httpbin.org"
   req.URL.Path = "/status/201"
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
