//
// Copyright (c) 2019 Intel Corporation
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
//

package controller

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/edgexfoundry/device-sdk-go/internal/common"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

func init() {
	lc := logger.NewClient("update_test", false, "./device-simple.log", "DEBUG")
	common.LoggingClient = lc
}

func TestAddRoute(t *testing.T) {

	tests := []struct {
		Name          string
		Route         string
		ErrorExpected bool
	}{
		{"Success", "/api/v1/test", false},
		{"Reserved Route", common.APIVersionRoute, true},
	}

	for _, test := range tests {
		r := mux.NewRouter()
		controller := NewRestController(r)
		controller.InitRestRoutes()

		err := controller.AddRoute(test.Route, func(http.ResponseWriter, *http.Request) {}, http.MethodPost)
		if test.ErrorExpected {
			assert.Error(t, err, "Expected an error")
		} else {
			if !assert.NoError(t, err, "Unexpected an error") {
				t.Fatal()
			}

			err = controller.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
				path, err := route.GetPathTemplate()
				if err != nil {
					return err
				}

				// Have to skip all the reserved routes that have previously been added.
				if controller.reservedRoutes[path] {
					return nil
				}

				routeMethods, err := route.GetMethods()
				if err != nil {
					return err
				}

				assert.Equal(t, test.Route, path)
				assert.Equal(t, http.MethodPost, routeMethods[0], "Expected POST Method")
				return nil
			})

			assert.NoError(t, err, "Unexpected error examining route")
		}
	}
}

func TestInitRestRoutes(t *testing.T) {
	r := mux.NewRouter()
	controller := NewRestController(r)
	controller.InitRestRoutes()

	err := controller.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		// Verify the route is reserved by attempting to add it as 'external' route.
		// If tests fails then the route was not added to the reserved list
		err = controller.AddRoute(path, func(http.ResponseWriter, *http.Request) {})
		assert.Error(t, err, path, fmt.Sprintf("Expected error for '%s'", path))
		return nil
	})

	assert.NoError(t, err, "Unexpected error examining route")
}
