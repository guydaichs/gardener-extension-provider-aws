// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package validation

import (
	"fmt"

	apisaws "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws"

	"github.com/gardener/gardener/pkg/apis/core"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateNetworking validates the network settings of a Shoot.
func ValidateNetworking(networking core.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if networking.Nodes == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("nodes"), "a nodes CIDR must be provided for AWS shoots"))
	}

	return allErrs
}

// ValidateWorkers validates the workers of a Shoot.
func ValidateWorkers(workers []core.Worker, zones []apisaws.Zone, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	awsZones := sets.NewString()
	for _, awsZone := range zones {
		awsZones.Insert(awsZone.Name)
	}

	for i, worker := range workers {
		if worker.Volume == nil {
			allErrs = append(allErrs, field.Required(fldPath.Index(i).Child("volume"), "must not be nil"))
		} else {
			allErrs = append(allErrs, validateVolume(worker.Volume, fldPath.Index(i).Child("volume"))...)
		}

		if len(worker.Zones) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Index(i).Child("zones"), "at least one zone must be configured"))
			continue
		}

		allErrs = append(allErrs, validateZones(worker.Zones, awsZones, fldPath.Index(i).Child("zones"))...)
	}

	return allErrs
}

func validateZones(zones []string, allowedZones sets.String, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, workerZone := range zones {
		if !allowedZones.Has(workerZone) {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), workerZone, fmt.Sprintf("supported values %v", allowedZones.UnsortedList())))
		}
	}
	return allErrs
}

func validateVolume(vol *core.Volume, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if vol.Type == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("type"), "must not be empty"))
	}
	if vol.Size == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("size"), "must not be empty"))
	}
	return allErrs
}
