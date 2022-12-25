#!/bin/bash

set -ex

kubectl delete po --ignore-not-found=true fleetcommand-agent

kubectl create configmap \
  fleetcommand-agent-config \
  -o yaml --dry-run=client \
  --from-file=pab.yml | kubectl apply -f -

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin
  namespace: default
---
apiVersion: v1
kind: Pod
metadata:
  name: fleetcommand-agent
spec:
  serviceAccountName: admin
  restartPolicy: Never
  containers:
  - name: fleetcommand-agent
    image: dangawne/points-are-bad-fleetcommand-agent
    args: ["run", "-f", "/builduser/fleetcommand/pab.yml"]
    imagePullPolicy: Always
    volumeMounts:
    - name: install-config
      mountPath: /builduser/fleetcommand/
  volumes:
  - name: install-config
    configMap:
      name: fleetcommand-agent-config
EOF
