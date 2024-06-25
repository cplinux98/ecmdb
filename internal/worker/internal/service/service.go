package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Duke1616/ecmdb/internal/runner"
	"github.com/Duke1616/ecmdb/internal/worker/internal/domain"
	"github.com/Duke1616/ecmdb/internal/worker/internal/repository"
	"github.com/ecodeclub/mq-api"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	RegisterWorker(ctx context.Context, req domain.Worker) error
	FindOrRegisterByName(ctx context.Context, req domain.Worker) (domain.Worker, error)
	FindOrRegisterByKey(ctx context.Context, req domain.Worker) (domain.Worker, error)
	ListWorker(ctx context.Context, offset, limit int64) ([]domain.Worker, int64, error)
}

type service struct {
	repo      repository.WorkerRepository
	runnerSvc runner.Service
	mq        mq.MQ
}

func NewService(mq mq.MQ, runnerSvc runner.Service, repo repository.WorkerRepository) Service {
	return &service{
		mq:        mq,
		repo:      repo,
		runnerSvc: runnerSvc,
	}
}

func (s *service) RegisterWorker(ctx context.Context, req domain.Worker) error {
	//TODO implement me
	panic("implement me")
}

func (s *service) FindOrRegisterByName(ctx context.Context, req domain.Worker) (domain.Worker, error) {
	worker, err := s.repo.FindByName(ctx, req.Name)
	if !errors.Is(err, repository.ErrUserNotFound) {
		if req.Status != worker.Status {
			_, err = s.repo.UpdateStatus(ctx, worker.Id, domain.Status.ToUint8(req.Status))
			if err != nil {
				return worker, fmt.Errorf("修改状态失败: %x", err)
			}
			worker.Status = req.Status
		}

		return worker, err
	}

	// 新增工作节点
	id, err := s.repo.CreateWorker(ctx, req)
	if err != nil {
		return domain.Worker{}, fmt.Errorf("创建节点失败: %x", err)
	}
	worker.Id = id

	// 新增 Topic
	if err = s.mq.CreateTopic(ctx, req.Topic, 1); err != nil {
		return domain.Worker{}, fmt.Errorf("创建Topic失败: %x", err)
	}

	// 新增 producer 监听
	if err = s.runnerSvc.CreateProducer(req.Topic); err != nil {
		return domain.Worker{}, fmt.Errorf("创建Topic失败: %x", err)
	}
	return worker, nil
}

func (s *service) FindOrRegisterByKey(ctx context.Context, req domain.Worker) (domain.Worker, error) {
	worker, err := s.repo.FindByKey(ctx, req.Key)
	if !errors.Is(err, repository.ErrUserNotFound) {
		if req.Status != worker.Status {
			_, err = s.repo.UpdateStatus(ctx, worker.Id, domain.Status.ToUint8(req.Status))
			if err != nil {
				return worker, fmt.Errorf("修改状态失败: %x", err)
			}
			worker.Status = req.Status
		}

		return worker, err
	}

	// 新增工作节点
	id, err := s.repo.CreateWorker(ctx, req)
	if err != nil {
		return domain.Worker{}, fmt.Errorf("创建节点失败: %x", err)
	}
	worker.Id = id

	// 新增 Topic
	if err = s.mq.CreateTopic(ctx, req.Topic, 1); err != nil {
		return domain.Worker{}, fmt.Errorf("创建Topic失败: %x", err)
	}

	// 新增 producer 监听
	if err = s.runnerSvc.CreateProducer(req.Topic); err != nil {
		return domain.Worker{}, fmt.Errorf("创建Topic失败: %x", err)
	}
	return worker, nil
}

func (s *service) ListWorker(ctx context.Context, offset, limit int64) ([]domain.Worker, int64, error) {
	var (
		eg    errgroup.Group
		ts    []domain.Worker
		total int64
	)
	eg.Go(func() error {
		var err error
		ts, err = s.repo.ListWorker(ctx, offset, limit)
		return err
	})

	eg.Go(func() error {
		var err error
		total, err = s.repo.Total(ctx)
		return err
	})
	if err := eg.Wait(); err != nil {
		return ts, total, err
	}
	return ts, total, nil
}
