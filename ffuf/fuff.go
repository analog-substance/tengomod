package ffuf

import (
	"context"

	ffuf "github.com/analog-substance/ffufwrap"
	"github.com/analog-substance/tengo/v2"
)

type module struct {
	ctx context.Context
}

func Module(ctx context.Context) map[string]tengo.Object {
	m := &module{
		ctx: ctx,
	}

	return map[string]tengo.Object{
		"fuzzer": &tengo.UserFunction{
			Name:  "fuzzer",
			Value: m.ffufFuzzer,
		},
		"strategy": &tengo.ImmutableMap{
			Value: map[string]tengo.Object{
				"default": &tengo.String{
					Value: string(ffuf.DefaultStrategy),
				},
				"greedy": &tengo.String{
					Value: string(ffuf.GreedyStrategy),
				},
				"basic": &tengo.String{
					Value: string(ffuf.BasicStrategy),
				},
				"advanced": &tengo.String{
					Value: string(ffuf.AdvancedStrategy),
				},
			},
		},
		"operator": &tengo.ImmutableMap{
			Value: map[string]tengo.Object{
				"or": &tengo.String{
					Value: string(ffuf.OrOperator),
				},
				"and": &tengo.String{
					Value: string(ffuf.AndOperator),
				},
			},
		},
		"mode": &tengo.ImmutableMap{
			Value: map[string]tengo.Object{
				"cluster_bomb": &tengo.String{
					Value: string(ffuf.ModeClusterBomb),
				},
				"pitch_fork": &tengo.String{
					Value: string(ffuf.ModePitchFork),
				},
				"sniper": &tengo.String{
					Value: string(ffuf.ModeSniper),
				},
			},
		},
		"format": &tengo.ImmutableMap{
			Value: map[string]tengo.Object{
				"all": &tengo.String{
					Value: string(ffuf.FormatAll),
				},
				"json": &tengo.String{
					Value: string(ffuf.FormatJSON),
				},
				"ejson": &tengo.String{
					Value: string(ffuf.FormatEJSON),
				},
				"html": &tengo.String{
					Value: string(ffuf.FormatHTML),
				},
				"md": &tengo.String{
					Value: string(ffuf.FormatMarkdown),
				},
				"csv": &tengo.String{
					Value: string(ffuf.FormatCSV),
				},
				"ecsv": &tengo.String{
					Value: string(ffuf.FormatECSV),
				},
			},
		},
	}
}

func (m *module) ffufFuzzer(args ...tengo.Object) (tengo.Object, error) {
	return newFfufFuzzer(m.ctx), nil
}
