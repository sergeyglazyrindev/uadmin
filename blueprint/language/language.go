package language

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/blueprint/language/migrations"
	"github.com/uadmin/uadmin/core"
	"mime/multipart"
	"strconv"
)

type Blueprint struct {
	core.Blueprint
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	languageAdminPage := core.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) { return nil, nil },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	languageAdminPage.PageName = "Languages"
	languageAdminPage.Slug = "language"
	languageAdminPage.BlueprintName = "language"
	languageAdminPage.Router = mainRouter
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(languageAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
	languagemodelAdminPage := core.NewGormAdminPage(
		languageAdminPage,
		func() (interface{}, interface{}) { return &core.Language{}, &[]*core.Language{} },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"EnglishName", "Name", "Flag", "Code", "RTL", "Default", "Active", "AvailableInGui"}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			defaultField, _ := form.FieldRegistry.GetByName("Default")
			defaultField.Validators.AddValidator("only_one_default_language", func(i interface{}, o interface{}) error {
				isDefault := i.(bool)
				if !isDefault {
					return nil
				}
				d := o.(*multipart.Form)
				ID := d.Value["ID"][0]
				uadminDatabase := core.NewUadminDatabase()
				lang := &core.Language{}
				uadminDatabase.Db.Where(&core.Language{Default: true}).First(lang)
				if lang.ID != 0 && ID != strconv.Itoa(int(lang.ID)) {
					return fmt.Errorf("only one default language could be configured")
				}
				return nil
			})
			return form
		},
	)
	languagemodelAdminPage.PageName = "Languages"
	languagemodelAdminPage.Slug = "language"
	languagemodelAdminPage.BlueprintName = "language"
	languagemodelAdminPage.Router = mainRouter
	languagemodelAdminPage.NoPermissionToAddNew = true
	err = languageAdminPage.SubPages.AddAdminPage(languagemodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing language blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	core.ProjectModels.RegisterModel(func() interface{} { return &core.Language{} })
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "language",
		Description:       "Language blueprint is responsible for managing languages used in the project",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
