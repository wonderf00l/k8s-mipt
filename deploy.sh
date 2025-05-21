docker build -t my-app .
minikube image load my-app
cd helm
helm delete my-app || true
helm install my-app .