package session

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session/v3"
)

const (
	managerKey = "gihub.com/orico-uc-backend/session/manager"
	storeKey   = "gihub.com/orico-uc-backend/session/store"
)

type Config struct {
	Secret   []byte
	Domain   string // optional
	Expired  int64
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
	Store    session.ManagerStore
}

func New(name string, config *Config) func(*gin.Context) {
	manager := session.NewManager(
		session.SetSign(config.Secret),
		session.SetCookieName(name),
		session.SetDomain(config.Domain),
		session.SetSecure(config.Secure),
		session.SetCookieLifeTime(config.MaxAge),
		session.SetExpired(config.Expired),
		session.SetSameSite(config.SameSite),
		session.SetStore(config.Store),
	)
	return func(c *gin.Context) {
		c.Set(managerKey, manager)
		c.Request.URL.RawQuery, _ = url.PathUnescape(c.Request.URL.RawQuery)
		c.Request.URL.RawQuery = escape(c.Request.URL.RawQuery)
		store, err := manager.Start(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			log.Fatalf("初始化session时redis连接错误: %+v", err)
		}
		c.Set(storeKey, store)
		c.Next()
	}
}

func escape(query string) string {
	q := query
	ascii := []string{"%2F", "%3A", "%3D", "%3F", "%20", "%25", "%26"}
	for i := 0; i < len(q); i++ {
		if q[i] == '%' && !slices.Contains(ascii, q[i:i+3]) {
			q = q[:i+1] + "25" + q[i+1:]
		}
	}
	return q
}

func NewWithContext(c *gin.Context) session.Store {
	if manager, ok := c.Get(managerKey); ok {
		store, err := manager.(*session.Manager).Start(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			panic(err)
		}
		return store
	}
	return nil
}

func Default(c *gin.Context) session.Store {
	store, ok := c.Get(storeKey)
	if ok {
		return store.(session.Store)
	}
	return nil
}

func Destroy(c *gin.Context) error {
	manager, ok := c.Get(managerKey)
	if ok {
		return manager.(*session.Manager).Destroy(c.Request.Context(), c.Writer, c.Request)
	}
	return errors.New("invalid session manager")
}

func Refresh(c *gin.Context) (session.Store, error) {
	manager, ok := c.Get(managerKey)
	if ok {
		return manager.(*session.Manager).Refresh(c.Request.Context(), c.Writer, c.Request)
	}
	return nil, errors.New("invalid session manager")
}
