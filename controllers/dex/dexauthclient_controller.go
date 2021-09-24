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

package dex

import (
	"context"
	"fmt"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dexv1 "github.com/xzzpig/dex-configurer/apis/dex/v1"
	"github.com/xzzpig/dex-configurer/utils"
)

const finalizerName = "dex.xzzpig.com/finalizer"
const labelName = "dex.xzzpig.com/md5."

// DexAuthClientReconciler reconciles a DexAuthClient object
type DexAuthClientReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dex.xzzpig.com,resources=dexauthclients,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dex.xzzpig.com,resources=dexauthclients/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dex.xzzpig.com,resources=dexauthclients/finalizers,verbs=update

func (r *DexAuthClientReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var dexClient dexv1.DexAuthClient
	if err := r.Get(ctx, req.NamespacedName, &dexClient); err != nil {
		log.V(1).Info("Unable to fetch DexAuthClient")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	defer r.Status().Update(ctx, &dexClient)

	if dexClient.Status.Finish {
		dexClient.Status.Reset()
		log.V(1).Info("Reset DexAuthClient")
	}

	if dexClient.Spec.SecretRef.NameSpace == "" {
		dexClient.Spec.SecretRef.NameSpace = dexClient.ObjectMeta.Namespace
	}

	var secret corev1.Secret
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: dexClient.Spec.SecretRef.NameSpace,
		Name:      dexClient.Spec.SecretRef.Name,
	}, &secret); err != nil {
		var message = fmt.Sprintf("Unable to get Secret(%v/%v)", dexClient.Spec.SecretRef.NameSpace, dexClient.Spec.SecretRef.Name)
		log.Error(err, message)
		dexClient.Status.Error(message)
		return ctrl.Result{}, err
	}
	dexConfigData := secret.Data[dexClient.Spec.SecretRef.Key]
	if len(dexConfigData) == 0 {
		var message = fmt.Sprintf("Unable to get Secret Data(%v/%v/%v)", dexClient.Spec.SecretRef.NameSpace, dexClient.Spec.SecretRef.Name, dexClient.Spec.SecretRef.Key)
		log.Error(fmt.Errorf(message), message)
		dexClient.Status.Error(message)
		return ctrl.Result{}, nil
	}
	dexConfig, err := GetConfigFromBytes(dexConfigData)
	if err != nil {
		var message = "Unable to parse DexConfig"
		log.Error(err, message)
		dexClient.Status.Error(message)
		return ctrl.Result{}, err
	}
	if dexClient.DeletionTimestamp.IsZero() {
		dexConfig.SetStaticClient(&dexClient.Spec)
		secretData := []byte(dexConfig.GoString())
		secret.Data[dexClient.Spec.SecretRef.Key] = secretData
		if err := r.Update(ctx, &secret); err != nil {
			var message = "Unable to Update Secret Data"
			log.Error(err, message)
			dexClient.Status.Error(message)
			return ctrl.Result{}, err
		}
		log.V(1).Info("Updated Secret Data")

		var deps appv1.DeploymentList
		if err := r.List(ctx, &deps, &client.ListOptions{
			Namespace: secret.Namespace,
		}); err != nil {
			var message = "Unable to List Pods"
			log.Error(err, message)
			dexClient.Status.Error(message)
			return ctrl.Result{}, err
		}
		for _, deployment := range deps.Items {
			for _, volume := range deployment.Spec.Template.Spec.Volumes {
				if volume.Secret == nil || volume.Secret.SecretName != secret.Name {
					continue
				}
				log.Info("Redeploy Deployment using Targt Secret", "pod", deployment.Namespace+"/"+deployment.Name)
				deployment.Spec.Template.Labels[labelName+secret.Name] = utils.Md5Bytes(secretData)
				if err := r.Update(ctx, &deployment); err != nil {
					var message = "Unable to Redeploy Deployment using Target Secret"
					log.Error(err, message)
					dexClient.Status.Error(message)
					return ctrl.Result{}, err
				}
			}
		}

		log.Info("DexAuthClient Update Success")
		dexClient.Status.DoSuccess()

		if !utils.ContainsString(dexClient.Finalizers, finalizerName) {
			dexClient.Finalizers = append(dexClient.Finalizers, finalizerName)
			if err := r.Update(ctx, &dexClient); err != nil {
				log.Error(err, "Unable to add DexAuthClient Finalizers")
				return ctrl.Result{}, err
			}
			log.Info("Add Finalizers to DexAuthClient")
		}
	} else {
		dexConfig.RemoveStaticClient(dexClient.Spec.Id)
		secretData := []byte(dexConfig.GoString())
		secret.Data[dexClient.Spec.SecretRef.Key] = secretData
		if err := r.Update(ctx, &secret); err != nil {
			var message = "Unable to update Secret Data"
			log.Error(err, message)
			dexClient.Status.Error(message)
			return ctrl.Result{}, err
		}
		log.V(1).Info("Removed StaticClient In Secret Data")

		var deps appv1.DeploymentList
		if err := r.List(ctx, &deps, &client.ListOptions{
			Namespace: secret.Namespace,
		}); err != nil {
			var message = "Unable to List Pods"
			log.Error(err, message)
			dexClient.Status.Error(message)
			return ctrl.Result{}, err
		}
		for _, deployment := range deps.Items {
			for _, volume := range deployment.Spec.Template.Spec.Volumes {
				if volume.Secret == nil || volume.Secret.SecretName != secret.Name {
					continue
				}
				log.Info("Redeploy Deployment using Targt Secret", "pod", deployment.Namespace+"/"+deployment.Name)
				deployment.Spec.Template.Labels[labelName+secret.Name] = utils.Md5Bytes(secretData)
				if err := r.Update(ctx, &deployment); err != nil {
					var message = "Unable to Redeploy Deployment using Target Secret"
					log.Error(err, message)
					dexClient.Status.Error(message)
					return ctrl.Result{}, err
				}
			}
		}

		if utils.ContainsString(dexClient.Finalizers, finalizerName) {
			dexClient.Finalizers = utils.RemoveString(dexClient.Finalizers, finalizerName)
			if err := r.Update(ctx, &dexClient); err != nil {
				log.Error(err, "Unable to add DexAuthClient Finalizers")
				return ctrl.Result{}, err
			}
			log.Info("Remove Finalizers from DexAuthClient")
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DexAuthClientReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dexv1.DexAuthClient{}).
		Complete(r)
}
