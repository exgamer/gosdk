package builder

import (
	"encoding/json"
	"encoding/xml"
	"github.com/exgamer/gosdk/pkg/config"
	"github.com/exgamer/gosdk/pkg/constants"
	"github.com/exgamer/gosdk/pkg/http/structures"
	"github.com/exgamer/gosdk/pkg/logger"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// NewPostHttpRequestBuilder - Новый построитель rest запросов для POST
func NewPostHttpRequestBuilder[E interface{}](url string) *HttpRequestBuilder[E] {
	return &HttpRequestBuilder[E]{
		url:                      url,
		method:                   "POST",
		timeout:                  30 * time.Second,
		transport:                loghttp.Transport{},
		responseDataPresentation: constants.JSON,
	}
}

// NewGetHttpRequestBuilder - Новый построитель rest запросов для GET
func NewGetHttpRequestBuilder[E interface{}](url string) *HttpRequestBuilder[E] {
	return &HttpRequestBuilder[E]{
		url:                      url,
		method:                   "GET",
		timeout:                  30 * time.Second,
		throwUnmarshalError:      false,
		transport:                loghttp.Transport{},
		responseDataPresentation: constants.JSON,
	}
}

// HttpRequestBuilder - Построитель rest запросов
type HttpRequestBuilder[E interface{}] struct {
	url                      string
	method                   string
	headers                  map[string]string
	throwUnmarshalError      bool
	body                     io.Reader
	appInfo                  *config.AppInfo
	timeout                  time.Duration
	transport                loghttp.Transport
	request                  *http.Request
	response                 *structures.HttpResponse[E]
	responseDataPresentation string
	result                   E
}

func (builder *HttpRequestBuilder[E]) SetResponseDataPresentation(dataPresentation string) *HttpRequestBuilder[E] {
	builder.responseDataPresentation = dataPresentation

	return builder
}

func (builder *HttpRequestBuilder[E]) SetHttpRequest(request *http.Request) *HttpRequestBuilder[E] {
	builder.request = request

	return builder
}

// SetRequestData - установить Доп данные для запроса (используется для логирования)
func (builder *HttpRequestBuilder[E]) SetRequestData(appInfo *config.AppInfo) *HttpRequestBuilder[E] {
	builder.appInfo = appInfo

	return builder
}

// SetRequestHeaders - установить заголовки запроса
func (builder *HttpRequestBuilder[E]) SetRequestHeaders(headers map[string]string) *HttpRequestBuilder[E] {
	builder.headers = headers

	return builder
}

// SetRequestBody - установить тело запроса
func (builder *HttpRequestBuilder[E]) SetRequestBody(body io.Reader) *HttpRequestBuilder[E] {
	builder.body = body

	return builder
}

// SetRequestTimeout - установить таймаут запроса
func (builder *HttpRequestBuilder[E]) SetRequestTimeout(timeout time.Duration) *HttpRequestBuilder[E] {
	builder.timeout = timeout

	return builder
}

// SetRequestTransport - установить параметры запроса
func (builder *HttpRequestBuilder[E]) SetRequestTransport(transport loghttp.Transport) *HttpRequestBuilder[E] {
	builder.transport = transport

	return builder
}

func (builder *HttpRequestBuilder[E]) do() error {
	client := http.Client{
		Timeout:   builder.timeout,
		Transport: &builder.transport,
	}

	if builder.request == nil {
		builder.request, _ = http.NewRequest(builder.method, builder.url, builder.body)
	}

	for n, v := range builder.headers {
		builder.request.Header.Set(n, v)
	}

	response, err := client.Do(builder.request)
	builder.response = &structures.HttpResponse[E]{
		Url:     builder.url,
		Method:  builder.method,
		Headers: builder.headers,
	}

	if err != nil {
		logger.LogError(err)

		return err
	}

	builder.response.Status = response.Status
	builder.response.StatusCode = response.StatusCode

	rBody, bErr := io.ReadAll(response.Body)

	if bErr != nil {
		logger.LogError(bErr)

		return bErr
	}

	builder.response.Body = rBody

	defer response.Body.Close()

	return nil
}

// Do - выполнить запрос
func (builder *HttpRequestBuilder[E]) Do() error {
	messageBuilder := strings.Builder{}
	start := time.Now()
	err := builder.do()

	if err != nil {
		logger.LogError(err)

		return err
	}

	execTime := time.Since(start)

	if err != nil {
		logger.LogError(err)
		messageBuilder.WriteString("Url: " + builder.method + " " + builder.url)
		messageBuilder.WriteString(" Error:" + err.Error())

		if builder.appInfo != nil {
			logger.FormattedLogWithAppInfo(builder.appInfo, messageBuilder.String())
		}

		return err
	}

	messageBuilder.WriteString("Url: " + builder.response.Method + " " + builder.response.Status + " " + builder.response.Url)
	messageBuilder.WriteString(" Exec time:" + execTime.String())

	if builder.response.StatusCode >= 400 {
		messageBuilder.WriteString(" Response:" + string(builder.response.Body))
		err := json.Unmarshal(builder.response.Body, &builder.response.ErrorsMap)

		if err != nil {
			log.Fatal(err)
		}
	}

	if builder.appInfo != nil {
		logger.FormattedLogWithAppInfo(builder.appInfo, messageBuilder.String())
	}

	return nil
}

// GetResult  Возвращает результат
func (builder *HttpRequestBuilder[E]) GetResult() (*structures.HttpResponse[E], error) {
	err := builder.Do()

	if err != nil {
		return nil, err
	}

	switch builder.responseDataPresentation {
	case constants.JSON:
		unMarshErr := json.Unmarshal(builder.response.Body, &builder.result)
		if unMarshErr != nil && builder.throwUnmarshalError {
			return nil, unMarshErr
		}
	case constants.XML:
		unMarshErr := xml.Unmarshal(builder.response.Body, &builder.result)
		if unMarshErr != nil && builder.throwUnmarshalError {
			return nil, unMarshErr
		}
	default:
		unMarshErr := json.Unmarshal(builder.response.Body, &builder.result)
		if unMarshErr != nil && builder.throwUnmarshalError {
			return nil, unMarshErr
		}
	}

	builder.response.Result = builder.result

	return builder.response, nil
}

// GetErrorByKey Возвращает ключ из ошибки
func (builder *HttpRequestBuilder[E]) GetErrorByKey(key string) string {
	message, found := builder.response.ErrorsMap[key]

	if found {
		return message.(string)
	}

	return "unknown"
}
