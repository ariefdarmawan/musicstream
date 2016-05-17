package main

import (
    "github.com/eaciit/knot/knot.v1"
    "eaciit/musicstream/web"
    "net/http"
    "github.com/eaciit/config"
    "github.com/eaciit/toolkit"
    "path/filepath"
    "os"
)

func readConfig() {
	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "../../config/app.json")
	config.SetConfigFile(wd)
    toolkit.Println("Applying config:",wd)
}

func main() {
    readConfig()
    app := webapp.App()
    port := toolkit.ToInt(config.Get("webapp_port").(float64), toolkit.RoundingAuto)
    knot.StartAppWithFn(app, toolkit.Sprintf("localhost:%d",port),
        map[string]knot.FnContent{"/":func(ctx *knot.WebContext)interface{}{
            http.Redirect(ctx.Writer, ctx.Request, "/music/index", 301)
            return nil
        }})
}