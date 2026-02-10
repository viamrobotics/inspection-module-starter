package inspectionmodule

import (
	"context"
	"fmt"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	generic "go.viam.com/rdk/services/generic"
	"go.viam.com/rdk/services/vision"
)

var Inspector = resource.NewModel("stations", "inspection-module", "inspector")

func init() {
	resource.RegisterService(generic.API, Inspector,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newInspectionModuleInspector,
		},
	)
}

type Config struct {
	Camera        string `json:"camera"`
	VisionService string `json:"vision"`
}

func (cfg *Config) Validate(path string) ([]string, []string, error) {
	if cfg.Camera == "" {
		return nil, nil, fmt.Errorf("camera is required")
	}
	if cfg.VisionService == "" {
		return nil, nil, fmt.Errorf("vision is required")
	}
	return []string{cfg.Camera, cfg.VisionService}, nil, nil
}

type inspectionModuleInspector struct {
	resource.AlwaysRebuild

	name   resource.Name
	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	detector vision.Service
}

func newInspectionModuleInspector(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}
	return NewInspector(ctx, deps, rawConf.ResourceName(), conf, logger)
}

func NewInspector(ctx context.Context, deps resource.Dependencies, name resource.Name, cfg *Config, logger logging.Logger) (resource.Resource, error) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	detector, err := vision.FromProvider(deps, cfg.VisionService)
	if err != nil {
		return nil, fmt.Errorf("failed to get vision service %q: %w", cfg.VisionService, err)
	}

	s := &inspectionModuleInspector{
		name:       name,
		logger:     logger,
		cfg:        cfg,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		detector:   detector,
	}
	return s, nil
}

func (s *inspectionModuleInspector) Name() resource.Name {
	return s.name
}

// detect calls the vision service and returns the label and confidence
// of the highest-confidence detection from the camera.
//
// TODO: Implement this method.
// 1. Call s.detector.DetectionsFromCamera() with s.cfg.Camera
// 2. If no detections, return "NO_DETECTION", 0, nil
// 3. Find the detection with the highest Score()
// 4. Return its Label() and Score()
func (s *inspectionModuleInspector) detect(ctx context.Context) (string, float64, error) {
	return "", 0, fmt.Errorf("not implemented: fill in the detect method")
}

func (s *inspectionModuleInspector) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := cmd["detect"]; ok {
		label, confidence, err := s.detect(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"label":      label,
			"confidence": confidence,
		}, nil
	}
	return nil, fmt.Errorf("unknown command: %v", cmd)
}

func (s *inspectionModuleInspector) Close(context.Context) error {
	s.cancelFunc()
	return nil
}
