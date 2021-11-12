package log

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/task4233/terraform-plugin-framework-demo/client"
)

type resourceLogType struct{}

func (r resourceLogType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"items": {
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"log": {
						Required: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"body": {
								Type:     types.StringType,
								Required: true,
							},
						}),
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"last_updated": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

type resourceLog struct {
	p provider
}

func (r resourceLogType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceLog{
		p: *(p.(*provider)),
	}, nil
}

func (r resourceLog) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan Order
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	items := make([]client.OrderItem, 0, len(plan.Items))
	for _, item := range plan.Items {
		if item.Log.Body.Null || item.Log.Body.Unknown {
			continue
		}
		items = append(items, client.OrderItem{
			Log: client.Log{
				Body: item.Log.Body.Value,
			},
		})
	}

	gotLogs, err := r.p.client.CreateLog(ctx, &client.Order{Items: items})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating log",
			fmt.Sprintf("Could not create log: %s", err.Error()),
		)
		return
	}

	gotItems := make([]OrderItem, 0, len(gotLogs.Items))
	for _, item := range gotLogs.Items {
		gotItems = append(gotItems, OrderItem{
			Log: Log{
				Body: types.String{
					Value: item.Log.Body,
				},
			},
		})
	}

	diags = resp.State.Set(ctx, Order{
		Items:       gotItems,
		LastUpdated: types.String{Value: string(time.Now().Format(time.RFC850))},
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceLog) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state Order
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	order, err := r.p.client.GetLogs(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed reading order",
			fmt.Sprintf("Failed client.GetLog: %s", err.Error()),
		)
		return
	}

	state.Items = make([]OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		state.Items = append(state.Items, OrderItem{
			Log: Log{
				Body: types.String{
					Value: item.Log.Body,
				},
			},
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceLog) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan Order
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state Order
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	items := make([]client.OrderItem, 0, len(plan.Items))
	for _, item := range plan.Items {
		items = append(items, client.OrderItem{
			Log: client.Log{
				Body: item.Log.Body.Value,
			},
		})
	}

	updatedOrder, err := r.p.client.UpdateLog(ctx, &client.Order{Items: items})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update order",
			fmt.Sprintf("Could not update: %s", err.Error()),
		)
		return
	}

	gotItems := make([]OrderItem, 0, len(updatedOrder.Items))
	for _, item := range updatedOrder.Items {
		gotItems = append(gotItems, OrderItem{
			Log: Log{
				Body: types.String{
					Value: item.Log.Body,
				},
			},
		})
	}

	diags = resp.State.Set(ctx, Order{
		Items:       gotItems,
		LastUpdated: types.String{Value: string(time.Now().Format(time.RFC850))},
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceLog) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Order
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.p.client.DeleteLog(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error delete order",
			fmt.Sprintf("Could not delete: %s", err.Error()),
		)
		return
	}

	resp.State.RemoveResource(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceLog) ImportState(context.Context, tfsdk.ImportResourceStateRequest, *tfsdk.ImportResourceStateResponse) {
	// TODO: no need to implement for now
	panic("not implemented")
}
