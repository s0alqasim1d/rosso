package dash

import (
   "bytes"
   "fmt"
   "net/http"
   "os"
   "testing"
)

var tests = []string{
   "mpd/amc.mpd",
   "mpd/paramount-lang.mpd",
   "mpd/paramount-role.mpd",
   "mpd/roku-eng.mpd",
}

func Test_Ext(t *testing.T) {
   for _, name := range tests {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      pre, err := New_Presentation(bytes.NewReader(text))
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      for _, rep := range pre.Representation() {
         fmt.Printf("%q\n", rep.Ext())
      }
      fmt.Println()
   }
}

func Test_Media(t *testing.T) {
   text, err := os.ReadFile("mpd/roku-eng.mpd")
   if err != nil {
      t.Fatal(err)
   }
   pre, err := New_Presentation(bytes.NewReader(text))
   if err != nil {
      t.Fatal(err)
   }
   base, err := http.NewRequest("", "http://example.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   for _, ref := range pre.Period.Adaptation_Set[0].Representation[0].Media() {
      req, err := http.NewRequest("", ref, nil)
      if err != nil {
         t.Fatal(err)
      }
      req.URL = base.URL.ResolveReference(req.URL)
      fmt.Println(req.URL)
   }
}
