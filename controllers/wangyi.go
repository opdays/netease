package controllers

import (
	"net/http"
	"strings"
	"io/ioutil"
	"strconv"
	"crypto/md5"
	"encoding/base64"
	json "github.com/bitly/go-simplejson"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type song struct {
	Id       int `json:"id"`
	Name     string `json:"song_name"`
	Url      string `json:"song_url"`
	SongId   int `json:"song_id"`
	Album    string `json:"album"`
	AlbumId  int `json:"album_id"`
	AlbumPic string `json:"song_pic"`
	Artists  string `json:"song_artists"`
}
//type songs struct {
//	Title string
//	Name  string
//	Id    int
//	I     int
//	H     bool
//}
//type list struct {
//	Title       string
//	Name        string
//	Description string
//	List        []songs
//}
type listsong []song

var High,_  = beego.AppConfig.Bool("high")
var Ip = beego.AppConfig.String("ip")

type SongSearchControler struct {
	beego.Controller
}
type SongDetailControler struct {
	beego.Controller
}
type SongCdnDetailControler struct {
	beego.Controller
}
type PlaylistControler struct {
	beego.Controller
}
type ArtistListControler struct {
	beego.Controller
}
func init()  {
	logs.SetLogger("console")
	logs.Info("进入controler init")
	logs.Debug("进入controler init")
	logs.Error("进入controler init")
}
//关键字搜索接口
func (c *SongSearchControler) Get() {

	key := c.GetString("key")
	list := SongSearch(key)
	c.Data["json"] = list
	logs.Debug(key)
	c.ServeJSON()
}

//查询歌曲id 返回url 一系列信息
func (c *SongDetailControler) Get() {
	id := c.GetString("id")
	_, s := songDetail(id, High)

	c.Data["json"] = s
	c.ServeJSON()

}
func (c *SongCdnDetailControler) Get() {
	id := c.GetString("id")
	s := GetSongDetail(id)

	c.Data["json"] = s
	c.ServeJSON()

}


func (c *PlaylistControler) Get(){
	id := c.GetString("id")
	list := listDetail(id)
	c.Data["json"] = list
	c.ServeJSON()
}
func (c *ArtistListControler) Get(){
	id := c.GetString("id")
	list := artistList(id)
	c.Data["json"] = list
	c.ServeJSON()
}




func HttpGet(url string) (body []byte) {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	reqest.Header.Set("Cookie", "appver=1.5.0.75771")
	reqest.Header.Set("Referer", "http://music.163.com/")
	reqest.Header.Set("X-Real-Ip", Ip)
	response, err := client.Do(reqest)
	defer response.Body.Close()

	if err != nil {
		return
	}

	if response.StatusCode == 200 {
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}
	}
	return
}
func HttpPost(url, data string) (body []byte) {
	client := &http.Client{}
	reqest, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return
	}
	reqest.Header.Set("Cookie", "appver=1.5.0.75771")
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Set("Referer", "http://music.163.com/")
	reqest.Header.Set("X-Real-Ip", Ip)
	response, err := client.Do(reqest)
	defer response.Body.Close()

	if err != nil {
		return
	}

	if response.StatusCode == 200 {
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}
	}
	return
}
func songDetail(songId string, h bool) (b bool, song song) {
	s := HttpGet("http://music.163.com/api/song/detail/?id=" + songId + "&ids=%5B" + songId + "%5D")
	j, err := json.NewJson(s)
	if err != nil {
		beego.Error("NewJson")
		return
	}
	code, err := j.Get("code").Int()
	if err != nil {
		beego.Error("CodeToInt")
		return
	}
	if code != 200 {
		b = false
		return
	}
	song.Name, err = j.Get("songs").GetIndex(0).Get("name").String()
	if err != nil {
		b = false
		return
	}
	song.AlbumId, err = j.Get("songs").GetIndex(0).Get("album").Get("id").Int()
	if !h {
		song.Url, err = j.Get("songs").GetIndex(0).Get("mp3Url").String()
		if song.Url == "" {
			song.Url = albumGetUrl(strconv.Itoa(song.AlbumId), songId, h)
		}
	} else {
		dfsId, _ := j.Get("songs").GetIndex(0).Get("hMusic").Get("dfsId").Int()
		if dfsId == 0 {
			dfsId, _ = j.Get("songs").GetIndex(0).Get("mMusic").Get("dfsId").Int()
		} else if dfsId == 0 {
			dfsId, _ = j.Get("songs").GetIndex(0).Get("lMusic").Get("dfsId").Int()
		}
		encrypted_song_id := encrypt_id(strconv.Itoa(dfsId))
		song.Url = "http://m1.music.126.net/" + encrypted_song_id + "/" + strconv.Itoa(dfsId) + ".mp3"
		if dfsId == 0 {
			song.Url = albumGetUrl(strconv.Itoa(song.AlbumId), songId, h)
		}
	}
	song.SongId , _= strconv.Atoi(songId)
	song.Url = strings.Replace(song.Url, "http://m", "http://p", -1)
	song.Album, err = j.Get("songs").GetIndex(0).Get("album").Get("name").String()
	song.AlbumPic, err = j.Get("songs").GetIndex(0).Get("album").Get("picUrl").String()
	song.Artists, err = j.Get("songs").GetIndex(0).Get("artists").GetIndex(0).Get("name").String()
	b = true
	return
}

