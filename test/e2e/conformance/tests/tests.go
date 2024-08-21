// Copyright (c) 2022 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/alibaba/higress/test/e2e/conformance/utils/ingress2gateway"
	"github.com/alibaba/higress/test/e2e/conformance/utils/suite"
	"github.com/stretchr/testify/require"
)

func Register(testcase suite.ConformanceTest) {
	if len(testcase.Features) == 0 {
		panic("must set at least one feature of the test case")
	}
	ConformanceTests = append(ConformanceTests, testcase)
}

var ConformanceTests []suite.ConformanceTest

func ToGatewayResource(t *testing.T, suite *suite.ConformanceTestSuite, ipath string) {
	gwResource, err := ingress2gateway.ConvertYAMLToClientObjects(ReadFromFile(ipath))
	if err != nil {
		require.NoErrorf(t, err, "error translate ingress to gatewayapi")
		t.Fatalf("failed to convert ingress to gateway resource: %v", err)
	}

	suite.Applier.MustApplyObjectsWithCleanup(t, suite.Client, suite.TimeoutConfig, gwResource, suite.Cleanup)
	t.Logf("🧳 Ingress: %s translate to Gateway resource successfully", ipath)
}

func ReadFromFile(fp string) []byte {
	basePath := "./conformance/"
	iPath := basePath + fp

	cmd := exec.Command("./conformance/utils/ingress2gateway/ingress2gateway", "print",
		"--providers", "higress", "--namespace", "higress-conformance-infra", "--input-file", iPath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}
	return output
}
