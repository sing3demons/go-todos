#!/bin/sh -e

kubectl apply -f k8s/postgres/
kubectl apply -f k8s/00-namespace.yml
kubectl apply -f k8s/redis-client.yml 
kubectl apply -f k8s/01-secret.yml 
kubectl apply -f k8s/02-pv-pvc.yml 
kubectl apply -f k8s/03-deploymet.yml
kubectl apply -f k8s/03-service.yml
kubectl apply -f k8s/04-ingress.yml