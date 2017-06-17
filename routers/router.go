// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"netease/controllers"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego"
)

func init() {
	beego.Any("/", func(ctx  *context.Context) {
		ctx.Output.Body([]byte("hello world!"))
	})
	nsSong := beego.NewNamespace("/song",
		beego.NSRouter("/search", &controllers.SongSearchControler{}),
		beego.NSRouter("/detail", &controllers.SongDetailControler{}),
		beego.NSRouter("/cdndetail", &controllers.SongCdnDetailControler{}),
	)
	nsPlaylist := beego.NewNamespace("/playlist",
		beego.NSRouter("", &controllers.PlaylistControler{}),
	)
	nsArtistList := beego.NewNamespace("/artist",
		beego.NSRouter("", &controllers.ArtistListControler{}),
	)
	beego.AddNamespace(nsSong,nsPlaylist,nsArtistList)
}
