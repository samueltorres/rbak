package auditor

import (
	"context"
	"sync"
	"time"

	admissionv1 "k8s.io/api/admission/v1"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type Auditor struct {
	cli         client.Client
	logger      logr.Logger
	workerCount int32
	reqs        chan admissionv1.AdmissionRequest
}

func New(cli client.Client, logger logr.Logger, workerCount int32) *Auditor {
	return &Auditor{
		cli:         cli,
		logger:      logger,
		workerCount: workerCount,
		reqs:        make(chan admissionv1.AdmissionRequest),
	}
}

func (a *Auditor) Start(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	for i := 0; i < int(a.workerCount); i++ {
		wg.Add(1)
		worker := NewWorker(a.logger, a.cli, wg, a.reqs)
		go worker.Start(ctx)
	}

	<-ctx.Done()
	wg.Wait()
	return nil
}

func (a *Auditor) Audit(ctx context.Context, req webhook.AdmissionRequest) error {
	cctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	select {
	case a.reqs <- req.AdmissionRequest:
		a.logger.Info("Handled the admission request in time")
		return nil
	case <-cctx.Done():
		a.logger.Info("Couldn't handle the admission request in time")
		return cctx.Err()
	}
}
