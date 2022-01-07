//go:build integration
// +build integration

// To enable compilation of this file in Goland, go to "Settings -> Go -> Vendoring & Build Tags -> Custom Tags" and add "integration"

/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service_binding

import (
	"fmt"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"

	. "github.com/apache/camel-k/e2e/support"
	v1 "github.com/apache/camel-k/pkg/apis/camel/v1"
)

func TestImageRegistryIsAMavenRepository(t *testing.T) {
	WithNewTestNamespace(t, func(ns string) {
		Expect(Kamel("install", "-n", ns).Execute()).To(Succeed())

		// Create integration that should decrypt foobar and log it
		name := "foobar-decryption"
		jar, err := filepath.Abs("files/sample-decryption-1.0.jar")
		assert.Nil(t, err)
		pom, err := filepath.Abs("files/sample-decryption-1.0.pom")
		assert.Nil(t, err)

		Expect(Kamel("run", "files/FoobarDecryption.java",
			"--name", name,
			"-d", fmt.Sprintf("file://%s", jar),
			"-d", fmt.Sprintf("file://%s", pom),
			"-n", ns,
		).Execute()).To(Succeed())

		Eventually(IntegrationPodPhase(ns, name), TestTimeoutMedium).Should(Equal(corev1.PodRunning))
		Eventually(IntegrationCondition(ns, name, v1.IntegrationConditionReady), TestTimeoutShort).Should(Equal(corev1.ConditionTrue))
		Eventually(IntegrationLogs(ns, name), TestTimeoutShort).Should(ContainSubstring("foobar"))

		// Clean up
		Expect(Kamel("delete", "--all", "-n", ns).Execute()).To(Succeed())
	})
}

func TestLocalFilesAreMountedInContainerInDefaultPath(t *testing.T) {
	WithNewTestNamespace(t, func(ns string) {
		// Expect(Kamel("install", "-n", ns).Execute()).To(Succeed())
		name := "laughing-route-default-path"

		Expect(Kamel("run", "files/LaughingRoute.java",
			"--name", name,
			"-p", "location=.?filename=laugh.txt",
			"-d", "file://files/laugh.txt",
			"-n", "camel-k",
		).Execute()).To(Succeed())

		Eventually(IntegrationPodPhase(ns, name), TestTimeoutMedium).Should(Equal(corev1.PodRunning))
		Eventually(IntegrationCondition(ns, name, v1.IntegrationConditionReady), TestTimeoutShort).Should(Equal(corev1.ConditionTrue))
		Eventually(IntegrationLogs(ns, name), TestTimeoutShort).Should(ContainSubstring("haha"))

		// Clean up
		Expect(Kamel("delete", "--all", "-n", ns).Execute()).To(Succeed())
	})
}

func TestLocalFilesAreMountedInContainerInCustomPath(t *testing.T) {
	WithNewTestNamespace(t, func(ns string) {
		Expect(Kamel("install", "-n", ns).Execute()).To(Succeed())
		name := "laughing-route-custom-path"
		customPath := "this/is/a/custom/path/"

		Expect(Kamel("run", "files/LaughingRoute.java",
			"--name", name,
			"-p", fmt.Sprintf("location=%s", customPath),
			"-d", fmt.Sprintf("file://files/laugh.txt:%slaugh.txt", customPath),
			"-n", ns,
		).Execute()).To(Succeed())

		Eventually(IntegrationPodPhase(ns, name), TestTimeoutMedium).Should(Equal(corev1.PodRunning))
		Eventually(IntegrationCondition(ns, name, v1.IntegrationConditionReady), TestTimeoutShort).Should(Equal(corev1.ConditionTrue))
		Eventually(IntegrationLogs(ns, name), TestTimeoutShort).Should(ContainSubstring("haha"))

		// Clean up
		Expect(Kamel("delete", "--all", "-n", ns).Execute()).To(Succeed())
	})
}

func TestLocalDirectoryIsMountedInContainer(t *testing.T) {
	WithNewTestNamespace(t, func(ns string) {
		Expect(Kamel("install", "-n", ns).Execute()).To(Succeed())
		name := "laughing-route-directory"

		Expect(Kamel("run", "files/LaughingRoute.java",
			"--name", name,
			"-p", "location=files/",
			"-d", fmt.Sprintf("file://files/laughs/:files/"),
			"-n", ns,
		).Execute()).To(Succeed())

		Eventually(IntegrationPodPhase(ns, name), TestTimeoutMedium).Should(Equal(corev1.PodRunning))
		Eventually(IntegrationCondition(ns, name, v1.IntegrationConditionReady), TestTimeoutShort).Should(Equal(corev1.ConditionTrue))
		Eventually(IntegrationLogs(ns, name), TestTimeoutShort).Should(ContainSubstring("haha"))
		Eventually(IntegrationLogs(ns, name), TestTimeoutShort).Should(ContainSubstring("hehe"))

		// Clean up
		Expect(Kamel("delete", "--all", "-n", ns).Execute()).To(Succeed())
	})
}
