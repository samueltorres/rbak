package auditor

import (
	"context"
	"strings"
	"sync"

	"github.com/go-logr/logr"
	rbakv1alpha1 "github.com/samueltorres/rbak/api/v1alpha1"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/serviceaccount"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Worker struct {
	wg     *sync.WaitGroup
	cli    client.Client
	logger logr.Logger
	reqs   chan admissionv1.AdmissionRequest
}

func NewWorker(logger logr.Logger, cli client.Client, wg *sync.WaitGroup, reqs chan admissionv1.AdmissionRequest) *Worker {
	return &Worker{
		logger: logger,
		cli:    cli,
		wg:     wg,
		reqs:   reqs,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case req := <-w.reqs:
			err := w.handleAdmission(ctx, req)
			if err != nil {
				w.logger.Error(err, "error handling admission request")
			}
		}
	}
}

func (w *Worker) handleAdmission(ctx context.Context, req admissionv1.AdmissionRequest) error {
	subject := userToSubject(req.UserInfo)
	username := strings.Replace(subject.Name, ":", "-", -1)

	ns := "kube-system"
	if subject.Kind == "ServiceAccount" {
		ns = subject.Namespace
	}

	rbacReport := rbakv1alpha1.RbacReport{
		ObjectMeta: v1.ObjectMeta{
			Name:      username,
			Namespace: ns,
		},
	}

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := controllerutil.CreateOrUpdate(ctx, w.cli, &rbacReport, func() error {
			rbacReport.Spec.Subject = subject
			rbacReport.Spec.Rules = addRule(rbacReport.Spec.Rules, req)
			return nil
		})
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func addRule(rules []rbakv1alpha1.Rules, req admissionv1.AdmissionRequest) []rbakv1alpha1.Rules {
	ruleIndex := -1
	for i := 0; i < len(rules); i++ {
		if !contains(rules[i].APIGroups, req.Kind.Group) {
			continue
		}
		if !contains(rules[i].Resources, req.Resource.Resource) {
			continue
		}
		if rules[i].Namespace != req.Namespace {
			continue
		}
		ruleIndex = i
		break
	}

	operation := strings.ToLower(string(req.Operation))

	if ruleIndex == -1 {
		rules = append(rules, rbakv1alpha1.Rules{
			Namespace: req.Namespace,
			APIGroups: []string{req.Kind.Group},
			Resources: []string{req.Resource.Resource},
			Verbs:     []string{operation},
		})
		return rules
	}

	if !contains(rules[ruleIndex].Verbs, operation) {
		rules[ruleIndex].Verbs = append(rules[ruleIndex].Verbs, operation)
	}

	return rules
}

func userToSubject(user authenticationv1.UserInfo) rbacv1.Subject {
	if ns, name, err := serviceaccount.SplitUsername(user.Username); err == nil {
		return rbacv1.Subject{Name: name, Namespace: ns, Kind: "ServiceAccount"}
	}
	return rbacv1.Subject{Name: user.Username, Kind: "User", APIGroup: rbacv1.GroupName}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
