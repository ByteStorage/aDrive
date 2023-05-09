package web

import (
	"aDrive/api"
	"aDrive/pkg/utils"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

var CurrentLeader string
var DB *gorm.DB

var (
	getCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_counter",
		Help: "A counter for get request",
	})
	putCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "put_counter",
		Help: "A counter for put request",
	})
	Place = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "place",
		Help: "My string metric",
	}, []string{"place"})
	UsedMem = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "used_mem",
		Help: "used memory",
	})
	UsedDisk = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "used_disk",
		Help: "used disk",
	})
	FreeMem = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "free_mem",
		Help: "total memory",
	})
	FreeDisk = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "free_disk",
		Help: "total disk",
	})
)

func init() {
	prometheus.NewRegistry().MustRegister(getCounter)
	prometheus.NewRegistry().MustRegister(putCounter)
	prometheus.MustRegister(Place)
	prometheus.MustRegister(UsedMem)
	prometheus.MustRegister(UsedDisk)
	prometheus.MustRegister(FreeMem)
	prometheus.MustRegister(FreeDisk)
	fmt.Println("init prometheus------------------")
}

func StartWeb() {
	initDatabase()
	store := cookie.NewStore([]byte("secret"))
	r := gin.Default()
	//设置跨域访问
	r.Use(Cors())
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/ping", func(c *gin.Context) {
		getCounter.Inc()
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/get/*filename", Get)
	r.GET("/visit/*filename", GetFormVisit)
	r.GET("/own/*path", Own)
	r.POST("/put", Put)
	r.GET("/showDashboard", ShowDashboard)
	r.GET("/mkdir/*path", Mkdir)
	r.GET("/delete/*path", Delete)
	r.GET("/list/*path", Index)
	r.POST("/login", Login)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	err := r.Run(":5862")
	if err != nil {
		panic(err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Own(c *gin.Context) {
	token := c.Query("token")
	username, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	ok, err := utils.ExpireToken(token)
	if err != nil || !ok {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	param := c.Param("path")
	path := "/" + username + param
	fmt.Println("username:", username)
	path, err = url.PathUnescape(path)
	decodedParam, err := url.QueryUnescape(path)
	if utils.Exit("decode param error", err) {
		return
	}
	fmt.Println("decodedParam:", decodedParam)
	listResp, err := api.List(CurrentLeader, decodedParam)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	token, err = utils.CreateToken(username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
		"list":  listResp,
	})
}

func Index(c *gin.Context) {
	param := c.Param("path")
	zap.L().Debug("param", zap.String("param", param))
	path, err := url.PathUnescape(param)
	decodedParam, err := url.QueryUnescape(path)
	if utils.Exit("decode param error", err) {
		return
	}
	fmt.Println("decodedParam:", decodedParam)
	listResp, err := api.List(CurrentLeader, decodedParam)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	token := c.Query("token")
	zap.L().Debug("token", zap.String("token", token))
	if token != "null" {
		username, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "token error",
			})
			return
		}
		ok, err := utils.ExpireToken(token)
		if err != nil || !ok {
			c.JSON(500, gin.H{
				"message": "token error",
			})
			return
		}
		token, err = utils.CreateToken(username)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	if path == "/" {
		for i := range listResp.Dirname {
			var user User
			err := DB.Model(&User{}).Where("username = ?", listResp.Dirname[i]).First(&user).Error
			if err != nil {
				c.JSON(500, gin.H{
					"message": err.Error(),
				})
				return
			}
			zap.L().Debug("user", zap.Any("user", user))
			if !user.Open {
				//删除
				listResp.Dirname = append(listResp.Dirname[:i], listResp.Dirname[i+1:]...)
			}
		}
	}
	c.JSON(200, gin.H{
		"list":  listResp,
		"token": token,
	})

}

func Get(c *gin.Context) {
	getCounter.Inc()
	token := c.Query("token")
	username, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	ok, err := utils.ExpireToken(token)
	if err != nil || !ok {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	fmt.Println("username:", username)
	param := c.Param("filename")
	param = "/" + username + param
	decodedParam, err := url.QueryUnescape(param)
	if utils.Exit("decode param error", err) {
		return
	}
	fmt.Println("decodedParam:", decodedParam)
	getResp, err := api.Get(CurrentLeader, decodedParam)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	filename := filepath.Base(param)
	file, err := ioutil.TempFile("", filename)
	if utils.Exit("create temp file error", err) {
		return
	}
	defer file.Close()
	_, err = file.Write(getResp.Data)
	if utils.Exit("write file error", err) {
		return
	}

	fmt.Println("filename:", filename)
	fmt.Println("文件下载成功")
	contentType := http.DetectContentType(getResp.Data)
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	c.Writer.Header().Set("Content-Type", contentType)
	c.File(file.Name())
	token, err = utils.CreateToken(username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
	})
}

func GetFormVisit(c *gin.Context) {
	getCounter.Inc()
	param := c.Param("filename")
	decodedParam, err := url.QueryUnescape(param)
	if utils.Exit("decode param error", err) {
		return
	}
	fmt.Println("decodedParam:", decodedParam)
	getResp, err := api.Get(CurrentLeader, decodedParam)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	filename := filepath.Base(param)
	file, err := ioutil.TempFile("", filename)
	if utils.Exit("create temp file error", err) {
		return
	}
	defer file.Close()
	_, err = file.Write(getResp.Data)
	if utils.Exit("write file error", err) {
		return
	}

	fmt.Println("filename:", filename)
	fmt.Println("文件下载成功")
	contentType := http.DetectContentType(getResp.Data)
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	c.Writer.Header().Set("Content-Type", contentType)

	token := c.Query("token")
	zap.L().Debug("token", zap.String("token", token))
	if token != "null" {
		username, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "token error",
			})
			return
		}
		ok, err := utils.ExpireToken(token)
		if err != nil || !ok {
			c.JSON(500, gin.H{
				"message": "token error",
			})
			return
		}
		token, err = utils.CreateToken(username)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	c.File(file.Name())
	c.JSON(200, gin.H{
		"token": token,
	})
}

func Put(c *gin.Context) {
	putCounter.Inc()
	token := c.Query("token")
	username, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	ok, err := utils.ExpireToken(token)
	if err != nil || !ok {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	param := c.PostForm("path")
	param = "/" + username + param
	decodedParam, err := url.QueryUnescape(param)
	if utils.Exit("decode param error", err) {
		return
	}
	fmt.Println("decodedParam:", decodedParam)
	path := utils.ModPath(decodedParam)
	file, err := c.FormFile("data")
	if utils.Exit("get file error", err) {
		return
	}
	fmt.Println(file.Filename)
	open, err := file.Open()
	if utils.Exit("open file error", err) {
		return
	}
	bytes, err := ioutil.ReadAll(open)
	if utils.Exit("read file error", err) {
		return
	}
	//读取文件内存大小
	needMemory := int64(len(bytes))
	var user User
	err = DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "用户不存在",
		})
		return
	}
	if user.Memory < needMemory {
		c.JSON(500, gin.H{
			"message": "用户内存不足",
		})
		return
	}
	user.Memory -= needMemory
	err = DB.Model(&User{}).Where("username = ?", username).Update("memory", user.Memory).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "用户内存更新失败",
		})
		return
	}
	join := path + file.Filename
	log.Println("absolutePath:", join)
	putResp, err := api.Put(CurrentLeader, join, bytes)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	token, err = utils.CreateToken(username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"token":   token,
		"putResp": putResp,
	})
}

func Login(c *gin.Context) {
	var user User
	err := c.ShouldBindWith(&user, binding.JSON)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "参数错误",
		})
		return
	}
	var find User
	err = DB.Where("username = ?", user.Username).First(&find).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "用户不存在",
		})
		return
	}
	if user.Password != find.Password {
		c.JSON(500, gin.H{
			"message": "密码错误",
		})
		return
	}
	token, err := utils.CreateToken(user.Username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "token生成失败",
		})
	}
	c.JSON(200, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}

func Delete(c *gin.Context) {
	token := c.Query("token")
	username, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	ok, err := utils.ExpireToken(token)
	if err != nil || !ok {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	fmt.Println("username:", username)
	fmt.Println("接收到POST请求")
	param := c.Param("path")
	param = "/" + username + param
	path := utils.ModPath(param)
	fmt.Println("path:", path)
	_, err = api.Delete(CurrentLeader, path)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	token, err = utils.CreateToken(username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
	})
}

func Mkdir(c *gin.Context) {
	token := c.Query("token")
	username, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	ok, err := utils.ExpireToken(token)
	if err != nil || !ok {
		c.JSON(500, gin.H{
			"message": "token error",
		})
		return
	}
	fmt.Println("username:", username)
	fmt.Println("接收到POST请求")
	param := c.Param("path")
	param = "/" + username + param
	zap.L().Info("param", zap.String("param", param))
	path := utils.ModPath(param)
	_, err = api.Mkdir(CurrentLeader, path)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	token, err = utils.CreateToken(username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
	})
}

func ShowDashboard(c *gin.Context) {
	fmt.Println("接收到GET请求")
	// 设置Grafana仪表板的URL
	url := "http://116.62.156.91:3000/d/vGS45kL4z/test"

	// 设置Grafana API的URL
	apiURL := fmt.Sprintf("%s/api/dashboards/%s", url, "vGS45kL4z")

	// 创建HTTP请求并添加Grafana API的Authorization头
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer ")
	req.Header.Set("Accept", "text/html")
	// 发送HTTP请求并读取响应内容
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// 将响应内容作为HTML发送到应用程序
	fmt.Println(string(body))
	c.JSON(200, gin.H{
		"message": string(body),
	})
}

func initDatabase() {
	db, err := gorm.Open(mysql.Open(""), &gorm.Config{})
	if err != nil {
		panic("application start fail")
	}
	DB = db
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}

type User struct {
	Username string `gorm:"type:varchar(20);not null;unique" json:"username"`
	Password string `gorm:"type:varchar(20);not null" json:"password"`
	Memory   int64  `gorm:"type:bigint;not null" json:"memory"`
	Open     bool   `gorm:"type:tinyint(1);not null" json:"open"`
}
