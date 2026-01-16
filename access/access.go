package access

import (
	"fingelpp/utils"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type Permission = string
type UserKey = string

const (
	PermissionEditLesson Permission = "editl"
	UserKeyAnonymous     UserKey    = "anonymous_user"
)

type AccessFile struct {
	permissions map[UserKey][]Permission
}

var CurrentAccessFile *AccessFile = loadAccessFile("access.txt")

func loadAccessFile(path string) *AccessFile {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &AccessFile{permissions: map[UserKey][]Permission{}}
		}
		panic(err)
	}

	permissions := map[UserKey][]Permission{}
	for line := range strings.SplitSeq(string(content), "\n") {
		line = strings.TrimSpace(line)
		split := strings.Split(line, ":")
		if len(split) == 0 {
			continue
		}
		key := split[0]
		permissions[key] = split[1:]

	}

	return &AccessFile{permissions: permissions}
}

func (acc *AccessFile) KeyHasPermission(key UserKey, permission Permission) bool {
	if acc.permissions[key] == nil {
		return false
	}
	return slices.Contains(acc.permissions[key], permission)
}

func (acc *AccessFile) HasPermission(c *gin.Context, permission Permission) bool {
	return acc.KeyHasPermission(GetKey(c), permission)
}

func (acc *AccessFile) KeyExist(key string) bool {
	return acc.permissions[key] != nil
}

func (acc *AccessFile) GetPerms(key UserKey) []Permission {
	return acc.permissions[key]
}

func (acc *AccessFile) EnforcePermission(c *gin.Context, perm Permission) bool {
	if !acc.HasPermission(c, perm) {
		utils.ReqError(c, http.StatusForbidden)
		return true
	}
	return false
}

func GetKey(c *gin.Context) UserKey {
	key, err := c.Cookie("AccessKey")
	if err != nil {
		return UserKeyAnonymous
	}
	return UserKey(key)
}

func Routes(r *gin.Engine) {

	r.GET("/access/key/:key", func(c *gin.Context) {
		key := c.Param("key")
		if !CurrentAccessFile.KeyExist(key) {
			utils.ReqError(c, http.StatusUnauthorized)
			return
		}
		c.SetCookie("AccessKey", key, 100000000, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.GET("/access/permissions", func(c *gin.Context) {
		perms := CurrentAccessFile.GetPerms(GetKey(c))
		var sb strings.Builder
		for _, perm := range perms {
			sb.WriteString(string(perm))
		}
		c.String(http.StatusOK, sb.String())
	})
}
