# operator scope

The daemonset-operator   default is namespaced that it only controll resource in one namespace.

It can change to cluster scoped to manage all  namespaces resource.

## change to cluster scoped 

- deploy/operator.yaml
  - Set `WATCH_NAMESPACE=""` to watch all namespaces instead of setting it to the pod’s namespace
  - Set `metadata.namespace` to define the namespace where the operator will be deployed.
- deploy/role.yaml
  - Use `ClusterRole` instead of `Role`
- deploy/role_binding.yaml
  - Use `ClusterRoleBinding` instead of `RoleBinding`
  - Use `ClusterRole` instead of `Role` for `roleRef`
  - Set the subject namespace to the namespace in which the operator is deployed.
- deploy/service_account.yaml
  - Set `metadata.namespace` to the namespace where the operator is deployed.