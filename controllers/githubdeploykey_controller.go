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

package controllers

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	credentialv1 "github.com/stobias123/git-credential-operator/api/v1"
	sshutil "github.com/stobias123/git-credential-operator/util"
)

const githubDeployKeyFinalizer = "finalizer.credential.github.com"

// GithubDeployKeyReconciler reconciles a GithubDeployKey object
type GithubDeployKeyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=credential.github.com,resources=githubdeploykeys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=credential.github.com,resources=githubdeploykeys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=credential.github.com,resources=githubdeploykeys/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *GithubDeployKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("githubdeploykey", req.NamespacedName)

	credsCR := &credentialv1.GithubDeployKey{}
	err := r.Get(ctx, req.NamespacedName, credsCR)
    if err != nil {
        // Error reading the object - requeue the request.
        reqLogger.Error(err, "Failed to get CredsCR.")
        return ctrl.Result{}, err
		}
		
		privateKey, publicKey, err := sshutil.GetSSHStrings()
		r.createGitHubDeployKey(ctx, *credsCR,publicKey)
		r.createGithubDeployKeySecret(ctx, *credsCR, privateKey)
		

		// ------------------------------------------------------------------------------------
		isCredentialMarkedforDeletion := credsCR.GetDeletionTimestamp() != nil
    if isCredentialMarkedforDeletion {
        if contains(credsCR.GetFinalizers(), githubDeployKeyFinalizer) {
            // Run finalization logic for githubDeployKeyFinalizer. If the
            // finalization logic fails, don't remove the finalizer so
            // that we can retry during the next reconciliation.
            if err := r.finalizeGithubDeployKey(reqLogger, credsCR); err != nil {
                return ctrl.Result{}, err
            }

            // Remove githubDeployKeyFinalizer. Once all finalizers have been
            // removed, the object will be deleted.
            controllerutil.RemoveFinalizer(credsCR, githubDeployKeyFinalizer)
            err := r.Update(ctx, credsCR)
            if err != nil {
                return ctrl.Result{}, err
            }
        }
        return ctrl.Result{}, nil
    }

    // Add finalizer for this CR
    if !contains(credsCR.GetFinalizers(), githubDeployKeyFinalizer) {
        if err := r.addFinalizer(reqLogger, credsCR); err != nil {
            return ctrl.Result{}, err
        }
    }
// ------------------------------------------------------------------------------------

	return ctrl.Result{}, nil
}

func (r *GithubDeployKeyReconciler) finalizeGithubDeployKey(reqLogger logr.Logger, credsCR *credentialv1.GithubDeployKey) error {
    // TODO(user): Add the cleanup steps that the operator
    // needs to do before the CR can be deleted. Examples
    // of finalizers include performing backups and deleting
    // resources that are not owned by this CR, like a PVC.
    reqLogger.Info("Successfully finalized GithubDeployKey")
    return nil
}

// createGitHubDeployKey creates the deploy key in github for the target repo.
func (r *GithubDeployKeyReconciler) createGitHubDeployKey(ctx context.Context, credsCR credentialv1.GithubDeployKey, publicKey []byte) error {
	token := os.Getenv("GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token },
	)

	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	repo, _, err := ghClient.Repositories.Get(ctx, credsCR.Spec.Organization, credsCR.Spec.Repo)
	fmt.Println(fmt.Sprintf("[Info] Found repo - %s", *repo.Name))

	// ---- Github
	publicKeyString := string(publicKey)
	if err != nil {
		fmt.Println("[Error] Problem creating public Key")
		return  err
	}


	//key, resp, err := ghClient.Repositories.CreateKey(ctx, repo.Owner.String(), credsCR.Spec.Repo, key)
	fmt.Println("[Info] Creating Key...")
	key, resp, err := ghClient.Repositories.CreateKey(
		ctx,
		credsCR.Spec.Organization,
		credsCR.Spec.Repo, 
		&github.Key{
			Title: &credsCR.Name,
			Key:   &publicKeyString,
		})
	fmt.Println(fmt.Sprintf("[Info] GH Response - %v", resp))
	fmt.Println(fmt.Sprintf("[Info] Created Key...%s", *key.Title))
	if err != nil {
		fmt.Println(fmt.Sprintf("[Error] Problem creating keys%s", err))
		return err
	}
	return nil
}

// createGithubDeployKeySecret generates an ssh public/private keypair and
// puts it into a kubernetes secret of type "ssh-auth" corresponding to the ssh key generated.
func (r *GithubDeployKeyReconciler) createGithubDeployKeySecret(ctx context.Context, credsCR credentialv1.GithubDeployKey, privateKey []byte) error {
	//sshKeySecret := &coreV1.Secret{}
	//var secretSpec coreV1.Secret
	secretSpec := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-githubkey", credsCR.Name),
			Annotations: map[string]string{
				"tekton.dev/git-0": "github.com",
			},
			Labels: map[string]string{
				"GithubOrganization": credsCR.Spec.Organization,
				"GithubRepo":         credsCR.Spec.Repo,
			},
			Namespace: credsCR.Namespace,
		},
		Type: "kubernetes.io/ssh-auth",
		Data: map[string][]byte{
			coreV1.SSHAuthPrivateKey: privateKey,
		},
	}

	err := r.Client.Create(ctx, secretSpec)
	if err != nil {
		fmt.Println("[Error] Error creating secret")
		return err
	}
	return nil
}

// addFinalizer adds... the finalizer...
func (r *GithubDeployKeyReconciler) addFinalizer(reqLogger logr.Logger, credsCR *credentialv1.GithubDeployKey) error {
    reqLogger.Info("Adding Finalizer for the GithubDeploy Key")
    controllerutil.AddFinalizer(credsCR, githubDeployKeyFinalizer)

    // Update CR
    err := r.Update(context.TODO(), credsCR)
    if err != nil {
        reqLogger.Error(err, "Failed to update GithubDeploy Key with finalizer")
        return err
    }
    return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GithubDeployKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&credentialv1.GithubDeployKey{}).
		Complete(r)
}

func contains(list []string, s string) bool {
    for _, v := range list {
        if v == s {
            return true
        }
    }
    return false
}