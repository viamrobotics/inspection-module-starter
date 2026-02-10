package main

import (
	"context"
	"flag"
	"fmt"

	inspectionmodule "inspectionmodule"

	"github.com/erh/vmodutils"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/services/generic"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	ctx := context.Background()
	logger := logging.NewLogger("cli")

	host := flag.String("host", "", "Machine address (required)")
	flag.Parse()

	if *host == "" {
		return fmt.Errorf("need -host flag (get address from Viam app)")
	}

	logger.Infof("Connecting to %s...", *host)
	machine, err := vmodutils.ConnectToHostFromCLIToken(ctx, *host, logger)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer machine.Close(ctx)

	cfg := &inspectionmodule.Config{
		Camera:        "inspection-cam",
		VisionService: "can-detector",
	}

	deps, err := vmodutils.MachineToDependencies(machine)
	if err != nil {
		return fmt.Errorf("failed to get dependencies: %w", err)
	}

	inspector, err := inspectionmodule.NewInspector(
		ctx,
		deps,
		generic.Named("inspector"),
		cfg,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create inspector: %w", err)
	}

	result, err := inspector.DoCommand(ctx, map[string]interface{}{"detect": true})
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	label := result["label"].(string)
	confidence := result["confidence"].(float64)
	logger.Infof("Detection: %s (%.1f%% confidence)", label, confidence*100)
	return nil
}
