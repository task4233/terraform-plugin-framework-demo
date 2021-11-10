package log

import "github.com/hashicorp/terraform-plugin-framework/types"

type Order struct {
	ID          types.String `tfsdk:"id"`
	Items       []OrderItem  `tfsdk:"items"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

type OrderItem struct {
	Log Log `tfsdk:"log"`
}

type Log struct {
	Body types.String `tfsdk:"body"`
}
