// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/peterbourgon/mergemap"
	"sigs.k8s.io/yaml"
)

var (
	_ function.Function = MergeFunction{}
)

func NewMergeFunction() function.Function {
	return MergeFunction{}
}

type MergeFunction struct{}

func (r MergeFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "deepmerge"
}

func (r MergeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Merge two maps received as arguments",
		MarkdownDescription: "Echoes given argument as result",
		Parameters: []function.Parameter{
			function.StringParameter{
				Description:         "",
				MarkdownDescription: "map to merge",
				Name:                "first",
			},
			function.StringParameter{
				Description:         "",
				MarkdownDescription: "the other map to merge",
				Name:                "second",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r MergeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var yaml1, yaml2, resultYaml string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &yaml1, &yaml2))

	var map1, map2, result map[string]any

	err := yaml.Unmarshal([]byte(yaml1), &map1)

	if err != nil {
		return
	}

	err = yaml.Unmarshal([]byte(yaml2), &map2)

	if err != nil {
		return
	}

	result = mergemap.Merge(map1, map2)
	if resp.Error != nil {
		return
	}

	resultYamlBytes, err := yaml.Marshal(result)
	if err != nil {
		return
	}

	resultYaml = string(resultYamlBytes)

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, resultYaml))
}
