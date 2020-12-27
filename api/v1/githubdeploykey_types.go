/*
Copyright 2020.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GithubDeployKeySpec defines the desired state of GithubDeployKey
type GithubDeployKeySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Repo is an example field of GithubDeployKey. Edit GithubDeployKey_types.go to remove/update
	Repo         string `json:"repo"`
	Organization string `json:"organization"`
}

// GithubDeployKeyStatus defines the observed state of GithubDeployKey
type GithubDeployKeyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	KeyCreated bool `json:"key_created"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GithubDeployKey is the Schema for the githubdeploykeys API
type GithubDeployKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GithubDeployKeySpec   `json:"spec,omitempty"`
	Status GithubDeployKeyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GithubDeployKeyList contains a list of GithubDeployKey
type GithubDeployKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GithubDeployKey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GithubDeployKey{}, &GithubDeployKeyList{})
}
