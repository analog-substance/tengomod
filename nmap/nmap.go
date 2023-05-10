package nmap

import "github.com/analog-substance/tengo/v2"

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"scanner":           &tengo.UserFunction{Name: "scanner", Value: nmapScanner},
		"timing_slowest":    &tengo.Int{Value: 0},
		"timing_sneaky":     &tengo.Int{Value: 1},
		"timing_polite":     &tengo.Int{Value: 2},
		"timing_normal":     &tengo.Int{Value: 3},
		"timing_aggressive": &tengo.Int{Value: 4},
		"timing_fastest":    &tengo.Int{Value: 5},
	}
}

// nmapScanner creates a new NmapScanner
// Represents 'nmap.scanner() NmapScanner'
func nmapScanner(args ...tengo.Object) (tengo.Object, error) {
	scanner, err := makeNmapScanner()
	if err != nil {
		return nil, err
	}

	return scanner, nil
}
