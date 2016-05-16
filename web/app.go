package webapp

import (
    _ "github.com/eaciit/dbox/dbc/mongo"
    "github.com/eaciit/dbox"
    //"github.com/eaciit/orm"
    "github.com/eaciit/knot/knot.v1"
    "github.com/eaciit/config"
    "os"
    "path/filepath"
    "github.com/eaciit/toolkit"
)

func App() *knot.App{
    readConfig()
    app := knot.NewApp("music")
    wd, _ := os.Getwd()
    wd += "/../"
    app.ViewsPath = wd+"views/"
    app.LayoutTemplate="_layout.html"
    app.Static("static",wd+"assets")
    app.Register(&Music{})
    return app
}

func readConfig() {
	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "../../config/app.json")
	config.SetConfigFile(wd)
	toolkit.Println("Applying config:", wd)
}

func connection() (c dbox.IConnection, e error) {
    //toolkit.Println("Connecting: ", config.GetDefault("db_host", "").(string))
	ci := &dbox.ConnectionInfo{
		config.GetDefault("db_host", "").(string),
		config.GetDefault("db_name", "").(string),
		config.GetDefault("db_user", "").(string),
		config.GetDefault("db_password", "").(string),
		toolkit.M{}}
	c, _ = dbox.NewConnection("mongo", ci)
	e = c.Connect()
	return
}