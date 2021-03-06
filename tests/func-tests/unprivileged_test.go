package tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"kubevirt.io/client-go/kubecli"
	testscore "kubevirt.io/kubevirt/tests"
)

var _ = Describe("Unprivileged tests", func() {
	virtClient, err := kubecli.GetKubevirtClient()
	testscore.PanicOnError(err)

	cfg := virtClient.Config()

	cfg.Impersonate = rest.ImpersonationConfig{
		UserName: "non-existent-user",
		Groups:   []string{"system:authenticated"},
	}

	unprivClient, err := kubecli.GetKubevirtClientFromRESTConfig(cfg)
	testscore.PanicOnError(err)

	It("should be able to read kubevirt-storage-class-defaults ConfigMap", func() {

		// Sanity check: can't read an arbitrary configmap (nonexistent)
		_, err = unprivClient.CoreV1().ConfigMaps(testscore.KubeVirtInstallNamespace).Get("non-existent-configmap", metav1.GetOptions{})
		Expect(apierrors.IsForbidden(err)).To(BeTrue())

		configmap, err := unprivClient.CoreV1().ConfigMaps(testscore.KubeVirtInstallNamespace).Get("kubevirt-storage-class-defaults", metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())

		Expect(configmap.Data["local-sc.volumeMode"]).To(Equal("Filesystem"))
		Expect(configmap.Data["local-sc.accessMode"]).To(Equal("ReadWriteOnce"))
	})
})
