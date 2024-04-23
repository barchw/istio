package ingressgateway_test

import (
	"context"
	"fmt"
	"github.com/kyma-project/istio/operator/internal/reconciliations/istio"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1alpha2 "github.com/kyma-project/istio/operator/api/v1alpha2"
	ingressgateway "github.com/kyma-project/istio/operator/pkg/lib/ingress_gateway"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"
)

var _ = Describe("Ingress Gateway Restarter", func() {
	Context("RequiresIngressGatewayRestart", func() {
		It("Should evaluate to true if new is nil and old is nil", func() {
			evaluator := ingressgateway.IngressGatewayRestartEvaluator{
				NewNumTrustedProxies: nil,
				OldNumTrustedProxies: new(int),
			}

			Expect(evaluator.RequiresIngressGatewayRestart()).To(BeTrue())
		})

		It("Should evaluate to true if new is not nil and old is nil", func() {
			evaluator := ingressgateway.IngressGatewayRestartEvaluator{
				NewNumTrustedProxies: new(int),
				OldNumTrustedProxies: nil,
			}

			Expect(evaluator.RequiresIngressGatewayRestart()).To(BeTrue())
		})

		It("Should evaluate to true if numTrustedProxies is different", func() {
			newNumTrustedProxies := 1
			oldNumTrustedProxies := 2

			evaluator := ingressgateway.IngressGatewayRestartEvaluator{
				NewNumTrustedProxies: &newNumTrustedProxies,
				OldNumTrustedProxies: &oldNumTrustedProxies,
			}

			Expect(evaluator.RequiresIngressGatewayRestart()).To(BeTrue())
		})

		It("Should evaluate to false if numTrustedProxies is the same", func() {
			numTrustedProxies := 1
			oldNumTrustedProxies := 1

			evaluator := ingressgateway.IngressGatewayRestartEvaluator{
				NewNumTrustedProxies: &oldNumTrustedProxies,
				OldNumTrustedProxies: &numTrustedProxies,
			}

			Expect(evaluator.RequiresIngressGatewayRestart()).To(BeFalse())

		})
	})

	Context("NewIngressGatewayEvaluator", func() {
		const (
			mockIstioTag             string = "1.16.1-distroless"
			lastAppliedConfiguration        = istio.LastAppliedConfiguration
		)

		It("Should return an error if getLastAppliedConfiguration fails", func() {
			predicate := ingressgateway.NewIngressGatewayRestartPredicate(&operatorv1alpha2.Istio{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						lastAppliedConfiguration: `{"config":{"numTrustedProxies":abc},"IstioTag":w}`,
					},
				},
			})
			_, err := predicate.NewIngressGatewayEvaluator(context.Background())

			Expect(err).To(HaveOccurred())
		})

		It("Should return nil for old numTrustedProxies if lastAppliedConfiguration is empty", func() {
			predicate := ingressgateway.NewIngressGatewayRestartPredicate(&operatorv1alpha2.Istio{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{},
				},
			})
			evaluator, err := predicate.NewIngressGatewayEvaluator(context.Background())

			Expect(err).NotTo(HaveOccurred())
			Expect(evaluator).NotTo(BeNil())
			Expect(evaluator.(ingressgateway.IngressGatewayRestartEvaluator).OldNumTrustedProxies).To(BeNil())
		})

		It("Should return correct not nil value for new and old numTrustedProxies", func() {
			istio := &operatorv1alpha2.Istio{
				Spec: operatorv1alpha2.IstioSpec{
					Config: operatorv1alpha2.Config{
						NumTrustedProxies: ptr.To(1),
					},
				},
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{lastAppliedConfiguration: fmt.Sprintf(`{"config":{"numTrustedProxies":1},"IstioTag":"%s"}`, mockIstioTag)},
				},
			}

			predicate := ingressgateway.NewIngressGatewayRestartPredicate(istio)
			evaluator, err := predicate.NewIngressGatewayEvaluator(context.Background())

			Expect(err).NotTo(HaveOccurred())
			Expect(evaluator).NotTo(BeNil())
			Expect(*(evaluator.(ingressgateway.IngressGatewayRestartEvaluator).NewNumTrustedProxies)).To(Equal(1))
			Expect(*(evaluator.(ingressgateway.IngressGatewayRestartEvaluator).OldNumTrustedProxies)).To(Equal(1))
		})

		It("Should return correct nil value for new and old numTrustedProxies", func() {
			istio := &operatorv1alpha2.Istio{
				Spec: operatorv1alpha2.IstioSpec{
					Config: operatorv1alpha2.Config{
						NumTrustedProxies: nil,
					},
				},
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{lastAppliedConfiguration: fmt.Sprintf(`{"IstioTag":"%s"}`, mockIstioTag)},
				},
			}

			predicate := ingressgateway.NewIngressGatewayRestartPredicate(istio)
			evaluator, err := predicate.NewIngressGatewayEvaluator(context.Background())

			Expect(err).NotTo(HaveOccurred())
			Expect(evaluator).NotTo(BeNil())
			Expect(evaluator.(ingressgateway.IngressGatewayRestartEvaluator).NewNumTrustedProxies).To(BeNil())
			Expect(evaluator.(ingressgateway.IngressGatewayRestartEvaluator).OldNumTrustedProxies).To(BeNil())
		})

	})
})
