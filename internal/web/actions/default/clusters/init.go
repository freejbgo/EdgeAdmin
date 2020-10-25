package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/clusterutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(clusterutils.NewClustersHelper()).
			Prefix("/clusters").
			Get("", new(IndexAction)).
			GetPost("/create", new(CreateAction)).
			Post("/sync", new(SyncAction)).
			Post("/checkChange", new(CheckChangeAction)).
			EndAll()
	})
}
