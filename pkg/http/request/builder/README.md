Пример http запроса

```go
package app

func (repo *UserRepository) GetUserByTokenAndRealm(token, userRealm string) (*structures2.User, *exception.AppException) {
	builder := httpRequestBuilder.NewGetHttpRequestBuilder[structures2.User](repo.restConfig.UserServiceHost + "/user/profile/by-token")
	builder.SetRequestHeaders(map[string]string{
		constants.HeaderUserRouteAuthorization: token,
		constants.HeaderUserRouteRealm:         userRealm,
	})
	builder.SetRequestData(repo.appInfo)
	result, err := builder.GetResult()

	if err != nil {
		return nil, exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	if result.StatusCode != 200 {
		return nil, exception.NewAppException(result.StatusCode, errors.New(builder.GetErrorByKey("message")), nil)
	}

	return &result.Result, nil
}

```