func albumGetUrl(albumId string, songId string, h bool) (url string) {
	s := HttpGet("http://music.163.com/api/album/" + albumId + "?id=" + albumId)
	j, err := json.NewJson(s)
	if err != nil {
		beego.Error("NewJson")
		return ""
	}
	code, err := j.Get("code").Int()
	if err != nil {
		beego.Error("CodeToInt")
		return ""
	}
	if code != 200 {
		return ""
	}
	songs := j.Get("album").Get("songs")
	arr, _ := songs.Array()
	for i := 0; i < len(arr); i++ {
		id, _ := songs.GetIndex(i).Get("id").Int()
		sid, _ := strconv.Atoi(songId)
		if id == sid {
			if !h {
				url, _ = songs.GetIndex(i).Get("mp3Url").String()
			} else {
				dfsId, _ := songs.GetIndex(i).Get("hMusic").Get("dfsId").Int()
				if dfsId == 0 {
					dfsId, _ = songs.GetIndex(i).Get("mMusic").Get("dfsId").Int()
				} else if dfsId == 0 {
					dfsId, _ = songs.GetIndex(i).Get("lMusic").Get("dfsId").Int()
				}
				encrypted_song_id := encrypt_id(strconv.Itoa(dfsId))
				url = "http://m1.music.126.net/" + encrypted_song_id + "/" + strconv.Itoa(dfsId) + ".mp3"
			}
		}
	}
	return
}
func GetSongDetail(songId string) (s song) {
	j, err := json.NewJson(HttpGet("http://music.163.com/api/song/detail/?id=" + songId + "&ids=%5B" + songId + "%5D"))
	if err != nil {
		beego.Error("NewJson")
		return
	}
	code, err := j.Get("code").Int()
	if err != nil {
		beego.Error("CodeToInt")
		return
	}
	if code != 200 {
		s = song{Name: "查询失败", Artists: "查询失败"}
		return
	}
	s.Name, err = j.Get("songs").GetIndex(0).Get("name").String()
	if err != nil {
		s = song{Name: "查询失败", Artists: "查询失败"}
		return
	}
	s.Album, _ = j.Get("songs").GetIndex(0).Get("album").GetIndex(0).Get("name").String()
	s.Artists, _ = j.Get("songs").GetIndex(0).Get("artists").GetIndex(0).Get("name").String()
	s.AlbumPic,_ = j.Get("songs").GetIndex(0).Get("album").Get("picUrl").String()
	s.SongId, _ = strconv.Atoi(songId)
	return
}

func SongSearch(songKey string) (list []song) {
	str := HttpPost("http://music.163.com/api/search/pc", "offset=0&limit=100&type=1&s=" + songKey)
	j, err := json.NewJson(str)
	if err != nil {
		beego.Error(err)
		return
	}
	code, _ := j.Get("code").Int()
	if code != 200 {
		return
	}
	songs := j.Get("result").Get("songs")
	songCount, _ := j.Get("result").Get("songCount").Int()
	if songCount > 100 {
		songCount = 100
	}
	for i := 0; i < songCount; i++ {
		song := GetSong(songs.GetIndex(i), false)
		song.Id = i + 1
		list = append(list, song)
	}
	return
}

