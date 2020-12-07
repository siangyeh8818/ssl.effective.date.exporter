# ssl.effective.date.exporter

Prometheus的客制exporter

功能是檢測SSL 憑證的有效天數

### How to build
```
skaffold build
```

### How to deploy
```
vi kubernetes/configmap.yaml
填入你要測的domain

vi kubernetes/deployment.yaml
填入redis的位址跟密碼

然後apply所有yaml
kubectl apply -f kubernetes/
```
