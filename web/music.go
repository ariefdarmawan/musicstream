package webapp

import (
	//"time"
	//"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type Music struct {
}

func (m *Music) Index(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputTemplate
	return struct{}{}
}

type searchModel struct{
	Keyword string 
	Fullsearch bool
	Song int
	Album int
	Artist int
	Take 	int
	Skip int
}

func (m *Music) Search(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	result := toolkit.NewResult()
	
	sm := new(searchModel)
	e := ctx.GetPayload(&sm)
	if e!=nil {
		return result.SetErrorTxt("Invalid parameter")
	}
	if sm.Keyword=="" {
		return result.SetErrorTxt("No keyword entered")
	}
	if sm.Fullsearch {
		sm.Keyword = "\""+sm.Keyword+"\""
	}
	
	var conn dbox.IConnection
	conn, e = connection()
	if e!=nil {
		return result.SetError(e)
	}
	defer conn.Close()
	
	mout := toolkit.M{}
	if sm.Song==1{
		songs := []toolkit.M{}
		csong, e := conn.NewQuery().
			Select("_id","song_name","artist_id","album_id","popularity").
			From("song").
			Where(dbox.Eq("$text",toolkit.M{}.Set("$search",sm.Keyword))).
			Order("-popularity").
			Skip(sm.Skip).
			Take(sm.Take).
			Cursor(nil)
		if e!=nil {
			return result.SetError(e)
		}
		defer csong.Close()
		
		csong.Fetch(&songs,0,false)
		if len(songs)>0{
			getSongs(conn, &songs)
		}
		mout.Set("song", songs)
	}
	
	if sm.Album==1{
		album := []toolkit.M{}
		calbum, e := conn.NewQuery().
			Select("_id","album_name","main_artist_id","popularity").
			From("album").
			Where(dbox.Eq("$text",toolkit.M{}.Set("$search",sm.Keyword))).
			Order("-popularity").
			Skip(sm.Skip).
			Take(sm.Take).
			Cursor(nil)
		if e!=nil {
			return result.SetError(e)
		}
		defer calbum.Close()
		
		calbum.Fetch(&album,0,false)
		if len(album)>0{
			getAlbum(conn,&album)
		} 
		mout.Set("album", album)
	}
	
	if sm.Artist==1{
		artist := []toolkit.M{}
		cartist, e := conn.NewQuery().
			Select("_id","artist_name","popularity").
			From("artist").
			Where(dbox.Eq("$text",toolkit.M{}.Set("$search",sm.Keyword))).
			Order("-popularity").
			Skip(sm.Skip).
			Take(sm.Take).
			Cursor(nil)
		if e!=nil {
			return result.SetError(e)
		}
		defer cartist.Close()
		
		cartist.Fetch(&artist,0,false)
		mout.Set("artist", artist)
	}
	
	result.SetData(mout)
	return result
}

func getSongs(conn dbox.IConnection, songs *[]toolkit.M){
	album := []toolkit.M{}
	artist := []toolkit.M{}
	for _, song := range *songs{
		calbum, _ := conn.NewQuery().SetConfig("multiexec",true).Select("_id","album_name").From("album").
			Where(dbox.Eq("_id",song.GetString("album_id"))).
			Cursor(nil)
		calbum.Fetch(&album,1,false)
		if len(album)>0{
			song["album_name"]=album[0].GetString("album_name")
		}
		calbum.Close()
		
		cartist, _ := conn.NewQuery().SetConfig("multiexec",true).Select("_id","artist_name").From("artist").
			Where(dbox.Eq("_id",song.GetString("artist_id"))).
			Cursor(nil)
		cartist.Fetch(&artist,1,false)
		if len(artist)>0{
			song["artist_name"]=artist[0].GetString("artist_name")
		}
		cartist.Close()
	}
}

func getAlbum(conn dbox.IConnection, songs *[]toolkit.M){
	artist := []toolkit.M{}
	for _, song := range *songs{
		cartist, _ := conn.NewQuery().SetConfig("multiexec",true).Select("_id","artist_name").From("artist").
			Where(dbox.Eq("_id",song.GetString("main_artist_id"))).
			Cursor(nil)
		cartist.Fetch(&artist,1,false)
		if len(artist)>0{
			song["artist_name"]=artist[0].GetString("artist_name")
		}
		cartist.Close()
	}
}
