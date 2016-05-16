package main

import (
	"os"

	"github.com/eaciit/config"
	"github.com/eaciit/toolkit"
	//"io"
	//"strings"
	"path/filepath"
	//"sync"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
)

var (
	path string = "/Users/ariefdarmawan/Dropbox/biz/eaciit/Melon/05. From Clients/DATA_POC_EACIIT/"
)

func main() {
	readConfig()
	//buildTag()
    buildAlbumAndArtist()
}

func checkX(e error, pre string) {
	if e != nil {
		toolkit.Println(pre, " Error: ", e)
		os.Exit(400)
	}
}

func readConfig() {
	wd, e := os.Getwd()
	checkX(e, "Working Directory")
	wd = filepath.Join(wd, "../../config/app.json")

	checkX(config.SetConfigFile(wd), "Config")
	toolkit.Println("Applying config:", wd)
}

func connection() (c dbox.IConnection, e error) {
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

func buildTag() {
	c, e := connection()
	checkX(e, "Connection")
	defer c.Close()

	//c.NewQuery().From("search_song").Delete().Exec(nil)
	//c.NewQuery().From("search_album").Delete().Exec(nil)
	//c.NewQuery().From("search_artist").Delete().Exec(nil)

	//wg := new(sync.WaitGroup)
    
    //wg.Add(2)
	c1, _ := c.NewQuery().From("popular").Cursor(nil)
	defer c1.Close()
    buildTagByCursor(c, c1, 1)
	
    c2, _ := c.NewQuery().From("populard").Cursor(nil)
	defer c2.Close()
    buildTagByCursor(c, c2, 2)
	//wg.Wait()
    toolkit.Println("Done")
}

func buildTagByCursor(c dbox.IConnection, cs dbox.ICursor, idx int) {
	//toolkit.Println(c.Info().Database)
    //defer wg.Done()
    iseof := false
    
    qsongsave := c.NewQuery().SetConfig("multiexec",true).From("song").Save()
    qalbumsave := c.NewQuery().SetConfig("multiexec",true).From("album").Save()
    qartistsave := c.NewQuery().SetConfig("multiexec",true).From("artist").Save()
    defer func(){
        qsongsave.Close()
        qalbumsave.Close()
        qartistsave.Close()
    }()
    
	for ;!iseof; {
		m := []toolkit.M{}
		checkX(cs.Fetch(&m, 1, false), "")
        if len(m) > 0 {
            songid := m[0].Get("_id", "").(string)
			//toolkit.Println("Evaluating song",songid)
			if songid != "" {
                msong := []toolkit.M{}
                csong, esong := c.NewQuery().From("song").Where(dbox.Eq("_id", songid)).Cursor(nil)
                if esong == nil {
                    csong.Fetch(&msong, 1, false)
                    if len(msong) > 0 {
                        if idx==1{
                            msong[0].Set("popularity",m[0].GetInt("c"))
                        } else {
                            msong[0].Set("popularity",msong[0].GetInt("popularity")+m[0].GetInt("c"))
                        }
                        qsongsave.Exec(toolkit.M{}.Set("data",msong[0]))
                        toolkit.Println(msong[0].Get("song_name"),":",msong[0].Get("popularity"))
                    } else {
                        toolkit.Println(songid, "could not be found")
                    }
                    csong.Close()
                } 
            } else {
                toolkit.Println(songid, "could not be found")
            }
		} else {
            iseof=true
        }
	}
}

func buildAlbumAndArtist() {
	c, e := connection()
	checkX(e, "Connection")
	defer c.Close()

	cs, _ := c.NewQuery().From("song").Where(dbox.Gt("popularity",0)).Cursor(nil)
	defer cs.Close()
	
	qalbumsave := c.NewQuery().SetConfig("multiexec",true).From("album").Save()
    qartistsave := c.NewQuery().SetConfig("multiexec",true).From("artist").Save()
    defer func(){
        qalbumsave.Close()
        qartistsave.Close()
    }()
	
	noteof := true
	for ;noteof;{
		songs := []toolkit.M{}
		cs.Fetch(&songs,1,false)
		if len(songs)>0{
			qalbum, _ := c.NewQuery().From("album").Where(dbox.Eq("_id", songs[0].GetString("album_id"))).Cursor(nil)
			albums := []toolkit.M{}
			qalbum.Fetch(&albums,1,false)
			if len(albums)>0{
				albums[0].Set("popularity",albums[0].GetInt("popularity")+songs[0].GetInt("popularity"))
				qalbumsave.Exec(toolkit.M{}.Set("data",albums[0]))
				toolkit.Println(albums[0].Get("album_name"),":",albums[0].Get("popularity"))
			}	
			qalbum.Close()		
			
			qartist, _ := c.NewQuery().From("artist").Where(dbox.Eq("_id", songs[0].GetString("artist_id"))).Cursor(nil)
			artists := []toolkit.M{}
			qartist.Fetch(&artists,1,false)
			if len(artists)>0{
				artists[0].Set("popularity",artists[0].GetInt("popularity")+songs[0].GetInt("popularity"))
				qartistsave.Exec(toolkit.M{}.Set("data",artists[0]))
				toolkit.Println(artists[0].Get("artist_name"),":",artists[0].Get("popularity"))
			}	
			qartist.Close()	
		} else {
			noteof=false
		}
	}
    
	//wg.Wait()
    toolkit.Println("Done")
}
