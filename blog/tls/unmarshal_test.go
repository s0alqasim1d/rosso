package tls

import (
   "bufio"
   "crypto/tls"
   "fmt"
   "net"
   "net/http"
   "net/url"
   "strings"
   "testing"
)

func Test_Unmarshal(t *testing.T) {
   var hello Client_Hello_Msg
   err := hello.UnmarshalBinary(android_client_hello[5:])
   fmt.Printf("%v %#v\n", err, hello)
}

// Android_API
//                   0   1   2   3   4
var android_client_hello = []byte("\x16\x03\x01\x00\xc2\x01\x00\x00\xbe\x03\x03;c\xba\xeb\v\xe0\x86u\xacLF\x8ft]\x9f4\x86g*\xd2\\n\x17\xa8\xb8\x816^O\xd1\xe0= \x84\x92\xf5\xeaB\x86\xad\xda~\xd49ML\"-\x00=\xdd\xdaZ\xa6=\xba\xdfL`r\xa7\xd2\xe1\x1eD\x00\x1c\xc0+\xc0,̩\xc0/\xc00̨\xc0\t\xc0\n\xc0\x13\xc0\x14\x00\x9c\x00\x9d\x00/\x005\x01\x00\x00Y\xff\x01\x00\x01\x00\x00\x00\x00\x14\x00\x12\x00\x00\x0fmail.google.com\x00\x17\x00\x00\x00#\x00\x00\x00\r\x00\x06\x00\x04\x04\x03\x04\x01\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\x10\x00\v\x00\t\bhttp/1.1\x00\v\x00\x02\x01\x00\x00\n\x00\b\x00\x06\x00\x1d\x00\x17\x00\x18")

func (c *conn) Write(b []byte) (int, error) {
   if c.done {
      return c.Conn.Write(b)
   }
   c.done = true
   return c.Conn.Write(android_client_hello)
}

type conn struct {
   done bool
   net.Conn
}

func Test_Android(t *testing.T) {
   // http.Request
   body := url.Values{
      "Email": {email},
      "Passwd": {passwd},
      "client_sig": {""},
      "droidguard_results": {"-"},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth",
      strings.NewReader(body),
   )
   if err != nil {
      t.Fatal(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   dial_conn, err := net.Dial("tcp", "android.googleapis.com:443")
   if err != nil {
      t.Fatal(err)
   }
   tls_conn := tls.Client(
      &conn{Conn: dial_conn}, &tls.Config{ServerName: "android.googleapis.com"},
   )
   if err := req.Write(tls_conn); err != nil {
      t.Fatal(err)
   }
   res, err := http.ReadResponse(bufio.NewReader(tls_conn), nil)
   if err != nil {
      t.Fatal(err)
   }
   if err := tls_conn.Close(); err != nil {
      t.Fatal(err)
   }
   if err := res.Body.Close(); err != nil {
      t.Fatal(err)
   }
   if res.StatusCode != http.StatusOK {
      t.Fatal(res)
   }
}
