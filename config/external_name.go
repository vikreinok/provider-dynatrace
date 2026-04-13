package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	"dynatrace_iam_group":                            config.IdentifierFromProvider,
	"dynatrace_iam_policy_boundary":                  config.IdentifierFromProvider,
	"dynatrace_iam_policy":                           config.IdentifierFromProvider,
	"dynatrace_iam_policy_bindings_v2":               config.IdentifierFromProvider,
	"dynatrace_management_zone_v2":                   config.IdentifierFromProvider,
	"dynatrace_alerting":                             config.IdentifierFromProvider,
	"dynatrace_openpipeline_v2_logs_pipelines":       config.IdentifierFromProvider,
	"dynatrace_openpipeline_v2_spans_pipelines":      config.IdentifierFromProvider,
	"dynatrace_openpipeline_v2_logs_pipelinegroups":  config.IdentifierFromProvider,
	"dynatrace_openpipeline_v2_spans_pipelinegroups": config.IdentifierFromProvider,
	"dynatrace_openpipeline_v2_logs_routing":         config.IdentifierFromProvider,
	"dynatrace_openpipeline_v2_spans_routing":        config.IdentifierFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
