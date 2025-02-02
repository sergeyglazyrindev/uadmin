package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/uadmin"
	sessionsblueprint "github.com/sergeyglazyrindev/uadmin/blueprint/sessions"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type AdminPopulateContextModelActionTestSuite struct {
	uadmin.TestSuite
}

func (suite *AdminPopulateContextModelActionTestSuite) TestPopulateContext() {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
	session := sessionAdapter.Create()
	user := core.GenerateUserModel()
	user.SetFirstName("testadmincontext")
	user.SetLastName("testadmincontext")
	user.SetUsername("testadmincontext")
	user.SetEmail("testadmincontext@example.com")
	user.SetActive(true)
	user.SetIsSuperUser(true)
	user.SetIsPasswordUsable(true)
	core.NewUadminDatabase().Db.Create(reflect.ValueOf(user).Interface())
	session.SetUser(user)
	session.Save()
	ctx.Request.AddCookie(&http.Cookie{
		Name:     core.CurrentConfig.D.Uadmin.AdminCookieName,
		Value:    url.QueryEscape(session.GetKey()),
		MaxAge:   int(core.CurrentConfig.D.Uadmin.SessionDuration),
		Path:     "/",
		Domain:   ctx.Request.URL.Host,
		SameSite: http.SameSiteLaxMode,
		Secure:   core.CurrentConfig.D.Uadmin.SecureCookie,
		HttpOnly: core.CurrentConfig.D.Uadmin.HTTPOnlyCookie,
	})
	adminContext := &core.AdminContext{}
	adminRequestParams := core.NewAdminRequestParams()
	adminRequestParams.NeedAllLanguages = true
	adminRequestParams.GenerateCSRFToken = true
	core.PopulateTemplateContextForAdminPanel(ctx, adminContext, adminRequestParams)
	sessionKey := session.GetKey()
	session, _ = sessionAdapter.GetByKey(sessionKey)
	csrfToken, _ := session.Get("csrf_token")
	assert.True(suite.T(), csrfToken != "")
	assert.True(suite.T(), adminContext.GetSessionKey() != "")
	assert.True(suite.T(), adminContext.UserExists)
	assert.True(suite.T(), adminContext.UserPermissionRegistry != nil)
	assert.True(suite.T(), len(adminContext.Languages) > 0)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPopulateAdminContext(t *testing.T) {
	uadmin.RunTests(t, new(AdminPopulateContextModelActionTestSuite))
}
