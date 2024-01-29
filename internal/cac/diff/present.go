package diff

import "golang.org/x/exp/maps"

func OnlyPresentKeys(source any, targetMap map[string]any) {
	var (
		sourceMap map[string]any
		ok        bool
	)
	if sourceMap, ok = source.(map[string]any); !ok {
		maps.Clear(targetMap)
	}

	for k := range targetMap {
		if _, ok := sourceMap[k]; !ok {
			delete(targetMap, k)
		}
	}
}
