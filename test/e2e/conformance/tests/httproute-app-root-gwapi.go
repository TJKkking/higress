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
	"github.com/alibaba/higress/test/e2e/conformance/utils/roundtripper"
	"github.com/alibaba/higress/test/e2e/conformance/utils/suite"
)

func init() {
	Register(HTTPRouteAppRootGateway)
}

var HTTPRouteAppRootGateway = suite.ConformanceTest{
	ShortName:   "HTTPRouteAppRootGateway",
	Description: "The Ingress translate to gateway resource in the higress-conformance-infra namespace uses the app root.",
	Features:    []suite.SupportedFeature{suite.Ingress2GatewayConformanceFeature},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		ToGatewayResource(t, suite, "tests/httproute-app-root.yaml")

		// https://gateway-api.sigs.k8s.io/reference/spec/#gateway.networking.k8s.io/v1.HTTPPathModifierType
		// Implementations SHOULD NOT add the port number in the ‘Location’ header in the following cases:
		// A Location header that will use HTTP (whether that is determined via the Listener protocol or the Scheme field) and use port 80.
		// A Location header that will use HTTPS (whether that is determined via the Listener protocol or the Scheme field) and use port 443.
		// but hhigress set it to 80
		testcases := []http.Assertion{
			{
				Meta: http.AssertionMeta{
					TargetBackend:   "infra-backend-v1",
					TargetNamespace: "higress-conformance-infra",
				},
				Request: http.AssertionRequest{
					ActualRequest: http.Request{
						Host:             "foo.com",
						Path:             "/",
						UnfollowRedirect: true,
					},
					RedirectRequest: &roundtripper.RedirectRequest{
						Scheme: "http",
						Host:   "foo.com",
						Path:   "/foo",
						Port:   "80",
					},
				},
				Response: http.AssertionResponse{
					ExpectedResponse: http.Response{
						StatusCode: 302,
					},
				},
			},
		}
		t.Run("HTTPRoute app root", func(t *testing.T) {
			for _, testcase := range testcases {
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, suite.GatewayAddress, testcase)
			}
		})
	},
}
