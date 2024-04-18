package web

import (
	"fmt"
	"github.com/Duke1616/ecmdb/internal/attribute"
	"github.com/Duke1616/ecmdb/internal/relation/internal/domain"
	"github.com/Duke1616/ecmdb/internal/relation/internal/service"
	"github.com/Duke1616/ecmdb/internal/resource"
	"github.com/Duke1616/ecmdb/pkg/ginx"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type RelationResourceHandler struct {
	svc          service.RelationResourceService
	attributeSvc attribute.Service
	resourceSvc  resource.Service
}

func NewRelationResourceHandler(svc service.RelationResourceService, attributeSvc attribute.Service,
	resourceSvc resource.Service) *RelationResourceHandler {
	return &RelationResourceHandler{
		svc:          svc,
		attributeSvc: attributeSvc,
		resourceSvc:  resourceSvc,
	}
}

func (h *RelationResourceHandler) RegisterRoute(server *gin.Engine) {
	g := server.Group("/relation/resource")
	// 资源关联关系
	g.POST("/create", ginx.WrapBody[CreateResourceRelationReq](h.CreateResourceRelation))
	g.POST("/list/all", ginx.WrapBody[Page](h.ListResourceRelation))

	// 新建关联，查询所有的关联信息
	g.POST("/list-name", ginx.WrapBody[ListResourceRelationByModelUidReq](h.ListResourceByModelUid))

	// 拓补图
	g.POST("/diagram", ginx.WrapBody[ListResourceDiagramReq](h.ListDiagram))

	// 列表展示
	g.POST("/list-src", ginx.WrapBody[ListResourceDiagramReq](h.ListSrcResource))
	g.POST("/list-dst", ginx.WrapBody[ListResourceDiagramReq](h.ListDstResource))
	g.POST("/list", ginx.WrapBody[ListResourceDiagramReq](h.List))

	// 列表聚合展示、通过聚合处理
	g.POST("/pipeline/list-src", ginx.WrapBody[ListResourceDiagramReq](h.ListSrcAggregated))
	g.POST("/pipeline/list-dst", ginx.WrapBody[ListResourceDiagramReq](h.ListDstAggregated))
	g.POST("/pipeline/all", ginx.WrapBody[ListResourceDiagramReq](h.ListAllAggregated))

	// 比对，查询已经关联的节点

}

func (h *RelationResourceHandler) CreateResourceRelation(ctx *gin.Context, req CreateResourceRelationReq) (ginx.Result, error) {
	resp, err := h.svc.CreateResourceRelation(ctx, domain.ResourceRelation{
		SourceModelUID:   req.SourceModelUID,
		TargetModelUID:   req.TargetModelUID,
		RelationTypeUID:  req.RelationTypeUID,
		SourceResourceID: req.SourceResourceID,
		TargetResourceID: req.TargetResourceID,
	})

	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Msg:  "创建资源关联关系成功",
		Data: resp,
	}, nil
}

func (h *RelationResourceHandler) ListResourceRelation(ctx *gin.Context, req Page) (ginx.Result, error) {
	m, _, err := h.svc.ListResourceRelation(ctx, req.Offset, req.Limit)
	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Msg:  "查询资源关联成功",
		Data: m,
	}, nil
}

// ListResourceByModelUid 查询模型下，可以进行关联的数据
func (h *RelationResourceHandler) ListResourceByModelUid(ctx *gin.Context, req ListResourceRelationByModelUidReq) (
	ginx.Result, error) {
	projection, err := h.attributeSvc.SearchAttributeFiled(ctx, req.ModelUid)
	if err != nil {
		return systemErrorResult, fmt.Errorf("查询字段属性失败: %w", err)
	}

	ids, err := h.svc.ListResourceIds(ctx, req.ModelUid, req.RelationType)
	if err != nil {
		return systemErrorResult, fmt.Errorf("查询 resource ids失败: %w", err)
	}

	resources, err := h.resourceSvc.ListResourceByIds(ctx, projection, ids)
	if err != nil {
		return systemErrorResult, fmt.Errorf("查询resources列表失败: %w", err)
	}

	return ginx.Result{
		Data: resources,
	}, nil
}

// ListDiagram 入参 BASE UID 和 resource id
func (h *RelationResourceHandler) ListDiagram(ctx *gin.Context, req ListResourceDiagramReq) (ginx.Result, error) {
	// 1. 查询SRC模型 UID 放到 SRC
	// 2. 查询DST模型 UID 放到 DST
	diagram, _, err := h.svc.ListDiagram(ctx, req.ModelUid, req.ResourceId)
	if err != nil {
		return ginx.Result{}, err
	}

	// 1. 查询模型的所有关联
	// 2.
	return ginx.Result{
		Data: diagram,
	}, nil
}

func (h *RelationResourceHandler) ListSrcResource(ctx *gin.Context, req ListResourceDiagramReq) (ginx.Result, error) {
	rs, err := h.svc.ListSrcResources(ctx, req.ModelUid, req.ResourceId)
	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Data: rs,
	}, nil
}

func (h *RelationResourceHandler) ListDstResource(ctx *gin.Context, req ListResourceDiagramReq) (ginx.Result, error) {
	rs, err := h.svc.ListDstResources(ctx, req.ModelUid, req.ResourceId)
	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Data: rs,
	}, nil
}

func (h *RelationResourceHandler) List(ctx *gin.Context, req ListResourceDiagramReq) (ginx.Result, error) {
	var (
		eg   errgroup.Group
		srcS []domain.ResourceRelation
		dstS []domain.ResourceRelation
	)

	eg.Go(func() error {
		var err error
		srcS, err = h.svc.ListSrcResources(ctx, req.ModelUid, req.ResourceId)
		return err
	})

	eg.Go(func() error {
		var err error
		dstS, err = h.svc.ListDstResources(ctx, req.ModelUid, req.ResourceId)
		return err
	})
	if err := eg.Wait(); err != nil {
		return systemErrorResult, err
	}
	result := append(srcS, dstS...)

	return ginx.Result{
		Data: result,
	}, nil
}

func (h *RelationResourceHandler) ListSrcAggregated(ctx *gin.Context, req ListResourceDiagramReq) (ginx.Result, error) {
	list, err := h.svc.ListSrcAggregated(ctx, req.ModelUid, req.ResourceId)
	if err != nil {
		return ginx.Result{}, err
	}

	return ginx.Result{
		Data: list,
	}, nil
}

func (h *RelationResourceHandler) ListDstAggregated(ctx *gin.Context, req ListResourceDiagramReq) (ginx.Result, error) {
	list, err := h.svc.ListDstAggregated(ctx, req.ModelUid, req.ResourceId)
	if err != nil {
		return ginx.Result{}, err
	}

	return ginx.Result{
		Data: list,
	}, nil
}

func (h *RelationResourceHandler) ListAllAggregated(ctx *gin.Context, req ListResourceDiagramReq) (ginx.Result, error) {
	var (
		eg   errgroup.Group
		srcS []domain.ResourceAggregatedData
		dstS []domain.ResourceAggregatedData
	)

	eg.Go(func() error {
		var err error
		srcS, err = h.svc.ListSrcAggregated(ctx, req.ModelUid, req.ResourceId)
		return err
	})

	eg.Go(func() error {
		var err error
		dstS, err = h.svc.ListDstAggregated(ctx, req.ModelUid, req.ResourceId)
		return err
	})
	if err := eg.Wait(); err != nil {
		return systemErrorResult, err
	}

	result := append(srcS, dstS...)
	return ginx.Result{
		Data: result,
	}, nil
}
