// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gcp

const (
	// See https://cloud.google.com/appengine/docs/flexible/python/migrating#modules
	// for the environment variables available in GAE environments
	gaeServiceEnv  = "GAE_SERVICE"
	gaeVersionEnv  = "GAE_VERSION"
	gaeInstanceEnv = "GAE_INSTANCE"
)

func (d *Detector) onAppEngine() bool {
	_, found := d.os.LookupEnv(gaeServiceEnv)
	return found
}

// AppEngineServiceName returns the service name of the app engine service.
func (d *Detector) AppEngineServiceName() (string, error) {
	if name, found := d.os.LookupEnv(gaeServiceEnv); found {
		return name, nil
	}
	return "", errEnvVarNotFound
}

// AppEngineServiceVersion returns the service version of the app engine service.
func (d *Detector) AppEngineServiceVersion() (string, error) {
	if version, found := d.os.LookupEnv(gaeVersionEnv); found {
		return version, nil
	}
	return "", errEnvVarNotFound
}

// AppEngineServiceInstance returns the service instance of the app engine service.
func (d *Detector) AppEngineServiceInstance() (string, error) {
	if instanceID, found := d.os.LookupEnv(gaeInstanceEnv); found {
		return instanceID, nil
	}
	return "", errEnvVarNotFound
}
