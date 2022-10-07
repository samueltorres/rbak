package webhook

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/samueltorres/rbak/internal/auditor"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func NewAuditorWebhook(auditor *auditor.Auditor, logger logr.Logger) *webhook.Admission {
	return &webhook.Admission{
		Handler: admission.HandlerFunc(func(ctx context.Context, req webhook.AdmissionRequest) webhook.AdmissionResponse {
			err := auditor.Audit(ctx, req)
			if err != nil {
				logger.Error(err, "error on webook")
				return webhook.Allowed("error on webhook")
			}
			return webhook.Allowed("ok")
		}),
	}
}
