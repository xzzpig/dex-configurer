# dex-configurer
An app to auto config dex auth client

## Feature
- `DexAuthClient`: Auto config DEX config file
- `DexClientOrder`: Auto create an oauth2 proxy and apply to the target ingress
- `Ingress`: Auto create a `DexClientOrder` by adding an annotation to the ingress,the annotation key is `dex.xzzpig.com/oauth-proxy` and value is the name of `DexProxyConfig`

## Install & Uninstall

### Install via Kubectl
```bash
kubectl apply -f https://github.com/xzzpig/dex-configurer/raw/master/deploy/bundle.yaml
```

### Unstall via kubectl
```bash
kubectl delete --all DexAuthClient && kubectl delete -f https://github.com/xzzpig/dex-configurer/raw/master/deploy/bundle.yaml
```

# Cookbook

## Manage Dex Config File by `DexAuthClient`
```yaml
apiVersion: dex.xzzpig.com/v1
kind: DexAuthClient
metadata:
  name: dexauthclient-sample
spec:
  id: example-app
  name: Example App
  redirectURIs:
  - http://127.0.0.1:5555/callback
  secret: abcdabcd
  secretRef:
    namespace: ""
    name: dex
    key: config.yaml
```

## Create OAuth Client by `DexClientOrder`
1. Create `DexProxyConfig`
```yaml
apiVersion: dex.xzzpig.com/v1
kind: DexProxyConfig
metadata:
  name: dexproxyconfig-sample
spec:
  ocid-issuer-url: https://auth.sample.com
  provider-display-name: Sample OCID Provider
  secretRef:
    namespace: ""
    name: dex
    key: config.yaml
  default-url:
    scheme: https
```
2. Create `DexClientOrder`
```yaml
apiVersion: dex.xzzpig.com/v1
kind: DexClientOrder
metadata:
  name: dexclientorder-sample
spec:
  config:
    name: dexproxyconfig-sample
  target-ingress:
    name: dexproxy-sample-ingress
```

3. Enjoy!
> App will Create an OAuth Proxy and Apply to Target Ingress

## Ingress Annotation Managed Order
1. Create `DexProxyConfig`
> Same To Before

2. Create/Modify Ingress
> Add an annotation to the `Ingress`,the annotation key is `dex.xzzpig.com/oauth-proxy` and value is the name of `DexProxyConfig`
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example
  annotations:
    dex.xzzpig.com/oauth-proxy`: dexclientorder-sample # Look Here
spec:
  rules:
  - host: example.yourdomain.com
    http:
      paths:
      - backend:
          service:
            name: example-service
            port:
              number: 5000
        path: /
        pathType: Prefix
```
3. Enjoy!
> App will Create a DexClientOrder
