package UrlController

import (
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/wxccs/tinyurl/app/Controllers"
	"github.com/wxccs/tinyurl/app/Models"
	"github.com/wxccs/tinyurl/global"
)

var base62Regex = regexp.MustCompile(`^[0-9a-zA-Z]+$`)

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

type ShortenResponse struct {
	ShortCode   string `json:"short_code"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

func Shorten(c *gin.Context) {
	funcName := "app.Controllers.UrlController.Shorten"

	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.WithField("func", funcName).Warn("invalid request body")
		Controllers.NewResponse(Controllers.RESPONSE_PARAMETERS_ERROR, "url is required", nil).ResponseJson(400, c)
		return
	}

	parsed, err := url.ParseRequestURI(req.URL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		global.Log.WithField("func", funcName).WithField("url", req.URL).Warn("invalid url")
		Controllers.NewResponse(Controllers.RESPONSE_PARAMETERS_ERROR, "invalid url, must be http or https", nil).ResponseJson(400, c)
		return
	}

	originalUrl := req.URL

	if isOwnShortUrl(c, parsed) {
		var existing Models.Url
		if result := global.DB.Where("short_code = ?", parsed.Path[1:]).First(&existing); result.Error == nil {
			originalUrl = existing.OriginalUrl
			global.Log.WithField("func", funcName).WithField("original_url", originalUrl).Info("resolved own short url to original")
		}
	}

	var existing Models.Url
	if result := global.DB.Where("original_url = ?", originalUrl).First(&existing); result.Error == nil {
		scheme := "http"
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		shortUrl := scheme + "://" + c.Request.Host + "/" + existing.ShortCode
		global.Log.WithField("func", funcName).WithField("short_code", existing.ShortCode).Info("returned existing short url")
		Controllers.NewResponse(Controllers.RESPONSE_OK, "success", ShortenResponse{
			ShortCode:   existing.ShortCode,
			ShortUrl:    shortUrl,
			OriginalUrl: originalUrl,
		}).ResponseJson(200, c)
		return
	}

	var urlModel Models.Url
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		code, genErr := global.Generator.Generate()
		if genErr != nil {
			global.Log.WithField("func", funcName).WithError(genErr).Error("failed to generate short code")
			Controllers.NewResponse(Controllers.RESPONSE_BACKEND_ERROR, "failed to generate short code", nil).ResponseJson(500, c)
			return
		}

		urlModel = Models.Url{
			ShortCode:   code,
			OriginalUrl: originalUrl,
		}

		result := global.DB.Create(&urlModel)
		if result.Error == nil {
			scheme := "http"
			if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
				scheme = "https"
			}
			shortUrl := scheme + "://" + c.Request.Host + "/" + code

			global.Log.WithField("func", funcName).WithField("short_code", code).Info("short url created")
			Controllers.NewResponse(Controllers.RESPONSE_OK, "success", ShortenResponse{
				ShortCode:   code,
				ShortUrl:    shortUrl,
				OriginalUrl: originalUrl,
			}).ResponseJson(200, c)
			return
		}

		global.Log.WithField("func", funcName).WithField("short_code", code).WithField("attempt", i+1).Warn("collision detected, retrying")
	}

	global.Log.WithField("func", funcName).Error("failed to generate unique short code after retries")
	Controllers.NewResponse(Controllers.RESPONSE_BACKEND_ERROR, "failed to generate unique short code", nil).ResponseJson(500, c)
}

func isOwnShortUrl(c *gin.Context, parsed *url.URL) bool {
	if parsed.Host != c.Request.Host {
		return false
	}
	path := parsed.Path
	if len(path) < 2 {
		return false
	}
	code := path[1:]
	return base62Regex.MatchString(code) && len(code) == global.Config.ShortURL.Length
}

func Redirect(c *gin.Context) {
	funcName := "app.Controllers.UrlController.Redirect"

	code := c.Param("code")
	if !base62Regex.MatchString(code) || len(code) != global.Config.ShortURL.Length {
		global.Log.WithField("func", funcName).WithField("code", code).Warn("invalid short code format")
		Controllers.NewResponse(Controllers.RESPONSE_PARAMETERS_ERROR, "invalid short url", nil).ResponseJson(400, c)
		return
	}

	var urlModel Models.Url
	result := global.DB.Where("short_code = ?", code).First(&urlModel)
	if result.Error != nil {
		global.Log.WithField("func", funcName).WithField("code", code).Warn("short code not found")
		Controllers.NewResponse(Controllers.RESPONSE_NOT_FOUND, "short url not found", nil).ResponseJson(400, c)
		return
	}

	global.Log.WithField("func", funcName).WithField("code", code).WithField("url", urlModel.OriginalUrl).Info("redirecting")
	c.Redirect(302, urlModel.OriginalUrl)
}