func GetSong(j *json.Json, url bool) (s song) {
	s.SongId, _ = j.Get("id").Int()
	s.Name, _ = j.Get("name").String()
	s.Artists, _ = j.Get("artists").GetIndex(0).Get("name").String()
	s.Album, _ = j.Get("album").Get("name").String()
	s.AlbumId, _ = j.Get("album").Get("id").Int()
	s.AlbumPic, _ = j.Get("album").Get("picUrl").String()
	s.AlbumPic += "?param=200y200"
	if url {
		s.Url, _ = j.Get("mp3Url").String()
	}
	return
}
func listDetail(listId string) (listsong []song) {
	s := HttpGet("http://music.163.com/api/playlist/detail?id=" + listId)
	j, err := json.NewJson(s)
	if err != nil {
		beego.Error("NewJson")
		return
	}
	code, err := j.Get("code").Int()
	if err != nil {
		beego.Error("CodeToInt")
		return
	}
	if code != 200 {
		return
	}
	l, _ := j.Get("result").Get("tracks").Array()
	h := len(l)

	for i := 0; i < h; i++ {
		var so song
		so.Id = i
		so.SongId, _ = j.Get("result").Get("tracks").GetIndex(i).Get("id").Int()
		so.Name, _ = j.Get("result").Get("tracks").GetIndex(i).Get("artists").GetIndex(0).Get("name").String()
		so.AlbumPic,_ = j.Get("result").Get("tracks").GetIndex(i).Get("album").Get("picUrl").String()
		so.AlbumPic += "?param=200y200"
		artists,_ :=  j.Get("result").Get("tracks").GetIndex(i).Get("artists").Array()
		so.Artists = ""
		for _,artist := range artists{
			if artist ,ok := artist.(map[string]interface{});ok {
				if name,ok := artist["name"].(string);ok {
					so.Artists += name
					so.Artists += ","
				}
			}

		}
		listsong = append(listsong, so)
	}
	return
}

func artistList(artistid string) (listsong []song){
	s := HttpGet("http://music.163.com/api/artist/" + artistid)
	j, err := json.NewJson(s)
	if err != nil {
		beego.Error("NewJson")
		return
	}
	code, err := j.Get("code").Int()
	if err != nil {
		beego.Error("CodeToInt")
		return
	}
	if code != 200 {
		return
	}
	hot := j.Get("hotSongs")
	arr,_ := hot.Array()
	l := len(arr)
	for i:=0 ;i<l;i++{
		var s  song
		s.Id = i
		s.SongId,_ = hot.GetIndex(i).Get("id").Int()
		s.Name,_ = hot.GetIndex(i).Get("name").String()
		s.Url = ""
		s.AlbumPic,_ = hot.GetIndex(i).Get("picUrl").String()
		s.AlbumPic += "?param=200y200"
		artists,_ := hot.GetIndex(i).Get("artists").Array()
		s.Artists = ""
		for _,artist := range artists{
			if artist ,ok := artist.(map[string]interface{});ok {
				if name,ok := artist["name"].(string);ok {
					s.Artists += name
					s.Artists += ","
				}
			}

		}
		listsong = append(listsong, s)
	}
	return
}
// https://github.com/yanunon/NeteaseCloudMusic/wiki/%E7%BD%91%E6%98%93%E4%BA%91%E9%9F%B3%E4%B9%90API%E5%88%86%E6%9E%90
func encrypt_id(id string) string {
	byte1 := []byte("3go8&$8*3*3h0k(2)2")
	byte2 := []byte(id)
	for i := 0; i < len(byte2); i++ {
		byte2[i] = byte2[i] ^ byte1[i % len(byte1)]
	}
	m := md5.New()
	m.Write(byte2)
	s := base64.StdEncoding.EncodeToString(m.Sum(nil))
	m.Reset()
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)
	return s
}
