package dash

import (
   "bytes"
   "fmt"
   "os"
   "strings"
   "testing"
)

var tests = map[string]bool{
   "mpd/amc.mpd": false, // ascending
   "mpd/paramount-lang.mpd": true, // random
   "mpd/paramount-role.mpd": false, // random
   "mpd/roku-eng.mpd": false, // descending
}

func Test_Video(t *testing.T) {
   for name, last := range tests {
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
      if last {
         target = Last_Index_Func(reps, func(a Representation) bool {
            return a.Bandwidth >= 999_999
         })
      } else {
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
