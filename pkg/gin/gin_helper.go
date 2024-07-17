package gin

import (
	"github.com/exgamer/gosdk/pkg/config"
	"github.com/exgamer/gosdk/pkg/constants"
	"github.com/exgamer/gosdk/pkg/exception"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	timeout "github.com/vearne/gin-timeout"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net/http"
	"time"
)

// InitRouter Базовая инициализация gin
func InitRouter(baseConfig *config.BaseConfig) *gin.Engine {
	if baseConfig.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Options
	router := gin.New()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "404 page not found"})
	})
	router.HandleMethodNotAllowed = true
	p := ginprometheus.NewPrometheus("ginHelpers")
	p.Use(router)
	router.Use(sentrygin.New(sentrygin.Options{}))
	router.Use(gin.Logger())
	router.Use(timeout.Timeout(timeout.WithTimeout(time.Duration(baseConfig.HandlerTimeout) * time.Second)))
	router.Use(gin.CustomRecovery(ErrorHandler))

	return router
}

// ErrorHandler Обработчик ошибок gin
func ErrorHandler(c *gin.Context, err any) {
	goErr := errors.Wrap(err, 2)
	details := make([]string, 0)

	for _, frame := range goErr.StackFrames() {
		details = append(details, frame.String())
	}

	sentry.CaptureException(goErr)
	c.JSON(http.StatusInternalServerError, gin.H{"message": goErr.Error(), "details": details})
}

func Error(c *gin.Context, exception *exception.AppException) {
	c.Set("exception", exception)
	c.Status(exception.Code)
}

func Success(c *gin.Context, data any) {
	c.Set("data", data)
}

func SetAppInfo(c *gin.Context, baseConfig *config.BaseConfig) {
	appInfo := config.AppInfo{}
	appInfo.RequestId = c.GetHeader(constants.RequestIdHeaderName)
	// если request id не пришел с заголовком, генерим его, чтобы прокидывать дальше при http запросах
	if appInfo.RequestId == "" {
		appInfo.GenerateRequestId()
		c.Request.Header.Add(constants.RequestIdHeaderName, appInfo.RequestId)
	}

	appInfo.LanguageCode = c.GetHeader(constants.LanguageHeaderName)

	if appInfo.LanguageCode == "" {
		appInfo.LanguageCode = constants.LangRu
	}

	appInfo.RequestUrl = c.Request.URL.Path
	appInfo.RequestMethod = c.Request.Method
	appInfo.RequestScheme = c.Request.URL.Scheme
	appInfo.RequestHost = c.Request.Host
	appInfo.AppEnv = baseConfig.AppEnv
	appInfo.ServiceName = baseConfig.Name
	c.Set("app_info", appInfo)
}

func GetAppInfo(c *gin.Context) *config.AppInfo {
	value, _ := c.Get("app_info")
	appInfo := value.(config.AppInfo)

	return &appInfo
}
