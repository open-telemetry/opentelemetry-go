package resource

import "context"

// compositeEntityDetector uses multiple sub-detectors to detect entities
// If more than one detector detects the entity of particular type all those
// detected entities will be merged together into one entity (one per entity type),
// provided that they have the same SchemaURL.
type compositeEntityDetector struct {
	detectors []Detector
}

func (p compositeEntityDetector) Detect(ctx context.Context) (*Resource, error) {
	merged := NewResource()
	for _, detector := range p.detectors {
		res, err := detector.Detect(ctx)
		if err != nil {
			return nil, err
		}
		merged, err = merge(merged, res, mergeOptions{allowMultipleOfSameType: true})
		if err != nil {
			return nil, err
		}
	}
	return merged, nil
}
