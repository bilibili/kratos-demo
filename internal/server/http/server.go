package http

import (
	"net/http"

	pb "kratos-demo/api"
	"kratos-demo/internal/model"
	"kratos-demo/internal/service"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

var (
	svc *service.Service
)

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine) {
	var (
		hc struct {
			Server *bm.ServerConfig
		}
	)
	if err := paladin.Get("http.toml").UnmarshalTOML(&hc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	svc = s
	engine = bm.DefaultServer(hc.Server)
	pb.RegisterDemoBMServer(engine, svc)
	initRouter(engine)
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)                 // engine自带的"/ping"接口，用于负载均衡检测服务健康状态
	g := e.Group("/kratos-demo") // e.Group 创建一组 "/kratos-demo" 起始的路由组
	{
		g.GET("/start", howToStart) // g.GET 创建一个 "kratos-demo/start" 的路由，使用GET方式请求，默认处理Handler为howToStart方法
		g.POST("start", howToStart) // g.POST 创建一个 "kratos-demo/start" 的路由，使用POST方式请求，默认处理Handler为howToStart方法

		// 路径参数有两个特殊符号":"和"*"
		// ":" 跟在"/"后面为参数的key，匹配两个/中间的值 或 一个/到结尾(其中不再包含/)的值
		// "*" 跟在"/"后面为参数的key，匹配从 /*开始到结尾的所有值，所有*必须写在最后且无法多个

		// NOTE：这是不被允许的！！！会和 /start 冲突！！！
		// g.GET("/:xxx")

		// NOTE: 可以拿到一个key为name的参数。注意只能匹配到/param1/felix，无法匹配/param1/felix/hao(该路径会404)
		g.GET("/param1/:name", pathParam)
		// NOTE: 可以拿到多个key参数。注意只能匹配到/param2/felix/hao/love，无法匹配/param2/felix或/param2/felix/hao
		g.GET("/param2/:name/:value/:felid", pathParam)
		// NOTE: 可以拿到一个key为name的参数 和 一个key为action的路径。
		// NOTE: 如/params3/felix/hello，action的值为"/hello"
		// NOTE: 如/params3/felix/hello/hi，action的值为"/hello/hi"
		// NOTE: 如/params3/felix/hello/hi/，action的值为"/hello/hi/"
		g.GET("/param3/:name/*action", pathParam)
	}
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// example for http request handler.
func howToStart(c *bm.Context) {
	k := &model.Kratos{
		Hello: "Golang 大法好 !!!",
	}
	c.JSON(k, nil)
}

func pathParam(c *bm.Context) {
	name, _ := c.Params.Get("name")
	value, _ := c.Params.Get("value")
	felid, _ := c.Params.Get("felid")
	action, _ := c.Params.Get("action")
	path := c.RoutePath // NOTE: 获取注册的路由原始地址，如: /kratos-demo/param1/:name
	c.JSONMap(map[string]interface{}{
		"name":   name,
		"value":  value,
		"felid":  felid,
		"action": action,
		"path":   path,
	}, nil)
}
