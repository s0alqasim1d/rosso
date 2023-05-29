package dash

import (
   "bytes"
   "fmt"
   "os"
   "strings"
   "testing"
)

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
      reps := pre.Represents().Filter(Audio)
      target := reps.Index(func(r Represent) bool {
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
      reps := pre.Represents().Filter(Video)
      if order == random {
         reps.Sort(func(a, b Represent) bool {
            return a.Bandwidth < b.Bandwidth
         })
      }
      var target int
      if order == descending {
         target = reps.Last_Index(func(a Represent) bool {
            return a.Bandwidth >= 999_999
         })
      } else {
         target = reps.Index(func(a Represent) bool {
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
      reps := pre.Represents().Filter(func(r Represent) bool {
         return Audio(r) || Video(r)
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

