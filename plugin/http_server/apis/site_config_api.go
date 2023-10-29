// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package apis

import (
	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
)

type SiteConfigVo struct {
	SiteName string `json:"siteName"`
	Logo     string `json:"logo"`
	AppName  string `json:"appName"`
}

func UpdateSiteConfig(c *gin.Context, ruleEngine typex.RuleX) {

	form := SiteConfigVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := service.UpdateSiteConfig(model.MSiteConfig{
		SiteName: form.SiteName,
		Logo:     form.Logo,
		AppName:  form.AppName,
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
func GetSiteConfig(c *gin.Context, ruleEngine typex.RuleX) {
	Model, err := service.GetSiteConfig()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(SiteConfigVo{
		SiteName: Model.SiteName,
		Logo:     __tempLogo,
		AppName:  Model.AppName,
	}))
}

// 这是用来测试的后期会删除
var __tempLogo = `
data:image/png;base64,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
`
