package tls

import (
   "fmt"
   "net/http"
   "net/url"
   "os"
   "strings"
   "testing"
   "time"
)

func sign_in(name string) ([]string, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   return strings.Split(string(data), "\n"), nil
}

func Test_UnmarshalText(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   account, err := sign_in(home + "/Documents/gmail.txt")
   if err != nil {
      t.Fatal(err)
   }
   body := url.Values{
      "Email": {account[0]},
      "Passwd": {account[1]},
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
