/*
Copyright 2021 xzzpig.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package networkingk8sio

import (
	"context"
	"strings"

	dexv1 "github.com/xzzpig/dex-configurer/apis/dex/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	AnnotationOAuthProxy = "dex.xzzpig.com/oauth-proxy"
)

// IngressReconciler reconciles a Ingress object
type IngressReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=networking.k8s.io.xzzpig.com,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io.xzzpig.com,resources=ingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.k8s.io.xzzpig.com,resources=ingresses/finalizers,verbs=update

func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var ingress netv1.Ingress
	if err := r.Get(ctx, req.NamespacedName, &ingress); err != nil {
		log.V(1).Info("Unable to fetch Ingress")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if ingress.Annotations == nil {
		return ctrl.Result{}, nil
	}

	configName := strings.TrimLeft(ingress.Annotations[AnnotationOAuthProxy], "/")

	if configName == "" {
		return ctrl.Result{}, nil
	}

	var orderList dexv1.DexClientOrderList
	r.List(ctx, &orderList, client.MatchingFields{TargetIngressKey: ingress.Name})
	for _, order := range orderList.Items {
		if order.Spec.TargetIngress.Name != ingress.Name || order.Spec.TargetIngress.Namespace != ingress.Namespace {
			continue
		}
		return ctrl.Result{}, nil
	}

	order := dexv1.DexClientOrder{}
	order.Namespace = ingress.Namespace
	order.Name = ingress.Name
	order.Spec.Config = dexv1.NamespacedName{
		Name: configName,
	}
	order.Spec.TargetIngress = dexv1.NamespacedName{
		Namespace: ingress.Namespace,
		Name:      ingress.Name,
	}

	if err := r.Create(ctx, &order); err != nil {
		log.Error(err, "Unable to Create DexClientOrder")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

const (
	TargetIngressKey = ".metadata.targetingress"
)

func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &dexv1.DexClientOrder{}, TargetIngressKey, func(rawObj client.Object) []string {
		order := rawObj.(*dexv1.DexClientOrder)
		return []string{order.Spec.TargetIngress.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&netv1.Ingress{}).
		Complete(r)
}
