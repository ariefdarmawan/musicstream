package main

import (
    "github.com/ariefdarmawan/flat"
    "github.com/eaciit/toolkit"
    "github.com/eaciit/config"
    "os"
    "io"
    "strings"
    "github.com/eaciit/dbox"
    "path/filepath"
    _ "github.com/eaciit/dbox/dbc/mongo"
)

var (
    path string = "/Users/ariefdarmawan/Dropbox/biz/eaciit/Melon/05. From Clients/DATA_POC_EACIIT/"   
) 

func main(){    
    readConfig()
    import2db("album", "album_id")
    import2db("artist","artist_id")
    import2db("song","song_id")
    //import2db("mu_user_stream","")
    //import2db("mu_user_down","")
}

func checkX(e error, pre string){
    if e!=nil {
        toolkit.Println(pre," Error: ",e)
        os.Exit(400)
    }
}

func readConfig(){
    wd, e := os.Getwd()
    checkX(e,"Working Directory")
    wd = filepath.Join(wd,"../../config/app.json")
    
    checkX(config.SetConfigFile(wd), "Config")
    toolkit.Println("Applying config:",wd)
}

func connection() (c dbox.IConnection, e error){
    ci := &dbox.ConnectionInfo{
        config.GetDefault("db_host","").(string),
        config.GetDefault("db_name","").(string),
        config.GetDefault("db_user","").(string),
        config.GetDefault("db_password","").(string),
        toolkit.M{}}
    c, _ = dbox.NewConnection("mongo",ci)
    e = c.Connect()
    return
}

func import2db(name string, idFieldName string){
    filename := path + name + ".txt"
    toolkit.Println("Reading file", filename)
    f := flat.New(filename, true, false)
    f.Delimeter = '|'
    checkX(f.Open(),"Import")
    defer f.Close()
    
    c, e := connection()
    checkX(e,"Connection")
    defer c.Close()
    
    c.NewQuery().From(name).Delete().Exec(nil)
    
    q := c.NewQuery().SetConfig("multiexec",true).From(name).Save()
    defer q.Close()
     
    isEOF := false
    i := 0
    for ;!isEOF;{
        i++
        m, e := f.ReadM()
        if e==io.EOF{
            isEOF=true   
        } else if e!=nil {
            checkX(e, "Read")
        } else {
            var id interface{}
            if idFieldName=="" {
               m.Set("_id",i)
               id = i
            } else {
                id = m.Get(idFieldName,"")
                m.Set("_id",id)
            }
            if m.Has("search_keyword"){
                keywords:=strings.Split(m.GetString("search_keyword"),",")
                newkeywords := []string{}
                for _, keyword := range keywords{
                    keyword = strings.ToLower(strings.Trim(keyword," "))
                    if keyword!=""{
                        newkeywords = append(newkeywords, keyword)
                    }
                }
                m.Set("popularity",0)
                m.Set("search_keyword",newkeywords)
            }
            toolkit.Printf("Saving %s record no: %d ID: %v \n", name, i, id)
            checkX(q.Exec(toolkit.M{}.Set("data",m)), toolkit.Sprintf("Saving %s data: %v",id, m))
        }
    }
}