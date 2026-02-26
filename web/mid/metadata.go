package mid

import "github.com/gin-gonic/gin"

const ExclusivelyTitleKey = "/gin/exclusively/title"

// Title set the title value in `*gin.Context`.
func Title(title string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ExclusivelyTitleKey, title)
		c.Next()
	}
}

// GetTitle returns the title value stored in `*gin.Context`, if any.
func GetTitle(c *gin.Context) string {
	return c.GetString(ExclusivelyTitleKey)
}
