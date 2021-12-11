package vip

import (
	"encoding/json"
	"fmt"
	"regexp"

	//"regexp"

	//	"time"
	//	"crypto/md5"
	//	"encoding/hex"
	//	"unicode/utf8"
	//	"strings"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	//	"github.com/buger/jsonparser"
)

var vip = core.NewBucket("vip")

//订单侠apikey
var apikey = vip.Get("apikey")

//商品详情
var title string = ""
var url string = ""
var market_price string = ""
var vip_price string = ""

//淘宝商品结构体
type Item struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		NoEvokeURL     string `json:"noEvokeUrl"`
		VipQuickAppURL string `json:"vipQuickAppUrl"`
		VipWxURL       string `json:"vipWxUrl"`
		DeeplinkURL    string `json:"deeplinkUrl"`
		LongURL        string `json:"longUrl"`
		Source         string `json:"source"`
		UlURL          string `json:"ulUrl"`
		URL            string `json:"url"`
		TraFrom        string `json:"traFrom"`
		NoEvokeLongURL string `json:"noEvokeLongUrl"`
		ItemInfo       struct {
			GoodsID               string        `json:"goodsId"`
			GoodsName             string        `json:"goodsName"`
			MarketPrice           string        `json:"marketPrice"`
			VipPrice              string        `json:"vipPrice"`
			Commission            string        `json:"commission"`
			CommissionRate        string        `json:"commissionRate"`
			Discount              string        `json:"discount"`
			CouponInfo            []interface{} `json:"couponInfo"`
			GoodsCarouselPictures []string      `json:"goodsCarouselPictures"`
			GoodsDetailPictures   []string      `json:"goodsDetailPictures"`
			BrandLogoFull         string        `json:"brandLogoFull"`
			GoodsThumbURL         string        `json:"goodsThumbUrl"`
			GoodsMainPicture      string        `json:"goodsMainPicture"`
			StoreInfo             struct {
				StoreName string `json:"storeName"`
				StoreID   string `json:"storeId"`
			} `json:"storeInfo"`
		} `json:"itemInfo"`
	} `json:"data"`
}

func init() {

	core.Server.GET("/vip/:sku", func(c *gin.Context) {
		sku := c.Param("sku")
		c.String(200, core.OttoFuncs["vip"](sku))
	})
	//添加命令
	core.AddCommand("", []core.Function{
		{ /*https://m.vip.com/product-1710612828-6919227375621606108.html?
			nmsns=shop_android-7.58.7-weixin&amp;nst=product&amp;nsbc=&amp;
			nct=card&amp;ncid=106458c5-abf9-3f0d-929a-aca9c47296e5&amp;
			nabtid=121&amp;nuid=423461972&amp;nchl_param=share:106458c5-abf9-3f0d-929a-aca9c47296e5:1639174223182&amp;
			mars_cid_a=106458c5-abf9-3f0d-929a-aca9c47296e5&amp;chl_type=share
			*/
			Rules: []string{"raw https?://m\\.vip\\.com/",
							"raw https?://t\\.vip\\.com/"},
			Handle: func(s core.Sender) interface{} {
				return getvip(s.GetContent())
			},
		},
	})
	core.OttoFuncs["vip"] = getvip //类似于向核心组件注册
}

func getvip(info string) string {
	fmt.Println("开始处理唯品会：")
	var rlt = ""
	params := regexp.MustCompile(`product\-(\d+)\-(\d+)\.html`).FindStringSubmatch(info)
	for _, param:=range params{
		fmt.Println(param)
	}
	if len(params) >= 2 {
		goodsId := params[2]
		url := getUrlConvert(goodsId)
		rlt = title + "\n市场价：" + market_price + "\n唯品价：" + vip_price + "\n惠购链接：" + url
	}
	return rlt
}

/*
通过商品id获取淘宝客推广链接
*/
func getUrlConvert(iids string) string {
	//根据id获取推广链接
	req := httplib.Get("http://api.tbk.dingdanxia.com/vip/id_privilege?" +
		"apikey=" + apikey +
		"&id=" + iids +
		"&itemInfo=true")
	data, _ := req.Bytes()
	fmt.Println("-------------------------------\n" + string(data))
	//itemURL, _ := jsonparser.GetString(data, "data","itemInfo","item_url")
	res := &Item{}
	json.Unmarshal([]byte(data), &res)
	if res.Data.ItemInfo.GoodsName != "" {
		title = res.Data.ItemInfo.GoodsName
	}
	market_price = res.Data.ItemInfo.MarketPrice
	vip_price = res.Data.ItemInfo.VipPrice
	//fmt.Println(res.Data.ItemURL)
	return res.Data.URL
}

// 创建一个错误处理函数，避免过多的 if err != nil{} 出现
func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}
