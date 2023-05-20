package url

import (
   "fmt"
   "net/url"
   "testing"
)

func Test_Path(t *testing.T) {
   ref := &url.URL{
      Scheme: "https",
      Host: "therokuchannel.roku.com",
   }
   err := set_path(ref, "/homescreen/content/https:%2F%2Fcontent.sr.roku.com")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(ref)
}
