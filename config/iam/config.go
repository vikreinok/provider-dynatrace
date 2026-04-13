package iam

import "github.com/crossplane/upjet/v2/pkg/config"

const shortGroup = "iam"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("dynatrace_iam_group", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("dynatrace_iam_policy_boundary", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("dynatrace_iam_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("dynatrace_iam_policy_bindings_v2", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.References["group"] = config.Reference{
			Type: "Group",
		}
		r.References["policy.id"] = config.Reference{
			Type: "Policy",
		}
		r.References["policy.boundaries"] = config.Reference{
			Type: "PolicyBoundary",
		}
	})

	// New resources added from iam.tf
	p.AddResourceConfigurator("dynatrace_management_zone_v2", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ManagementZoneV2"
	})
	p.AddResourceConfigurator("dynatrace_alerting", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.References["management_zone"] = config.Reference{
			Type: "ManagementZoneV2",
		}
	})
	p.AddResourceConfigurator("dynatrace_openpipeline_v2_logs_pipelines", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("dynatrace_openpipeline_v2_spans_pipelines", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("dynatrace_openpipeline_v2_logs_pipelinegroups", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.References["member_pipelines"] = config.Reference{
			Type: "V2LogsPipelines",
		}
	})
	p.AddResourceConfigurator("dynatrace_openpipeline_v2_spans_pipelinegroups", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.References["member_pipelines"] = config.Reference{
			Type: "V2SpansPipelines",
		}
	})
	p.AddResourceConfigurator("dynatrace_openpipeline_v2_logs_routing", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.References["routing_entries.routing_entry.pipeline_id"] = config.Reference{
			Type: "V2LogsPipelines",
		}
	})
	p.AddResourceConfigurator("dynatrace_openpipeline_v2_spans_routing", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.References["routing_entries.routing_entry.pipeline_id"] = config.Reference{
			Type: "V2SpansPipelines",
		}
	})
}
