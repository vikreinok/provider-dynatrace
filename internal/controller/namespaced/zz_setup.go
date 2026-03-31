// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	group "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/group"
	policy "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/policy"
	policybindingsv2 "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/policybindingsv2"
	policyboundary "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/policyboundary"
	providerconfig "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		group.Setup,
		policy.Setup,
		policybindingsv2.Setup,
		policyboundary.Setup,
		providerconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		group.SetupGated,
		policy.SetupGated,
		policybindingsv2.SetupGated,
		policyboundary.SetupGated,
		providerconfig.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
