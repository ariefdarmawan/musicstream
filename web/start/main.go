package main

import (
    "github.com/eaciit/knot/knot.v1"
    "eaciit/musicstream/web"
)

func main() {
    app := webapp.App()
    knot.StartApp(app, "localhost:8027")
}