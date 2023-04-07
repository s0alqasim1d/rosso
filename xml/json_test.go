package xml

import (
   "os"
   "strings"
   "testing"
)

const src = `
<regionalRating>
   <rating>TV-PG</rating>
   <region>CA</region>
</regionalRating>
`

func Test_Indent(t *testing.T) {
   err := Indent(os.Stdout, strings.NewReader(src), "", " ")
   if err != nil {
      t.Fatal(err)
   }
}
