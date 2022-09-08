package webhook

import (
	"context"

	"github.com/samueltorres/rbak/internal/auditor"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func NewAuditorWebhook(auditor *auditor.Auditor) *webhook.Admission {
	return &webhook.Admission{
		Handler: admission.HandlerFunc(func(ctx context.Context, req webhook.AdmissionRequest) webhook.AdmissionResponse {
			err := auditor.Audit(ctx, req)
			if err != nil {
				return webhook.Allowed("there was an error")
			}
			return webhook.Allowed("ok")
		}),
	}
}
