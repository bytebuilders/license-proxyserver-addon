/*
Copyright AppsCode Inc. and Contributors.

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

package controller

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/fluxcd/helm-controller/api/v2beta1"
	verifier "go.bytebuilders.dev/license-verifier"
	"go.bytebuilders.dev/license-verifier/apis/licenses/v1alpha1"
	licenseclient "go.bytebuilders.dev/license-verifier/client"
	"go.bytebuilders.dev/license-verifier/info"
	v1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/clock"
	ocmv1 "open-cluster-management.io/api/cluster/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	ClusterClaimClusterID             = "id.k8s.io"
	ClusterClaimLicense               = "licenses.appscode.com"
	LicenseSecret                     = "license-proxyserver-licenses"
	LicenseProxyServerHelmReleaseName = "license-proxyserver"
	LicenseProxyServerNamespace       = "kubeops"
)

func NewLicenseReconciler(hubClient client.Client) *LicenseReconciler {
	return &LicenseReconciler{
		HubClient: hubClient,
		clock:     clock.RealClock{},
	}
}

var _ reconcile.Reconciler = &LicenseReconciler{}

type LicenseReconciler struct {
	HubClient client.Client
	clock     clock.RealClock
}

// SetupWithManager sets up the controller with the Manager.
func (r *LicenseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ocmv1.ManagedCluster{}).Watches(&ocmv1.ManagedCluster{}, &handler.EnqueueRequestForObject{}).Complete(r)
}

func (r *LicenseReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Start reconciling")
	managed := &ocmv1.ManagedCluster{}
	err := r.HubClient.Get(ctx, request.NamespacedName, managed)
	if err != nil {
		return reconcile.Result{}, err
	}

	if len(managed.Status.ClusterClaims) > 1 {
		var cid string
		var features []string

		for _, claim := range managed.Status.ClusterClaims {
			if claim.Name == ClusterClaimClusterID {
				cid = claim.Value
			}
			if claim.Name == ClusterClaimLicense {
				features = strings.Split(claim.Value, ",")
			}
		}

		err = licenseHelper(ctx, r.HubClient, cid, features, managed.Name)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func licenseHelper(ctx context.Context, HubClient client.Client, cid string, features []string, clusterName string) error {
	lps := v2beta1.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      LicenseProxyServerHelmReleaseName,
			Namespace: LicenseProxyServerNamespace,
		},
	}

	err := HubClient.Get(context.TODO(), client.ObjectKey{Name: lps.Name, Namespace: lps.Namespace}, &lps)
	if err == nil {
		baseURL, token, err := getProxyServerURLAndToken(lps)
		if err != nil {
			return err
		}

		l, err := getLicense(baseURL, token, cid, features)
		if err != nil {
			return err
		}

		// get secret
		sec := v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      LicenseSecret,
				Namespace: clusterName,
			},
		}
		err = HubClient.Get(context.TODO(), client.ObjectKey{Name: sec.Name, Namespace: sec.Namespace}, &sec)
		if err == nil {
			// update secret
			sec.Data[l.PlanName] = l.Data
			err = HubClient.Update(context.TODO(), &sec)
			if err != nil {
				return err
			}
		} else if err != nil && kerr.IsNotFound(err) {
			// create secret
			data := make(map[string][]byte)
			data[l.PlanName] = l.Data
			sec.Data = data
			err = HubClient.Create(ctx, &sec)
			if err != nil {
				return err
			} else {
				return nil
			}
		} else if err != nil {
			return err
		}

		return nil
	}

	return err
}

func getProxyServerURLAndToken(lps v2beta1.HelmRelease) (string, string, error) {
	val := make(map[string]interface{})
	jsonByte, err := json.Marshal(lps.Spec.Values)
	if err != nil {
		return "", "", err
	}
	if err = json.Unmarshal(jsonByte, &val); err != nil {
		return "", "", err
	}
	baseURL, _, err := unstructured.NestedString(val, "platform", "baseURL")
	if err != nil {
		return "", "", err
	}
	token, _, err := unstructured.NestedString(val, "platform", "token")
	if err != nil {
		return "", "", err
	}

	return baseURL, token, nil
}

func getLicense(baseURL, token, cid string, features []string) (*v1alpha1.License, error) {
	lc, err := licenseclient.NewClient(baseURL, token, cid)
	if err != nil {
		return nil, err
	}

	nl, _, err := lc.AcquireLicense(features)
	if err != nil {
		return nil, err
	}

	caData, err := info.LoadLicenseCA()
	if err != nil {
		return nil, err
	}
	caCert, err := info.ParseCertificate(caData)
	if err != nil {
		return nil, err
	}

	l, err := verifier.ParseLicense(verifier.ParserOptions{
		ClusterUID: cid,
		CACert:     caCert,
		License:    nl,
	})
	if err != nil {
		return nil, err
	}

	return &l, nil
}
