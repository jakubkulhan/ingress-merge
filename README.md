# Merge Ingress Controller

Merge [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) Controller combines multiple ingress 
resources into a new one.

## Motivation

Different ingress controllers [behave differently](https://github.com/kubernetes/ingress-nginx/issues/1539#issue-266008311) 
when creating _services_/_load balancers_ satisfying the ingress resources of the managed class. Some create single service/LB 
for all ingress resources, some merge resources according to hosts or TLS certificates, other create separate service/LB 
per ingress resource.

E.g. AWS ALB Ingress Controller creates a new load balancer for each ingress resource. This can become quite costly. 
There is [an issue](https://github.com/kubernetes-sigs/aws-alb-ingress-controller/issues/298) to support reusing ALBs 
across ingress resources, however, it won't be implemented anytime soon.

Merge Ingress Controller allows you to create ingress resources that will be combined together to create a new ingress
resource that will be managed by different controller.

## Setup

Install via [Helm](https://www.helm.sh/):

```sh
helm install --namespace kube-system --name ingress-merge ./helm
```

## Example

Create multiple ingresses & one config map that will provide parameters for the result ingress:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: foo-ingress
  annotations:
    kubernetes.io/ingress.class: merge
    merge.ingress.kubernetes.io/config: merged-ingress
spec:
  rules:
  - host: foo.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: foo-svc
            port: 
              number: 80

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: bar-ingress
  annotations:
    kubernetes.io/ingress.class: merge
    merge.ingress.kubernetes.io/config: merged-ingress
spec:
  rules:
  - host: bar.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: bar-svc
            port: 
              number: 80

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: merged-ingress
data:
  annotations: |
    kubernetes.io/ingress.class: other
```

Merge Ingress Controller will create new ingress resource named by the config map with rules combined together:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: merged-ingress
  annotations:
    kubernetes.io/ingress.class: other
spec:
  rules:
  - host: bar.example.com
    http:
      paths:
      - path: /
        backend:
          service:
            name: bar-svc
            port:
              number: 80

  - host: foo.example.com
    http:
      paths:
      - path: /
        backend:
          service:
            name: foo-svc
            port:
              number: 80
```

## Annotations

| Annotation | Default Value | Description | Example |
|------------|---------------|-------------|---------|
| `kubernetes.io/ingress.class` | | Use `merge` for this controller to take over. | `kubernetes.io/ingress.class: merge` | 
| `merge.ingress.kubernetes.io/config` | | Name of the [`ConfigMap`](https://kubernetes.io/docs/tutorials/configuration/) resource that will be used to merge this ingress with others. Because ingresses do not support to reference services across namespaces, neither does this reference. All ingresses to be merged, the config map & the result ingress use the same namespace. | `merge.ingress.kubernetes.io/config: merged-ingress` | 
| `merge.ingress.kubernetes.io/priority` | `0` | Rules from ingresses with higher priority come in the result ingress rules first. | `merge.ingress.kubernetes.io/priority: 10` |
| `merge.ingress.kubernetes.io/result` | | Marks ingress created by the controller. If all source ingress resources are deleted, this ingress is deleted as well. | `merge.ingress.kubernetes.io/result: "true"` |

## Configuration keys

| Key | Default Value | Description | Example |
|-----|---------------|-------------|---------|
| `name` | _name of the `ConfigMap`_ | Name of the result ingress resource. | `name: my-merged-ingress` |
| `labels` | | YAML/JSON-serialized labels to be applied to the result ingress. | `labels: '{"app": "loadbalancer", "env": "prod"}'` |
| `annotations` | `{"merge.ingress.kubernetes.io/result":"true"}` | YAML/JSON-serialized labels to be applied to the result ingress. `merge.ingress.kubernetes.io/result` with value `true` will be always added to the annotations. | `annotations: '{"kubernetes.io/ingress.class": "alb"}` |
| `backend` | | Default backend for the result ingress (`spec.backend`). Source ingresses **must not** specify default backend (such ingresses won't be merged). | `backend: '{"serviceName": "default-backend-svc", "servicePort": 80}` |

## License

Licensed under MIT license. See `LICENSE` file.
