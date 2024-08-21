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
	"testing"

	"github.com/alibaba/higress/test/e2e/conformance/utils/http"
	"github.com/alibaba/higress/test/e2e/conformance/utils/suite"
)

func init() {
	Register(HTTPRouteRewriteHostGateway)
}

var HTTPRouteRewriteHostGateway = suite.ConformanceTest{
	ShortName:   "HTTPRouteRewriteHostGateway",
	Description: "A single Ingress in the higress-conformance-infra namespace uses the rewrite host.",
	Features:    []suite.SupportedFeature{suite.Ingress2GatewayConformanceFeature},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		ToGatewayResource(t, suite, "tests/httproute-rewrite-host.yaml")

		testcases := []http.Assertion{
			{
				Request: http.AssertionRequest{
					ActualRequest: http.Request{
						Path: "/foo",
						Host: "foo.com",
					},
					ExpectedRequest: &http.ExpectedRequest{
						Request: http.Request{
							Host: "higress.io",
							Path: "/foo",
						},
					},
				},

				Response: http.AssertionResponse{
					ExpectedResponse: http.Response{
						StatusCode: 200,
					},
				},

				Meta: http.AssertionMeta{
					TargetBackend:   "infra-backend-v1",
					TargetNamespace: "higress-conformance-infra",
				},
			}, {
				Request: http.AssertionRequest{
					ActualRequest: http.Request{
						Path: "/foo",
						Host: "bar.com",
					},
					ExpectedRequest: &http.ExpectedRequest{
						Request: http.Request{
							Host: "higress.io",
							Path: "/foo",
						},
					},
				},

				Response: http.AssertionResponse{
					ExpectedResponse: http.Response{
						StatusCode: 200,
					},
				},

				Meta: http.AssertionMeta{
					TargetBackend:   "infra-backend-v2",
					TargetNamespace: "higress-conformance-infra",
				},
			},
		}

		t.Run("Rewrite HTTPRoute Host", func(t *testing.T) {
			for _, testcase := range testcases {
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, suite.GatewayAddress, testcase)
			}
		})
	},
}
