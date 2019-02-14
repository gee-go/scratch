```go

type Dim struct {
    Key, Value uint32
}

type Tree interface {
    // Insert tags into the tree.
    // NOTE:
    //    tags must be sorted!
    //
    //    On first insert for a full tag set, insert "fingers" for each unique key seen.
    //    These fingers are ordered by key and can be used for group by and other queries
    Insert(tags []Dim)

}

```
