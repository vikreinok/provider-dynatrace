// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	alerting "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/alerting"
	group "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/group"
	managementzonev2 "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/managementzonev2"
	policy "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/policy"
	policybindingsv2 "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/policybindingsv2"
	policyboundary "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/policyboundary"
	v2logspipelinegroups "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/v2logspipelinegroups"
	v2logspipelines "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/v2logspipelines"
	v2logsrouting "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/v2logsrouting"
	v2spanspipelinegroups "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/v2spanspipelinegroups"
	v2spanspipelines "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/v2spanspipelines"
	v2spansrouting "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/iam/v2spansrouting"
	providerconfig "github.com/vikreinok/provider-dynatrace/internal/controller/namespaced/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		alerting.Setup,
		group.Setup,
		managementzonev2.Setup,
		policy.Setup,
		policybindingsv2.Setup,
		policyboundary.Setup,
		v2logspipelinegroups.Setup,
		v2logspipelines.Setup,
		v2logsrouting.Setup,
		v2spanspipelinegroups.Setup,
		v2spanspipelines.Setup,
		v2spansrouting.Setup,
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
		alerting.SetupGated,
		group.SetupGated,
		managementzonev2.SetupGated,
		policy.SetupGated,
		policybindingsv2.SetupGated,
		policyboundary.SetupGated,
		v2logspipelinegroups.SetupGated,
		v2logspipelines.SetupGated,
		v2logsrouting.SetupGated,
		v2spanspipelinegroups.SetupGated,
		v2spanspipelines.SetupGated,
		v2spansrouting.SetupGated,
		providerconfig.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
