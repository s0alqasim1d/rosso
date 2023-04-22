package tls

import (
   "fmt"
   "net/http"
   "net/url"
   "strings"
   "testing"
   "time"
)

func Test_UnmarshalText(t *testing.T) {
   body := url.Values{
      "Email": {email},
      "Passwd": {passwd},
      "client_sig": {""},
      "droidguard_results": {"-"},
   }.Encode()
   for _, test := range tests {
      hello, err := Parse(test.in)
      if err != nil {
         t.Fatal(err)
      }
      req, err := http.NewRequest(
         "POST", "https://android.googleapis.com/auth",
         strings.NewReader(body),
      )
      if err != nil {
         t.Fatal(err)
      }
      req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
      res, err := hello.Transport().RoundTrip(req)
      if err != nil {
         t.Fatal(err)
      }
      defer res.Body.Close()
      fmt.Println(res.Status, test.in)
      time.Sleep(time.Second)
   }
}
