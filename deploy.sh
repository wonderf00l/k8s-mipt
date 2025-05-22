docker build -t my-app .

# minikube start
minikube image load my-app

# istioctl install --set profile=demo -y
kubectl label namespace default istio-injection=enabled

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install prometheus prometheus-community/kube-prometheus-stack

cd helm
helm delete my-app || true
helm install my-app .