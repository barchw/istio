package istio

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	operatorv1alpha1 "github.com/kyma-project/istio/operator/api/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Installation struct {
	Client         IstioClient
	IstioVersion   string
	IstioImageBase string
}

const (
	LastAppliedConfiguration string = "operator.kyma-project.io/lastAppliedConfiguration"
)

// PerformInstall runs Istio installation with merged Istio Operator manifest file when the trigger requires an installation.
func (i *Installation) PerformInstall(ctx context.Context, trigger IstioCRChange, istioCR *operatorv1alpha1.Istio, kubeClient client.Client) error {
	if !trigger.NeedsIstioInstall() {
		ctrl.Log.Info("Install of Istio was skipped")
		return nil
	}

	mergedIstioOperatorPath, err := merge(istioCR, i.Client.defaultIstioOperatorPath, i.Client.workingDir, TemplateData{IstioVersion: i.IstioVersion, IstioImageBase: i.IstioImageBase})
	if err != nil {
		return err
	}

	err = i.Client.Install(mergedIstioOperatorPath)
	if err != nil {
		return err
	}

	return nil
}
