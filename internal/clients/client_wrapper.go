package clients

import (
	"context"
	"strings"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConditionFormatterClient wraps a client.Client and intercepts Status() updates
// to remove newlines from condition messages.
type ConditionFormatterClient struct {
	client.Client
}

func (c *ConditionFormatterClient) Status() client.SubResourceWriter {
	return &ConditionFormatterSubResourceWriter{
		SubResourceWriter: c.Client.Status(),
	}
}

type ConditionFormatterSubResourceWriter struct {
	client.SubResourceWriter
}

func formatConditionMessages(obj client.Object) {
	cond, ok := obj.(resource.Conditioned)
	if !ok {
		return
	}

	// We only process specific condition types used by Crossplane to avoid
	// unintentionally modifying unrelated conditions.
	conditionTypes := []xpv1.ConditionType{
		xpv1.TypeReady,
		xpv1.TypeSynced,
		xpv1.ConditionType("AsyncOperation"),
		xpv1.ConditionType("LastAsyncOperation"),
	}

	for _, ct := range conditionTypes {
		c := cond.GetCondition(ct)
		if c.Type == ct && strings.Contains(c.Message, "\n") {
			// Replace newlines with spaces to keep it on a single line
			c.Message = strings.ReplaceAll(c.Message, "\n", " | ")
			cond.SetConditions(c)
		}
	}
}

func (sw *ConditionFormatterSubResourceWriter) Update(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
	formatConditionMessages(obj)
	return sw.SubResourceWriter.Update(ctx, obj, opts...)
}

func (sw *ConditionFormatterSubResourceWriter) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
	formatConditionMessages(obj)
	return sw.SubResourceWriter.Patch(ctx, obj, patch, opts...)
}

func (sw *ConditionFormatterSubResourceWriter) Create(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceCreateOption) error {
	// Usually conditions aren't set during Create, but we format just in case.
	formatConditionMessages(obj)
	return sw.SubResourceWriter.Create(ctx, obj, subResource, opts...)
}

// Ensure the wrapper satisfies the SubResourceWriter interface (some methods might require runtime pkg)
// If Apply is called, it might not be a client.Object, but we just delegate it.
