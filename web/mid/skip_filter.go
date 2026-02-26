package mid

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/sets"
)

// SkipFilter 跳过授权过滤器
// NOTE: 跳过登陆验证必定跳过权限验证
type SkipFilter struct {
	authorizes  sets.Set[string] // 跳过登陆验证
	permissions sets.Set[string] // 跳过权限验证
}

func NewSkipFilter() *SkipFilter {
	return &SkipFilter{
		authorizes:  sets.New[string](),
		permissions: sets.New[string](),
	}
}

// 增加 跳过登陆验证(同时权限验证也跳过)
func (s *SkipFilter) AddAuthorize(method string, paths ...string) *SkipFilter {
	for _, path := range paths {
		v := FormatMethodUri(method, path)
		s.authorizes.Insert(v)
		s.permissions.Insert(v)
	}
	return s
}

// 增加 跳过权限验证
func (s *SkipFilter) AddPermission(method string, paths ...string) *SkipFilter {
	for _, path := range paths {
		v := FormatMethodUri(method, path)
		s.permissions.Insert(v)
	}
	return s
}

func (s *SkipFilter) SkipAuthorize(c *gin.Context) bool {
	return s.authorizes.Contains(FormatMethodUri(c.Request.Method, c.FullPath()))
}

func (s *SkipFilter) SkipPermission(c *gin.Context) bool {
	return s.permissions.Contains(FormatMethodUri(c.Request.Method, c.FullPath()))
}

func (s *SkipFilter) ListAuthorize() []string {
	return slices.Collect(s.authorizes.All())
}

func (s *SkipFilter) ListPermission() []string {
	return slices.Collect(s.permissions.All())
}

func FormatMethodUri(method, path string) string {
	return method + "@" + path
}
