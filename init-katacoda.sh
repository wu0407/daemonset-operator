#!/bin/bash

wget https://golang.org/dl/go1.14.6.linux-amd64.tar.gz -O /root/go1.14.6.linux-amd64.tar.gz
wget https://github.com/operator-framework/operator-sdk/releases/download/v0.17.2/operator-sdk-v0.17.2-x86_64-linux-gnu -O /root/operator-sdk
cd root
tar xf go1.14.6.linux-amd64.tar.gz
rm -rf /usr/local/go
mv /root/go /usr/local/go

minikube start --wait=false

cd daemonset-operator/deploy

kubectl apply -f role.yaml -f service_account.yaml -f role_binding.yaml -f crds/ds.xiaoqing.com_xdaemonsets_crd.yaml

/root/operator-sdk build daemonset-operator:v0.1
sed -i 's/REPLACE_IMAGE/daemonset-operator:v0.1/' operator.yaml
kubectl apply -f operator.yaml