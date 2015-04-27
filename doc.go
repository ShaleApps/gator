/*
package gator (or valigator) is a library that validates structs using struct tags.  Here is a usage example:

    // NOTE: Gator tags will only be recognized on fields that are public.
    type BigStruct struct {
        Required        string      `gator:”nonzero”`
        Email           string      `gator:”email”`
        Website         string      `gator:"url"`
        IPAddress       string      `gator:"ip"`
        PetName         string      `gator:"alpha"`
        Phone           string      `gator:"num | len(10)"`
        Password        string      `gator:”alphanum | minlen(5) | maxlen(15)"`
        DayOfWeek       int         `gator:”gte(0) | lt(7)”`
        Lat             float64     `gator:”lat”`
        Lng             float64     `gator:”lon”`
        TennisScore string          `gator:”in(love,15,30,40)”`
        NewSuperHero    string      `gator:”notin(Superman,Batman,The Flash)”`
        Zipcode         string      `gator:”match(^\d{5}(?:[-\s]\d{4})?$)”`
        Color           string      `gator:”hexcolor”`
        Ages            []int       `gator:”each( gt(18) | lt(35) )”`
    }
    b := &BigStruct{
        Password: "TOOOOOOOOOOOOOOO LONG",
    }
    if err := gator.NewStruct(b).Validate(); err != nil {
        fmt.Println(err)
    }

Validation logic can be deserialized by gator using a query string:

    type WebsiteListing struct {
        Url      string
        Username string
    }
    website := &WebsiteListing{
        Url:      "https//news.ycombinator.com",
        Username: "hello1",
    },
    g := gator.NewQueryStr(website, "Url=url&Username=alphanum|minlen(5)|maxlen(10)")
    if err := g.Validate(); err != nil {
        fmt.Println(err)
    }

Custom tags can be added.  Tokens are added statically and affect all Gators.  Here is an example:

    gator.RegisterStructTagToken("pword", func(s string) gator.Func {
        return gator.Matches(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{4,8}$`)
    })
    type User struct {
        Email    string `gator:”email”`
        Password string `gator:"pword"`
    }
    u := &User{
        Email:    "gator@example.com",
        Password: "ASDF12345",
    }
    if err := gator.NewStruct(u).Validate(); err != nil {
        fmt.Println(err)
    }

*/
package gator
