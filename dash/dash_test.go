package dash

import (
   "2a.pages.dev/rosso/slices"
   "bytes"
   "fmt"
   "net/http"
   "os"
   "strings"
   "testing"
)

func Test_Media(t *testing.T) {
   text, err := os.ReadFile("mpd/roku.mpd")
   if err != nil {
      t.Fatal(err)
   }
   med, err := New_Media(bytes.NewReader(text))
   if err != nil {
      t.Fatal(err)
   }
   base, err := http.NewRequest("", "http://example.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   for _, ref := range med.Period.Adaptation_Set[0].Representation[0].Media() {
      req, err := http.NewRequest("", ref, nil)
      if err != nil {
         t.Fatal(err)
      }
      req.URL = base.URL.ResolveReference(req.URL)
      fmt.Println(req.URL)
   }
}
var tests = []string{
   "mpd/amc.mpd",
   "mpd/paramount.mpd",
   "mpd/roku.mpd",
}

func Test_Info(t *testing.T) {
   for _, name := range tests {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      med, err := New_Media(bytes.NewReader(text))
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      reps := slices.Delete(med.Represents(), func(r Represent) bool {
         if Audio(r) {
            return false
         }
         if Video(r) {
            return false
         }
         return true
      })
      for i, rep := range reps {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

func Test_Ext(t *testing.T) {
   for _, name := range tests {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      med, err := New_Media(bytes.NewReader(text))
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      for _, rep := range med.Represents() {
         fmt.Printf("%q\n", rep.Ext())
      }
      fmt.Println()
   }
}

func Test_Video(t *testing.T) {
   for _, name := range tests {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      med, err := New_Media(bytes.NewReader(text))
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      reps := slices.Delete(med.Represents(), Not(Video))
      slices.Sort(reps, func(a, b Represent) int {
         return int(b.Bandwidth - a.Bandwidth)
      })
      target := slices.Index(reps, func(a Represent) bool {
         return a.Bandwidth <= 9_000_000
      })
      for i, rep := range reps {
         if i >= 1 {
            fmt.Println()
         }
         if i == target {
            fmt.Print("!")
         }
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

func Test_Audio(t *testing.T) {
   for _, name := range tests {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      med, err := New_Media(bytes.NewReader(text))
      if err != nil {
         t.Fatal(err)
      }
      reps := slices.Delete(med.Represents(), Not(Audio))
      target := slices.Index(reps, func(r Represent) bool {
         if strings.HasPrefix(r.Adaptation.Lang, "en") {
            return strings.Contains(r.Codecs, "mp4a.")
         }
         return false
      })
      fmt.Println(name)
      for i, rep := range reps {
         if i == target {
            fmt.Print("!")
         }
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

