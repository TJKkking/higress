// Copyright (c) 2022 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"testing"

	"github.com/alibaba/higress/test/e2e/conformance/utils/http"
	"github.com/alibaba/higress/test/e2e/conformance/utils/roundtripper"
	"github.com/alibaba/higress/test/e2e/conformance/utils/suite"
)

func init() {
	Register(HttpForceRedirectHttpsGateway)
}

var HttpForceRedirectHttpsGateway = suite.ConformanceTest{
	ShortName:   "HttpForceRedirectHttpsGateway",
	Description: " The ingress in the higress-conformance-infra namespace enforces server-side HTTPS with forced redirection.",
	// todo add the feature
	Features: []suite.SupportedFeature{suite.HTTPConformanceFeature},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		ToGatewayResource(t, suite, "tests/httproute-force-redirect-https.yaml")

		testcases := []http.Assertion{
			{
				Meta: http.AssertionMeta{
					TargetBackend:   "infra-backend-v1",
					TargetNamespace: "higress-conformance-infra",
				},
				Request: http.AssertionRequest{
					ActualRequest: http.Request{
						Host:             "test.com",
						Path:             "/test",
						UnfollowRedirect: true,
					},
					RedirectRequest: &roundtripper.RedirectRequest{
						Scheme: "https",
						Host:   "test.com",
						Path:   "/test",
					},
				},
				Response: http.AssertionResponse{
					ExpectedResponse: http.Response{
						StatusCode: 308,
					},
				},
			},
		}

		t.Run("HttpForceRedirectHttps", func(t *testing.T) {
			for _, testcase := range testcases {
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, suite.GatewayAddress, testcase)
			}
		})

	},
}
