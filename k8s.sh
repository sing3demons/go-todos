#!/bin/sh -e

# Docker Desktop Mac, from version 18.06.0-ce
# kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.2.0/deploy/static/provider/aws/deploy.yaml

kubectl apply -f k8s/postgres/
kubectl apply -f k8s/00-namespace.yml
kubectl apply -f k8s/redis-client.yml 
kubectl apply -f k8s/01-secret.yml 
kubectl apply -f k8s/02-pv-pvc.yml 
kubectl apply -f k8s/03-deploymet.yml
kubectl apply -f k8s/03-service.yml
kubectl apply -f k8s/04-ingress.yml