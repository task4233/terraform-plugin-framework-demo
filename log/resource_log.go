package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/task4233/note-v2-terraform/client"
)

type resourceLogType struct{}

func (r resourceLogType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
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
	// create前のStateと後のStateを一緒にしないとバグる
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Order
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	items := make([]client.OrderItem, len(plan.Items))
	{
		idx := 0
		for _, item := range plan.Items {
			if item.Log.Body.Null || item.Log.Body.Unknown {
				continue
			}
			items[idx] = client.OrderItem{
				Log: client.Log{
					Body: item.Log.Body.Value,
				},
			}
			idx++
		}
	}

	log := client.Order{
		Items: items,
	}

	gotLogs, err := r.p.client.CreateLog(ctx, &log)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating log",
			fmt.Sprintf("Could not create log: %s", err.Error()),
		)
	}

	gotItems := make([]OrderItem, len(gotLogs.Items))
	for idx, item := range gotLogs.Items {
		gotItems[idx] = OrderItem{
			Log: Log{
				Body: types.String{
					Value: item.Log.Body,
				},
			},
		}
	}

	result := Order{
		ID:          types.String{Value: "0"}, // set 0 for now
		Items:       gotItems,
		LastUpdated: types.String{Value: string(time.Now().Format(time.RFC850))},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceLog) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state Order
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	// orderID := state.ID.Value
	order, err := r.p.client.GetLogs(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed reading order",
			fmt.Sprintf("Failed client.GetLog: %s", err.Error()),
		)
	}

	state.Items = []OrderItem{}
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

	fmt.Fprintf(os.Stderr, "[Read]\n")
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

	items := make([]client.OrderItem, len(plan.Items))
	for idx := range plan.Items {
		items[idx] = client.OrderItem{
			Log: client.Log{
				Body: plan.Items[idx].Log.Body.Value,
			},
		}
	}

	orderID := state.ID.Value
	order, err := r.p.client.UpdateLog(ctx, orderID, &client.Order{
		Items: items,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update order",
			fmt.Sprintf("Could not update: %s", err.Error()),
		)
		return
	}

	lis := make([]OrderItem, len(order.Items))
	for idx := range order.Items {
		lis[idx] = OrderItem{
			Log: Log{
				Body: types.String{
					Value: order.Items[idx].Log.Body,
				},
			},
		}
	}

	var result = Order{
		ID:          types.String{Value: orderID},
		Items:       lis,
		LastUpdated: types.String{Value: string(time.Now().Format(time.RFC850))},
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Fprintf(os.Stderr, "[Update]\n")
}

func (r resourceLog) Delete(context.Context, tfsdk.DeleteResourceRequest, *tfsdk.DeleteResourceResponse) {
	fmt.Fprintf(os.Stderr, "[Delete]\n")
}

func (r resourceLog) ImportState(context.Context, tfsdk.ImportResourceStateRequest, *tfsdk.ImportResourceStateResponse) {
	fmt.Fprintf(os.Stderr, "[ImportState]\n")
}
