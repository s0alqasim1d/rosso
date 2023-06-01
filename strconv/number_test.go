package strconv

import (
   "fmt"
   "testing"
)

func Test_String(t *testing.T) {
   if s := fmt.Sprint(Cardinal(123)); s != "123" {
      t.Fatal(s)
   }
   if s := fmt.Sprint(Cardinal(1234)); s != "1.23 thousand" {
      t.Fatal(s)
   }
   if s := fmt.Sprint(Size(123)); s != "123 byte" {
      t.Fatal(s)
   }
   if s := fmt.Sprint(Size(1234)); s != "1.23 kilobyte" {
      t.Fatal(s)
   }
   if s := fmt.Sprint(1234/Rate(10)); s != "123 byte/s" {
      t.Fatal(s)
   }
   if s := fmt.Sprint(12345/Rate(10)); s != "1.23 kilobyte/s" {
      t.Fatal(s)
   }
   if s := fmt.Sprint(Percent(1234)/10000); s != "12.34%" {
      t.Fatal(s)
   }
}
