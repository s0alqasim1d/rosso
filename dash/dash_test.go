package dash

import (
   "bytes"
   "fmt"
   "net/http"
   "os"
   "strings"
   "testing"
)

const (
   ascending = iota
   descending
   random
)

var tests = map[string]int{
   "mpd/amc.mpd": ascending,
   "mpd/paramount.mpd": random,
   "mpd/roku.mpd": descending,
}

func Test_Audio(t *testing.T) {
   for name := range tests {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      pre, err := New_Presentation(bytes.NewReader(text))
      if err != nil {
         t.Fatal(err)
      }
      reps := Filter(pre.Representation(), Audio)
      target := Index_Func(reps, func(r Representation) bool {
         if !strings.HasPrefix(r.Adaptation.Lang, "en") {
            return false
         }
         if !strings.Contains(r.Codecs, "mp4a.") {
            return false
         }
         if r.Role() == "description" {
            return false
         }
         return true
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

func Test_Ext(t *testing.T) {
   for name := range tests {
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

func Test_Info(t *testing.T) {
   for name := range tests {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      pre, err := New_Presentation(bytes.NewReader(text))
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      reps := Filter(pre.Representation(), func(r Representation) bool {
         return Audio(r) || Video(r)
      })
      for _, rep := range reps {
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

func Test_Media(t *testing.T) {
   text, err := os.ReadFile("mpd/roku.mpd")
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

func Test_Video(t *testing.T) {
   for name, order := range tests {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      pre, err := New_Presentation(bytes.NewReader(text))
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      reps := Filter(pre.Representation(), Video)
      var target int
      switch order {
      case ascending:
         target = Index_Func(reps, func(a Representation) bool {
            return a.Bandwidth >= 999_999
         })
      case descending:
         target = Last_Index_Func(reps, func(a Representation) bool {
            return a.Bandwidth >= 999_999
         })
      case random:
         Sort_Func(reps, func(a, b Representation) bool {
            return a.Bandwidth < b.Bandwidth
         })
         target = Index_Func(reps, func(a Representation) bool {
            return a.Bandwidth >= 999_999
         })
      }
      for i, rep := range reps {
         if i == target {
            fmt.Print("!")
         }
         fmt.Println(rep)
      }
      fmt.Println()
   }
}
