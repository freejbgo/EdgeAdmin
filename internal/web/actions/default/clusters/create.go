package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CreateAction struct {
	actionutils.ParentAction
}

func (this *CreateAction) Init() {
	this.Nav("", "cluster", "create")
}

func (this *CreateAction) RunGet(params struct{}) {
	hasDomainsResp, err := this.RPC().DNSDomainRPC().ExistAvailableDomains(this.AdminContext(), &pb.ExistAvailableDomainsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasDomains"] = hasDomainsResp.Exist

	this.Show()
}

func (this *CreateAction) RunPost(params struct {
	Name string

	// SSH相关
	GrantId    int64
	InstallDir string

	// DNS相关
	DnsDomainId int64
	DnsName     string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入集群名称")

	// 检查DNS名称
	if len(params.DnsName) > 0 {
		if !domainutils.ValidateDomainFormat(params.DnsName) {
			this.FailField("dnsName", "请输入正确的DNS子域名")
		}

		// 检查是否已经被使用
		resp, err := this.RPC().NodeClusterRPC().CheckNodeClusterDNSName(this.AdminContext(), &pb.CheckNodeClusterDNSNameRequest{
			NodeClusterId: 0,
			DnsName:       params.DnsName,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if resp.IsUsed {
			this.FailField("dnsName", "此DNS子域名已经被使用，请换一个再试")
		}
	}

	// TODO 检查DnsDomainId的有效性

	createResp, err := this.RPC().NodeClusterRPC().CreateNodeCluster(this.AdminContext(), &pb.CreateNodeClusterRequest{
		Name:        params.Name,
		GrantId:     params.GrantId,
		InstallDir:  params.InstallDir,
		DnsDomainId: params.DnsDomainId,
		DnsName:     params.DnsName,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	this.CreateLog(oplogs.LevelInfo, "创建节点集群：%d", createResp.ClusterId)

	this.Success()
}
