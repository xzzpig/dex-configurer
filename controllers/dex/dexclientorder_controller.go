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
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dexv1 "github.com/xzzpig/dex-configurer/apis/dex/v1"
	"github.com/xzzpig/dex-configurer/utils"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DexClientOrderReconciler reconciles a DexClientOrder object
type DexClientOrderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dex.xzzpig.com,resources=dexclientorders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dex.xzzpig.com,resources=dexclientorders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dex.xzzpig.com,resources=dexclientorders/finalizers,verbs=update
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets;configmaps;services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dex.xzzpig.com,resources=dexproxyconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dex.xzzpig.com,resources=dexproxyconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dex.xzzpig.com,resources=dexproxyconfigs/finalizers,verbs=update

func (r *DexClientOrderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var order dexv1.DexClientOrder
	if err := r.Get(ctx, req.NamespacedName, &order); err != nil {
		log.V(1).Info("Unable to fetch DexClientOrder")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	defer r.Status().Update(ctx, &order)

	//查找带出其他对象
	//DexProxyConfig
	var config dexv1.DexProxyConfig
	if err := r.Get(ctx, types.NamespacedName(order.Spec.Config), &config); err != nil {
		var message = "Unable to fetch DexProxyConfig"
		log.Error(err, message)
		order.Status.Message = message
		order.Status.Created = false
		return ctrl.Result{}, err
	}
	//Ingress
	var targetIngress netv1.Ingress
	if order.Spec.TargetIngress.Namespace == "" {
		order.Spec.TargetIngress.Namespace = order.Namespace
	}
	if err := r.Get(ctx, types.NamespacedName(order.Spec.TargetIngress), &targetIngress); err != nil {
		var message = "Unable to fetch Target Ingress"
		log.Error(err, message)
		order.Status.Message = message
		order.Status.Created = false
		return ctrl.Result{}, err
	}

	//生成填充可选字段
	if order.Spec.ClientId == "" {
		order.Spec.ClientId = strings.ReplaceAll(utils.UnMarshalCamel(strings.ReplaceAll(order.Name, "-", "_")), "_", "-")
		if err := r.Update(ctx, &order); err != nil {
			var message = "Unable to Update DexClientOrder ClientId"
			log.Error(err, message, "ClientId", order.Spec.ClientId)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Updated DexClientOrder ClientId", "ClientId", order.Spec.ClientId)
		return ctrl.Result{Requeue: true}, nil
	}
	if order.Spec.ClientName == "" {
		order.Spec.ClientName = utils.MarshalCamel(strings.ReplaceAll(order.Spec.ClientId, "-", "_"))
		if err := r.Update(ctx, &order); err != nil {
			var message = "Unable to Update DexClientOrder ClientName"
			log.Error(err, message, "ClientName", order.Spec.ClientName)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Updated DexClientOrder ClientName", "ClientName", order.Spec.ClientName)
		return ctrl.Result{Requeue: true}, nil
	}
	if order.Spec.ClientSecret == "" {
		order.Spec.ClientSecret = utils.Md5(utils.RandString(8))
		if err := r.Update(ctx, &order); err != nil {
			var message = "Unable to Update DexClientOrder ClientSecret"
			log.Error(err, message, "ClientSecret", order.Spec.ClientSecret)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Updated DexClientOrder ClientSecret", "ClientSecret", order.Spec.ClientSecret)
		return ctrl.Result{Requeue: true}, nil
	}
	if order.Spec.RedirectUrl.Scheme == "" {
		order.Spec.RedirectUrl.Scheme = config.Spec.DefaultUrl.Scheme
		if order.Spec.RedirectUrl.Scheme == "" {
			var message = "Empty DexClientOrder RedirectUrl Scheme and Not Found in DexProxyConfig DefaultUrl"
			err := errors.New(strings.ToLower(message))
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		} else {
			if err := r.Update(ctx, &order); err != nil {
				var message = "Unable to Update DexClientOrder RedirectUrl Scheme"
				log.Error(err, message)
				order.Status.Message = message
				order.Status.Created = false
				return ctrl.Result{}, err
			}
			log.Info("Updated DexClientOrder  RedirectUrl Scheme")
			return ctrl.Result{Requeue: true}, nil
		}
	}
	if order.Spec.RedirectUrl.Port == 0 || order.Spec.RedirectUrl.Path == "" {
		if order.Spec.RedirectUrl.Port == 0 {
			if config.Spec.DefaultUrl.Port != 0 {
				order.Spec.RedirectUrl.Port = config.Spec.DefaultUrl.Port
			} else if order.Spec.RedirectUrl.Scheme == dexv1.SchemeHTTPS {
				order.Spec.RedirectUrl.Port = 443
			} else {
				order.Spec.RedirectUrl.Port = 80
			}
		}
		if order.Spec.RedirectUrl.Path == "" {
			if config.Spec.DefaultUrl.Path == "" {
				order.Spec.RedirectUrl.Path = "/oauth2"
			} else {
				order.Spec.RedirectUrl.Path = config.Spec.DefaultUrl.Path
			}
		}
		if err := r.Update(ctx, &order); err != nil {
			var message = "Unable to Update DexClientOrder RedirectUrl Path&Port"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Updated DexClientOrder RedirectUrl Path&Port")
		return ctrl.Result{Requeue: true}, nil
	}
	if order.Spec.Config.Namespace != "" {
		order.Spec.Config.Namespace = ""
		if err := r.Update(ctx, &order); err != nil {
			var message = "Unable to Update DexClientOrder Config Namespace"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Updated DexClientOrder Config Namespace")
		return ctrl.Result{Requeue: true}, nil
	}
	if config.Spec.CookieSecret == "" {
		config.Spec.CookieSecret = utils.RandString(16)
		if err := r.Update(ctx, &config); err != nil {
			var message = "Unable to Update DexProxyConfig CookieSecret"
			log.Error(err, message, "CookieSecret", config.Spec.CookieSecret)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Updated DexProxyConfig CookieSecret", "CookieSecret", config.Spec.CookieSecret)
		return ctrl.Result{Requeue: true}, nil
	}
	if config.Spec.ProxyImage == "" {
		config.Spec.ProxyImage = "quay.io/oauth2-proxy/oauth2-proxy:v7.4.0"
		if err := r.Update(ctx, &config); err != nil {
			var message = "Unable to Update DexProxyConfig ProxyImage"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Updated DexProxyConfig ProxyImage")
		return ctrl.Result{Requeue: true}, nil
	}
	if config.Spec.SecretRef.NameSpace == "" {
		config.Spec.SecretRef.NameSpace = config.Namespace
	}
	if config.Spec.ProviderDisplayName == "" {
		config.Spec.ProviderDisplayName = "Dex Login"
	}
	if config.Spec.AuthCacheKey == "" {
		config.Spec.AuthCacheKey = "$cookie__oauth2_proxy"
	}

	if order.Spec.RedirectUrl.Host == "" {
		for _, rule := range targetIngress.Spec.Rules {
			if rule.Host != "" {
				order.Spec.RedirectUrl.Host = rule.Host
				if err := r.Update(ctx, &order); err != nil {
					var message = "Unable to Update DexClientOrder RedirectUrl Host"
					log.Error(err, message)
					order.Status.Message = message
					order.Status.Created = false
					return ctrl.Result{}, err
				}
				log.Info("Updated DexClientOrder RedirectUrl Host")
				return ctrl.Result{Requeue: true}, nil
			}
		}
		var message = "DexClientOrder RedirectUrl Host is Empty but No Host Found in Target Ingress"
		err := errors.New(message)
		log.Error(err, message)
		order.Status.Message = message
		order.Status.Created = false
		return ctrl.Result{}, err
	}

	//处理Client
	var client dexv1.DexAuthClient
	clientNew := false
	if order.Status.RefObjects.ClientRef.Namespace == "" {
		order.Status.RefObjects.ClientRef.Namespace = order.Namespace
	}
	if order.Status.RefObjects.ClientRef.Name == "" {
		order.Status.RefObjects.ClientRef.Name = order.Name + "-client"
		for {
			if err := r.Get(ctx, types.NamespacedName(order.Status.RefObjects.ClientRef), &dexv1.DexAuthClient{}); err != nil {
				break
			}
			order.Status.RefObjects.ClientRef.Name = order.Name + "-client-" + strings.ToLower(utils.RandString(4))
		}
		client.Name = order.Status.RefObjects.ClientRef.Name
		client.Namespace = order.Status.RefObjects.ClientRef.Namespace
		clientNew = true
	} else {
		if err := r.Get(ctx, types.NamespacedName(order.Status.RefObjects.ClientRef), &client); err != nil {
			log.Error(err, "Unable to fetch DexAuthClient,Will Create New One")
			clientNew = true
		}
	}
	client.Spec.Id = order.Spec.ClientId
	client.Spec.Name = order.Spec.ClientName
	client.Spec.Secret = order.Spec.ClientSecret
	client.Spec.RedirectURIs = append([]string{order.Spec.RedirectUrl.GoString()}, order.Spec.ExtraRedirectUrls...)
	client.Spec.SecretRef = config.Spec.SecretRef
	if err := ctrl.SetControllerReference(&order, &client, r.Scheme); err != nil {
		var message = "Unable to Set Controller Reference on DexAuthClient"
		log.Error(err, message)
		order.Status.Message = message
		order.Status.Created = false
		return ctrl.Result{}, err
	}

	if clientNew {
		if err := r.Create(ctx, &client); err != nil {
			var message = "Unable to Create DexAuthClient"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("DexAuthClient Created")
	} else {
		if err := r.Update(ctx, &client); err != nil {
			var message = "Unable to Update DexAuthClient"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("DexAuthClient Updated")
	}

	//处理Deployment
	var deployment appv1.Deployment
	deploymentNew := false
	if order.Status.RefObjects.DeploymentRef.Namespace == "" {
		order.Status.RefObjects.DeploymentRef.Namespace = order.Namespace
	}
	if order.Status.RefObjects.DeploymentRef.Name == "" {
		order.Status.RefObjects.DeploymentRef.Name = order.Name + "-dex-proxy"
		for {
			if err := r.Get(ctx, types.NamespacedName(order.Status.RefObjects.DeploymentRef), &appv1.Deployment{}); err != nil {
				break
			}
			order.Status.RefObjects.DeploymentRef.Name = order.Name + "-dex-proxy-" + strings.ToLower(utils.RandString(4))
		}
		deployment.Name = order.Status.RefObjects.DeploymentRef.Name
		deployment.Namespace = order.Status.RefObjects.DeploymentRef.Namespace
		deploymentNew = true
	} else {
		if err := r.Get(ctx, types.NamespacedName(order.Status.RefObjects.DeploymentRef), &deployment); err != nil {
			log.Error(err, "Unable to fetch Deployment,Will Create New One")
			deploymentNew = true
		}
	}
	createConfigMap := func() error {
		var configMap corev1.ConfigMap
		configMapNew := false
		if err := r.Get(ctx, types.NamespacedName(order.Status.RefObjects.DeploymentRef), &configMap); err != nil {
			configMapNew = true
			configMap.Namespace = order.Status.RefObjects.DeploymentRef.Namespace
			configMap.Name = order.Status.RefObjects.DeploymentRef.Name
			if err := ctrl.SetControllerReference(&order, &configMap, r.Scheme); err != nil {
				var message = "Unable to Set Controller Reference on ConfigMap"
				log.Error(err, message)
				order.Status.Message = message
				order.Status.Created = false
				return err
			}
		}
		configMap.Data = make(map[string]string)
		configMap.Data["oauth2_proxy.cfg"] = `email_domains = [ "*" ]
upstreams = [ "file:///dev/null" ]`
		if configMapNew {
			if err := r.Create(ctx, &configMap); err != nil {
				var message = "Unable to Create ConfigMap"
				log.Error(err, message)
				order.Status.Message = message
				order.Status.Created = false
				return err
			}
			log.Info("ConfigMap Created")
		} else {
			if err := r.Update(ctx, &configMap); err != nil {
				var message = "Unable to Update ConfigMap"
				log.Error(err, message)
				order.Status.Message = message
				order.Status.Created = false
				return err
			}
			log.Info("ConfigMap Updated")
		}
		return nil
	}
	constructDeployment := func() error {
		deployment.Name = order.Status.RefObjects.DeploymentRef.Name
		deployment.Namespace = order.Status.RefObjects.DeploymentRef.Namespace
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: make(map[string]string),
		}
		deployment.Spec.Selector.MatchLabels["manager"] = "dex-configurer"
		deployment.Spec.Selector.MatchLabels["name"] = deployment.Name

		deployment.Spec.Template.Labels = make(map[string]string)
		for k, v := range deployment.Spec.Selector.MatchLabels {
			deployment.Spec.Template.Labels[k] = v
		}

		if len(deployment.Spec.Template.Spec.Containers) == 0 {
			deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, corev1.Container{})
		}
		deployment.Spec.Template.Spec.Containers[0].Args = []string{
			"--http-address=0.0.0.0:4180",
			"--metrics-address=0.0.0.0:44180",
			"--cookie-name=_oauth2_proxy",
			"--oidc-issuer-url=" + config.Spec.OcidIssuerUrl,
			"--provider=oidc",
			"--provider-display-name=" + config.Spec.ProviderDisplayName,
			"--redirect-url=" + order.Spec.RedirectUrl.GoString(),
			"--config=/etc/oauth2_proxy/oauth2_proxy.cfg",
		}
		for _, group := range order.Spec.AllowedGroups {
			deployment.Spec.Template.Spec.Containers[0].Args = append(deployment.Spec.Template.Spec.Containers[0].Args, "--allowed-group="+group)
		}
		deployment.Spec.Template.Spec.Containers[0].Args = append(deployment.Spec.Template.Spec.Containers[0].Args, order.Spec.ExtraArguments...)

		deployment.Spec.Template.Spec.Containers[0].Env = []corev1.EnvVar{
			{
				Name:  "OAUTH2_PROXY_CLIENT_ID",
				Value: order.Spec.ClientId,
			},
			{
				Name:  "OAUTH2_PROXY_CLIENT_SECRET",
				Value: order.Spec.ClientSecret,
			},
			{
				Name:  "OAUTH2_PROXY_COOKIE_SECRET",
				Value: config.Spec.CookieSecret,
			},
		}
		deployment.Spec.Template.Spec.Containers[0].Image = config.Spec.ProxyImage
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			FailureThreshold: 3,
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/ping",
					Port:   intstr.FromString("http"),
					Scheme: corev1.URISchemeHTTP,
				},
			},
			PeriodSeconds:    10,
			SuccessThreshold: 1,
			TimeoutSeconds:   1,
		}
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			FailureThreshold: 3,
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/ping",
					Port:   intstr.FromString("http"),
					Scheme: corev1.URISchemeHTTP,
				},
			},
			PeriodSeconds:    10,
			SuccessThreshold: 1,
			TimeoutSeconds:   1,
		}
		deployment.Spec.Template.Spec.Containers[0].Name = deployment.Name
		deployment.Spec.Template.Spec.Containers[0].Ports = []corev1.ContainerPort{
			{
				Name:          "http",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: 4180,
			}, {
				Name:          "metrics",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: 44180,
			},
		}
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "configmain",
				MountPath: "/etc/oauth2_proxy",
			},
		}
		deployment.Spec.Template.Spec.DNSPolicy = corev1.DNSClusterFirst
		deployment.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyAlways
		deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: "configmain",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: deployment.Name,
						},
					},
				},
			},
		}

		if err := ctrl.SetControllerReference(&order, &deployment, r.Scheme); err != nil {
			var message = "Unable to Set Controller Reference on Deployment"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return err
		}
		return nil
	}
	if err := createConfigMap(); err != nil {
		return ctrl.Result{}, err
	}
	if err := constructDeployment(); err != nil {
		return ctrl.Result{}, err
	}

	if deploymentNew {
		if err := r.Create(ctx, &deployment); err != nil {
			var message = "Unable to Create Deployment"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Deployment Created")
	} else {
		if err := r.Update(ctx, &deployment); err != nil {
			var message = "Unable to Update Deployment"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Deployment Updated")
	}

	//创建Service
	var service corev1.Service
	servicetNew := false
	if order.Status.RefObjects.ServiceRef.Namespace == "" {
		order.Status.RefObjects.ServiceRef.Namespace = order.Namespace
	}
	if order.Status.RefObjects.ServiceRef.Name == "" {
		order.Status.RefObjects.ServiceRef.Name = order.Name + "-svc"
		for {
			if err := r.Get(ctx, types.NamespacedName(order.Status.RefObjects.ServiceRef), &corev1.Service{}); err != nil {
				break
			}
			order.Status.RefObjects.ServiceRef.Name = order.Name + "-svc-" + strings.ToLower(utils.RandString(4))
		}
		service.Name = order.Status.RefObjects.ServiceRef.Name
		service.Namespace = order.Status.RefObjects.ServiceRef.Namespace
		servicetNew = true
	} else {
		if err := r.Get(ctx, types.NamespacedName(order.Status.RefObjects.ServiceRef), &service); err != nil {
			log.Error(err, "Unable to fetch Service,Will Create New One")
			servicetNew = true
		}
	}
	constructService := func() error {
		service.Name = order.Status.RefObjects.ServiceRef.Name
		service.Namespace = order.Status.RefObjects.ServiceRef.Namespace

		service.Spec.Ports = []corev1.ServicePort{
			{
				Name:       "http",
				Port:       80,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromString("http"),
			},
			{
				Name:       "metrics",
				Port:       44180,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromString("metrics"),
			},
		}
		service.Spec.Selector = make(map[string]string)
		service.Spec.Selector["manager"] = "dex-configurer"
		service.Spec.Selector["name"] = deployment.Name
		service.Spec.Type = corev1.ServiceTypeClusterIP

		if err := ctrl.SetControllerReference(&order, &service, r.Scheme); err != nil {
			var message = "Unable to Set Controller Reference on Service"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return err
		}
		return nil
	}
	if err := constructService(); err != nil {
		return ctrl.Result{}, err
	}
	if servicetNew {
		if err := r.Create(ctx, &service); err != nil {
			var message = "Unable to Create Service"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Service Created")
	} else {
		if err := r.Update(ctx, &service); err != nil {
			var message = "Unable to Update Service"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Service Updated")
	}

	//创建Ingress
	var ingress netv1.Ingress
	ingressNew := false
	if order.Status.RefObjects.IngressRef.Namespace == "" {
		order.Status.RefObjects.IngressRef.Namespace = order.Namespace
	}
	if order.Status.RefObjects.IngressRef.Name == "" {
		order.Status.RefObjects.IngressRef.Name = order.Name + "-ingress"
		for {
			if err := r.Get(ctx, types.NamespacedName(order.Status.RefObjects.IngressRef), &netv1.Ingress{}); err != nil {
				break
			}
			order.Status.RefObjects.IngressRef.Name = order.Name + "-ingress-" + strings.ToLower(utils.RandString(4))
		}
		ingress.Name = order.Status.RefObjects.IngressRef.Name
		ingress.Namespace = order.Status.RefObjects.IngressRef.Namespace
		ingressNew = true
	} else {
		if err := r.Get(ctx, types.NamespacedName(order.Status.RefObjects.IngressRef), &ingress); err != nil {
			log.Error(err, "Unable to fetch Ingress,Will Create New One")
			ingressNew = true
		}
	}
	constructIngress := func() error {
		ingress.Name = order.Status.RefObjects.IngressRef.Name
		ingress.Namespace = order.Status.RefObjects.IngressRef.Namespace

		pathType := netv1.PathTypeImplementationSpecific

		ingress.Spec.Rules = []netv1.IngressRule{
			{
				Host: order.Spec.RedirectUrl.Host,
				IngressRuleValue: netv1.IngressRuleValue{
					HTTP: &netv1.HTTPIngressRuleValue{
						Paths: []netv1.HTTPIngressPath{
							{
								Path:     order.Spec.RedirectUrl.GetPath(),
								PathType: &pathType,
								Backend: netv1.IngressBackend{
									Service: &netv1.IngressServiceBackend{
										Name: service.Name,
										Port: netv1.ServiceBackendPort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		}

		if err := ctrl.SetControllerReference(&order, &ingress, r.Scheme); err != nil {
			var message = "Unable to Set Controller Reference on Ingress"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return err
		}
		return nil
	}
	if err := constructIngress(); err != nil {
		return ctrl.Result{}, err
	}
	if ingressNew {
		if err := r.Create(ctx, &ingress); err != nil {
			var message = "Unable to Create Ingress"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Ingress Created")
	} else {
		if err := r.Update(ctx, &ingress); err != nil {
			var message = "Unable to Update Ingress"
			log.Error(err, message)
			order.Status.Message = message
			order.Status.Created = false
			return ctrl.Result{}, err
		}
		log.Info("Ingress Updated")
	}

	//处理Target Ingress的Annotations
	if targetIngress.Annotations == nil {
		targetIngress.Annotations = map[string]string{}
	}
	targetIngress.Annotations["nginx.ingress.kubernetes.io/auth-signin"] = order.Spec.RedirectUrl.AuthSignin()
	targetIngress.Annotations["nginx.ingress.kubernetes.io/auth-url"] = fmt.Sprintf("http://%v.%v.svc.cluster.local/oauth2/auth", service.Name, service.Namespace)
	if config.Spec.AuthCacheEnabled {
		targetIngress.Annotations["nginx.ingress.kubernetes.io/auth-cache-key"] = config.Spec.AuthCacheKey
	}
	if err := r.Update(ctx, &targetIngress); err != nil {
		var message = "Unable to Update Target Ingress Annotations"
		log.Error(err, message)
		order.Status.Message = message
		order.Status.Created = false
		return ctrl.Result{}, err
	}
	log.Info("Target Ingress Annotations Updated")

	order.Status.Message = ""
	order.Status.Created = true

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DexClientOrderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dexv1.DexClientOrder{}).
		Owns(&dexv1.DexAuthClient{}).
		Complete(r)
}
