package dash

import (
   "2a.pages.dev/rosso/slices"
   "bytes"
   "fmt"
   "os"
   "strings"
   "testing"
)

func Test_Video(t *testing.T) {
   for name, sorted := range tests {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      pre, err := New_Presentation(bytes.NewReader(text))
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      reps := slices.Delete(pre.Represents(), Not(Video))
      if !sorted {
         slices.Sort(reps, func(a, b Represent) int {
            return int(b.Bandwidth - a.Bandwidth)
         })
      }
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

// true if sorted
var tests = map[string]bool{
   "mpd/amc.mpd": false,
   "mpd/paramount.mpd": false,
   "mpd/roku.mpd": true,
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
      reps := slices.Delete(pre.Represents(), func(r Represent) bool {
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
      reps := slices.Delete(pre.Represents(), Not(Audio))
      target := slices.Index(reps, func(r Represent) bool {
         if strings.HasPrefix(r.Adaptation.Lang, "en") {
            if strings.Contains(r.Codecs, "mp4a.") {
               return r.Role() != "description"
            }
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
