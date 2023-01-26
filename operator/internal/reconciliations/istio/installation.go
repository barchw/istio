package istio

import (
	"context"

	operatorv1alpha1 "github.com/kyma-project/istio/operator/api/v1alpha1"
	"github.com/kyma-project/istio/operator/pkg/lib/gatherer"
	"github.com/masterminds/semver"
	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Installation struct {
	Client         IstioClient
	IstioVersion   string
	IstioImageBase string
}

// Reconcile setup configuration and runs an Istio installation with merged Istio Operator manifest file.
func (i *Installation) Reconcile(ctx context.Context, istioCR *operatorv1alpha1.Istio, kubeClient client.Client) (ctrl.Result, error) {
	installedVersions, err := gatherer.ListInstalledIstioRevisions(ctx, kubeClient)
	if err != nil {
		return ctrl.Result{}, err
	}
	if len(installedVersions) > 0 {
		// compare versions and make a default revision
		if semver.MustParse(i.IstioVersion).LessThan(installedVersions["default"]) {
			return ctrl.Result{}, nil
		}
	}

	mergedIstioOperatorPath, err := merge(istioCR, i.Client.defaultIstioOperatorPath, i.Client.workingDir, TemplateData{IstioVersion: i.IstioVersion, IstioImageBase: i.IstioImageBase})
	if err != nil {
		return ctrl.Result{}, err
	}

	err = i.Client.Install(mergedIstioOperatorPath)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func isIstioInstalled(kubeClient client.Client) bool {
	var istiodList appsv1.DeploymentList
	err := kubeClient.List(context.Background(), &istiodList, client.MatchingLabels(gatherer.IstiodAppLabel))
	if err != nil {
		return false
	}

	return len(istiodList.Items) > 0
}
