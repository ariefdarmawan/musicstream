package main

import (
    "github.com/ariefdarmawan/flat"
    "github.com/eaciit/toolkit"
    "os"
    "io"
    "github.com/eaciit/dbox"
    _ "github.com/eaciit/dbox/dbc/mongo"
)

func main(){    
    import2db("album", "ambul_id")
    import2db("artist","artist_id")
    import2db("song","song_id")
    import2db("mu_user_stream","mbp_id")
    import2db("mu_user_down","mbp_id")
}

func checkX(e error, pre string){
    if e!=nil {
        toolkit.Println(pre," Error: ",e)
        os.Exit(400)
    }
}

func connection() dbox.IConnection{
    return nil
}

func import2db(name string, idFieldName string){
    f := flat.New(name, true, false)
    checkX(f.Open(),"Import")
    defer f.Close()
    
    c := connection()
    defer c.Close()
    q := c.NewQuery().From(name).Save()
    defer q.Close()
     
    isEOF := false
    for ;!isEOF;{
        m, e := f.ReadM()
        if e==io.EOF{
            isEOF=true   
        } else if e!=nil {
            checkX(e, "Read")
        } else {
            id := m.Get(idFieldName)
            m.Set("_id",id)
            checkX(q.Exec(m), toolkit.Sprintf("Saving %s",id))
        }
    }
}