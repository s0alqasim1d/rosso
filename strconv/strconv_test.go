package strconv

import (
   "bytes"
   "encoding/xml"
   "fmt"
   "html"
   "mime/quotedprintable"
   "net/http"
   "net/http/httputil"
   "net/url"
   "strconv"
   "strings"
   "testing"
   "text/template"
)

func Test_Encode(t *testing.T) {
   s := "\x01¶'"
   escapes := []string{
      // github.com/golang/go/blob/go1.20.2/src/encoding/xml/xml.go#L1902
      func() string {
         var b strings.Builder
         xml.EscapeText(&b, []byte(s))
         return b.String()
      }(),
      // github.com/golang/go/blob/go1.20.2/src/html/escape.go#L178
      html.EscapeString(s),
      // github.com/golang/go/blob/go1.20.2/src/mime/quotedprintable/writer.go#L31
      // github.com/golang/go/blob/go1.20.2/src/mime/quotedprintable/writer.go#L112
      func() string {
         var b strings.Builder
         w := quotedprintable.NewWriter(&b)
         w.Write([]byte(s))
         w.Close()
         return b.String()
      }(),
      // github.com/golang/go/blob/go1.20.2/src/net/url/url.go#L281
      url.PathEscape(s),
      // github.com/golang/go/blob/go1.20.2/src/net/url/url.go#L275
      url.QueryEscape(s),
      // github.com/golang/go/blob/go1.20.2/src/strconv/quote.go#L128
      strconv.Quote(s),
      // github.com/golang/go/blob/go1.20.2/src/strconv/quote.go#L141
      strconv.QuoteToASCII(s),
      // github.com/golang/go/blob/go1.20.2/src/strconv/quote.go#L155
      strconv.QuoteToGraphic(s),
      // github.com/golang/go/blob/go1.20.2/src/text/template/funcs.go#L611
      func() string {
         var b strings.Builder
         template.HTMLEscape(&b, []byte(s))
         return b.String()
      }(),
      // github.com/golang/go/blob/go1.20.2/src/text/template/funcs.go#L671
      func() string {
         var b strings.Builder
         template.JSEscape(&b, []byte(s))
         return b.String()
      }(),
   }
   for _, escape := range escapes {
      fmt.Println(escape)
   }
}

func Test_Binary(t *testing.T) {
   src, dst, err := round_trip("https://picsum.photos/1")
   if err != nil {
      t.Fatal(err)
   }
   if !bytes.Equal(src, dst) {
      t.Fatal(dst)
   }
}

func round_trip(s string) ([]byte, []byte, error) {
   res, err := http.Get(s)
   if err != nil {
      return nil, nil, err
   }
   src, err := httputil.DumpResponse(res, true)
   if err != nil {
      return nil, nil, err
   }
   if err := res.Body.Close(); err != nil {
      return nil, nil, err
   }
   dst, err := decode(Encode(src))
   if err != nil {
      return nil, nil, err
   }
   return src, dst, nil
}

func Test_Text(t *testing.T) {
   src, dst, err := round_trip("http://httpbin.org/get")
   if err != nil {
      t.Fatal(err)
   }
   if !bytes.Equal(src, dst) {
      t.Fatal(dst)
   }
}
