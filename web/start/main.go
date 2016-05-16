package main

import (
    "github.com/eaciit/knot/knot.v1"
    "eaciit/musicstream/web"
    "net/http"
)

func main() {
    app := webapp.App()
    knot.StartAppWithFn(app, "localhost:8027",
        map[string]knot.FnContent{"/":func(ctx *knot.WebContext)interface{}{
            http.Redirect(ctx.Writer, ctx.Request, "/music/index", 301)
            return nil
        }})
}