package common

import (
	"context"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type R12nContext struct {
	context.Context
	Client                 client.Client
	Log                    logr.Logger
	SetControllerReference func(controlled metav1.Object) error
}
