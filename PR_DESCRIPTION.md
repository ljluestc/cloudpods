# Feature: Support Annotations in DeploymentCreateOptions

## Description
This PR addresses issue #21376 by adding support for passing annotations when creating deployments via `DeploymentCreateOptions`. This enhancement enables advanced configurations, such as specifying GPU requirements for self-built Kubernetes clusters using annotations (e.g., `aliyun.com/gpu-mem`).

## Changes

### pkg/mcclient/options/k8s

#### [pod_template.go](file:///home/calelin/dev/cloudpods/pkg/mcclient/options/k8s/pod_template.go)
- Introduced `K8sAnnotationOptions` struct to handle parsing of `key=value` annotation strings.
- Implemented `Params()` and `Attach()` methods for `K8sAnnotationOptions` to correctly format and attach annotations to the request parameters.

#### [deployment.go](file:///home/calelin/dev/cloudpods/pkg/mcclient/options/k8s/deployment.go)
- Embedded `K8sAnnotationOptions` into the `DeploymentCreateOptions` struct.
- Updated the `Params()` method of `DeploymentCreateOptions` to invoke `K8sAnnotationOptions.Attach()`, ensuring annotations are included in the API request.

## Verification
- **Automated Tests**: Verified with a local unit test `TestDeploymentCreateOptions_Params` which confirmed that annotations are correctly parsed and added to the JSON parameters under the `annotations` key.
- **Manual Verification**: Since this is a client-side option change, verify that the generated JSON payload now includes the `annotations` field when `Annotations` are provided in the options.

## Related Issues
- Fixes #21376
