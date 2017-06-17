# netease
## go版网易api 结合python3版musicbox网易命令行接口 和 go-music翻译的
### 此接口可以听收费歌曲返回的mp3地址不过期, 可以存数据库
### 演示地址opdays.com:8080

| URL        | function           | Method  |
| ------------- |:-------------:| -----:|
| /song/search?key=邓紫棋      | 搜索邓紫棋的歌曲id和图片列表 | GET |
| /song/detail?id=     | 返回指定id歌曲的mp3url      |   GET |
| /playlist?id=  | 返回指定id歌单的歌曲id和图片列表    |    GET |
| /artist?id=  | 返回指定id歌手的歌曲id和图片列表    |    GET |
