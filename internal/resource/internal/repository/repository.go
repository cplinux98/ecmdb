package repository

import (
	"context"
	"github.com/Duke1616/ecmdb/internal/resource/internal/domain"
	"github.com/Duke1616/ecmdb/internal/resource/internal/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type ResourceRepository interface {
	CreateResource(ctx context.Context, req domain.Resource) (int64, error)
	FindResourceById(ctx context.Context, fields []string, id int64) (domain.Resource, error)
	ListResource(ctx context.Context, fields []string, modelUid string, offset, limit int64) ([]domain.Resource, error)
	Total(ctx context.Context, modelUid string) (int64, error)
	ListResourcesByIds(ctx context.Context, fields []string, ids []int64) ([]domain.Resource, error)
	ListExcludeResourceByIds(ctx context.Context, fields []string, modelUid string, offset, limit int64, ids []int64) ([]domain.Resource, error)
}

type resourceRepository struct {
	dao dao.ResourceDAO
}

func NewResourceRepository(dao dao.ResourceDAO) ResourceRepository {
	return &resourceRepository{
		dao: dao,
	}
}

func (r *resourceRepository) CreateResource(ctx context.Context, req domain.Resource) (int64, error) {
	return r.dao.CreateResource(ctx, r.toEntity(req))
}

func (r *resourceRepository) FindResourceById(ctx context.Context, fields []string, id int64) (domain.Resource, error) {
	rs, err := r.dao.FindResourceById(ctx, fields, id)
	return r.toDomain(rs), err
}

func (r *resourceRepository) ListResourcesByIds(ctx context.Context, fields []string, ids []int64) ([]domain.Resource, error) {
	rrs, err := r.dao.ListResourcesByIds(ctx, fields, ids)

	return slice.Map(rrs, func(idx int, src dao.Resource) domain.Resource {
		return r.toDomain(src)
	}), err
}

func (r *resourceRepository) ListResource(ctx context.Context, fields []string, modelUid string, offset, limit int64) ([]domain.Resource, error) {
	rrs, err := r.dao.ListResource(ctx, fields, modelUid, offset, limit)

	return slice.Map(rrs, func(idx int, src dao.Resource) domain.Resource {
		return r.toDomain(src)
	}), err
}

func (r *resourceRepository) Total(ctx context.Context, modelUid string) (int64, error) {
	return r.dao.Count(ctx, modelUid)
}

func (r *resourceRepository) ListExcludeResourceByIds(ctx context.Context, fields []string, modelUid string, offset, limit int64, ids []int64) ([]domain.Resource, error) {
	rrs, err := r.dao.ListExcludeResourceByids(ctx, fields, modelUid, offset, limit, ids)

	return slice.Map(rrs, func(idx int, src dao.Resource) domain.Resource {
		return r.toDomain(src)
	}), err
}

func (r *resourceRepository) toEntity(req domain.Resource) dao.Resource {
	return dao.Resource{
		ModelUID: req.ModelUID,
		Name:     req.Name,
		Data:     req.Data,
	}
}

func (r *resourceRepository) toDomain(src dao.Resource) domain.Resource {
	return domain.Resource{
		ID:       src.ID,
		ModelUID: src.ModelUID,
		Data:     src.Data,
		Name:     src.Name,
	}
}
