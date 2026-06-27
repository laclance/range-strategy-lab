package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"range-strategy-lab/internal/lab"
)

const defaultCSVPath = "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv"

var futuresRangeUniverseDiscoveryConfigForRun = lab.DefaultFuturesRangeUniverseDiscoveryAuditConfig
var futuresRangeUniverseBreakoutRetestAcceptanceBaselineConfigForRun = lab.DefaultFuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig
var futuresRangeUniverseStructuredCompressionBaselineConfigForRun = lab.DefaultFuturesRangeUniverseStructuredCompressionBaselineConfig
var futuresRangeUniverseStructuredCompressionOptimizationConfigForRun = lab.DefaultFuturesRangeUniverseStructuredCompressionOptimizationConfig
var futuresRangeUniverseStructuredCompressionStrategyReplayConfigForRun = lab.DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig
var futuresRangeUniverseStructuredCompressionWalkForwardConfigForRun = lab.DefaultFuturesRangeUniverseStructuredCompressionWalkForwardConfig
var futuresHigherTFNestedRangeRotationAuditConfigForRun = lab.DefaultFuturesHigherTFNestedRangeRotationAuditConfig
var futuresRangeFirstOccupancyRotationV1OptimizationConfigForRun = lab.DefaultFuturesRangeFirstOccupancyRotationV1OptimizationConfig
var futuresRangeContextTriageAuditConfigForRun = lab.DefaultFuturesRangeContextTriageAuditConfig
var futuresRangeStateConstructionLoopAuditConfigForRun = lab.DefaultFuturesRangeStateConstructionLoopAuditConfig
var futuresRangeContextRouterAuditConfigForRun = lab.DefaultFuturesRangeContextRouterAuditConfig
var futuresRangeRouterRotationPremiseAuditConfigForRun = lab.DefaultFuturesRangeRouterRotationPremiseAuditConfig

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	return runWithArgs(os.Args[1:])
}

func runWithArgs(args []string) error {
	fs := flag.NewFlagSet("rangelab", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	csvPath := fs.String("csv", defaultCSVPath, "5m BTCUSDT candle CSV")
	sourceProduct := fs.String("source-product", "", "source product: binance-usdm-futures or binance-spot")
	allowSpotComparison := fs.Bool("allow-spot-comparison", false, "allow an explicitly labeled Binance spot comparison run")
	outDir := fs.String("out-dir", "results/smoke", "output directory")
	startBalance := fs.Float64("start-balance", 1000, "starting balance")
	riskPct := fs.Float64("risk-pct", 0.01, "fraction of equity risked at stop")
	maxNotionalPct := fs.Float64("max-notional-pct", 1.0, "maximum entry notional as fraction of equity")
	feePct := fs.Float64("fee-pct", 0.0004, "fee fraction per side")
	slippagePct := fs.Float64("slippage-pct", 0.000116, "slippage fraction per side")
	maxHoldBars := fs.Int("max-hold-bars", 24, "default max hold bars stamped on placeholder signals")
	detector := fs.Bool("detector", false, "write detector-only range diagnostics")
	detectorSweep := fs.Bool("detector-sweep", false, "write detector sweep/audit diagnostics")
	detectorDurabilitySweep := fs.Bool("detector-durability-sweep", false, "write non-trading detector durability sweep diagnostics")
	detectorContextRefinementAudit := fs.Bool("detector-context-refinement-audit", false, "write non-trading detector context refinement diagnostics")
	holdInsideDirectionalEdgeAudit := fs.Bool("hold-inside-directional-edge-audit", false, "write non-trading hold-inside directional edge diagnostics")
	holdInsideMidlineTransitionAudit := fs.Bool("hold-inside-midline-transition-audit", false, "write non-trading hold-inside midline transition diagnostics")
	holdInsideMidlineReactionAudit := fs.Bool("hold-inside-midline-reaction-audit", false, "write non-trading hold-inside midline reaction diagnostics")
	holdInsideMidlineTouchPrototype := fs.Bool("hold-inside-midline-touch-prototype", false, "run offline hold-inside midline touch prototype")
	futuresImpulseAbsorptionAudit := fs.Bool("futures-impulse-absorption-audit", false, "write non-trading futures impulse absorption diagnostics")
	futuresRangeCandidateDiscoveryAudit := fs.Bool("futures-range-candidate-discovery-audit", false, "write non-trading futures range candidate discovery diagnostics")
	futuresRangeUniverseDiscoveryAudit := fs.Bool("futures-range-universe-discovery-audit", false, "write non-trading futures range universe discovery diagnostics")
	futuresCleanBreakoutBaselineBacktest := fs.Bool("futures-clean-breakout-baseline-backtest", false, "run offline futures clean breakout baseline backtest")
	futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest := fs.Bool("futures-range-universe-breakout-retest-acceptance-baseline-backtest", false, "run offline futures range universe breakout-retest acceptance baseline backtest")
	futuresRangeUniverseStructuredCompressionBaselineBacktest := fs.Bool("futures-range-universe-structured-compression-baseline-backtest", false, "run offline futures range universe structured compression baseline backtest")
	futuresRangeUniverseStructuredCompressionOptimization := fs.Bool("futures-range-universe-structured-compression-optimization", false, "run offline futures range universe structured compression optimization")
	futuresRangeUniverseStructuredCompressionStrategyReplay := fs.Bool("futures-range-universe-structured-compression-strategy-replay", false, "run offline futures range universe structured compression strategy replay")
	futuresRangeUniverseStructuredCompressionWalkForwardRobustness := fs.Bool("futures-range-universe-structured-compression-walk-forward-robustness", false, "run offline futures range universe structured compression walk-forward robustness")
	futuresHigherTFNestedRangeRotationAudit := fs.Bool("futures-higher-tf-nested-range-rotation-audit", false, "write non-trading futures higher-timeframe nested range rotation diagnostics")
	futuresRangeFirstOccupancyRotationV1Optimization := fs.Bool("futures-range-first-occupancy-rotation-v1-optimization", false, "run offline futures range-first occupancy rotation v1 optimization")
	futuresRangeContextTriageAudit := fs.Bool("futures-range-context-triage-audit", false, "write non-trading futures range context triage diagnostics")
	futuresRangeStateConstructionLoopAudit := fs.Bool("futures-range-state-construction-loop-audit", false, "write non-trading futures range-state construction loop diagnostics")
	futuresRangeContextRouterAudit := fs.Bool("futures-range-context-router-audit", false, "write non-trading futures range context router diagnostics")
	futuresRangeRouterRotationPremiseAudit := fs.Bool("futures-range-router-rotation-premise-audit", false, "write non-trading futures range router rotation premise diagnostics")
	srAudit := fs.Bool("sr-audit", false, "write go-sr support/resistance audit diagnostics")
	srBoundaryAudit := fs.Bool("sr-boundary-audit", false, "write non-trading SR boundary quality diagnostics")
	srBoundaryInspect := fs.Bool("sr-boundary-inspect", false, "write compact non-trading SR boundary candidate comparison diagnostics")
	srRejectionTimingAudit := fs.Bool("sr-rejection-timing-audit", false, "write compact non-trading SR rejection timing diagnostics")
	srConfirmationTimingAudit := fs.Bool("sr-confirmation-timing-audit", false, "write compact non-trading SR confirmation timing diagnostics")
	srFalseBreakReclaimTimingAudit := fs.Bool("sr-false-break-reclaim-timing-audit", false, "write compact non-trading SR false-break reclaim timing diagnostics")
	compressionBreakoutAudit := fs.Bool("compression-breakout-audit", false, "write compact non-trading compression breakout diagnostics")
	rangeRegimeDurabilityAudit := fs.Bool("range-regime-durability-audit", false, "write compact non-trading range regime durability diagnostics")
	detectorLookbackDays := fs.Int("detector-lookback-days", 20, "range detector trailing lookback in days")
	detectorPercentile := fs.Float64("detector-percentile", 0.30, "range detector low-compression percentile threshold")
	detectorMinConsecutiveBars := fs.Int("detector-min-consecutive-bars", 12, "range detector confirmed raw-active bars before active")
	detectorUseBollinger := fs.Bool("detector-use-bollinger", true, "include Bollinger20 width in range detector")
	detectorUseADX := fs.Bool("detector-use-adx", false, "include ADX14 in range detector")
	if err := fs.Parse(args); err != nil {
		return err
	}

	sourceProductWasSet := false
	outDirWasSet := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "source-product" {
			sourceProductWasSet = true
		}
		if f.Name == "out-dir" {
			outDirWasSet = true
		}
	})
	if *futuresRangeStateConstructionLoopAudit && !outDirWasSet {
		*outDir = "results/futures-range-state-construction-loop-audit"
	}
	if *futuresRangeContextRouterAudit && !outDirWasSet {
		*outDir = "results/futures-range-context-router-audit"
	}
	if *futuresRangeRouterRotationPremiseAudit && !outDirWasSet {
		*outDir = "results/futures-range-router-rotation-premise-audit"
	}

	product := *sourceProduct
	if *csvPath == defaultCSVPath && !sourceProductWasSet {
		product = lab.SourceProductBinanceUSDMFutures
	}
	if *csvPath != defaultCSVPath && !sourceProductWasSet {
		return fmt.Errorf("non-default -csv path requires explicit -source-product=%s or -source-product=%s", lab.SourceProductBinanceUSDMFutures, lab.SourceProductBinanceSpot)
	}

	candles, sourceManifest, err := lab.LoadResearchSourceCSV(*csvPath, lab.SourceValidationOptions{
		Product:             product,
		AllowSpotComparison: *allowSpotComparison,
	})
	if err != nil {
		return err
	}
	if len(candles) == 0 {
		return fmt.Errorf("no candles loaded")
	}

	cfg := lab.BacktestConfig{
		StartBalance:   *startBalance,
		RiskPct:        *riskPct,
		MaxNotionalPct: *maxNotionalPct,
		FeePct:         *feePct,
		SlippagePct:    *slippagePct,
		MaxHoldBars:    *maxHoldBars,
	}

	tradeProducingFlagSelected := *holdInsideMidlineTouchPrototype || *futuresCleanBreakoutBaselineBacktest || *futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest || *futuresRangeUniverseStructuredCompressionBaselineBacktest || *futuresRangeUniverseStructuredCompressionOptimization || *futuresRangeUniverseStructuredCompressionStrategyReplay || *futuresRangeUniverseStructuredCompressionWalkForwardRobustness || *futuresRangeFirstOccupancyRotationV1Optimization
	if *futuresHigherTFNestedRangeRotationAudit {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-higher-tf-nested-range-rotation-audit requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if tradeProducingFlagSelected {
			return fmt.Errorf("-futures-higher-tf-nested-range-rotation-audit cannot be combined with trade-producing prototype/backtest/optimization/replay/walk-forward flags")
		}
	}
	if *futuresRangeContextTriageAudit {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-context-triage-audit requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if tradeProducingFlagSelected {
			return fmt.Errorf("-futures-range-context-triage-audit cannot be combined with trade-producing prototype/backtest/optimization/replay/walk-forward flags")
		}
	}
	if *futuresRangeStateConstructionLoopAudit {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-state-construction-loop-audit requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if tradeProducingFlagSelected {
			return fmt.Errorf("-futures-range-state-construction-loop-audit cannot be combined with trade-producing prototype/backtest/optimization/replay/walk-forward flags")
		}
		if *detector || *detectorSweep || *detectorDurabilitySweep || *detectorContextRefinementAudit || *holdInsideDirectionalEdgeAudit || *holdInsideMidlineTransitionAudit || *holdInsideMidlineReactionAudit || *futuresImpulseAbsorptionAudit || *futuresRangeCandidateDiscoveryAudit || *futuresRangeUniverseDiscoveryAudit || *futuresHigherTFNestedRangeRotationAudit || *futuresRangeContextTriageAudit || *futuresRangeContextRouterAudit || *futuresRangeRouterRotationPremiseAudit || *srAudit || *srBoundaryAudit || *srBoundaryInspect || *srRejectionTimingAudit || *srConfirmationTimingAudit || *srFalseBreakReclaimTimingAudit || *compressionBreakoutAudit || *rangeRegimeDurabilityAudit {
			return fmt.Errorf("-futures-range-state-construction-loop-audit cannot be combined with other audit, detector, source-expansion, or diagnostic flags")
		}
	}
	if *futuresRangeContextRouterAudit {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-context-router-audit requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if tradeProducingFlagSelected {
			return fmt.Errorf("-futures-range-context-router-audit cannot be combined with trade-producing prototype/backtest/optimization/replay/walk-forward flags")
		}
		if *detector || *detectorSweep || *detectorDurabilitySweep || *detectorContextRefinementAudit || *holdInsideDirectionalEdgeAudit || *holdInsideMidlineTransitionAudit || *holdInsideMidlineReactionAudit || *futuresImpulseAbsorptionAudit || *futuresRangeCandidateDiscoveryAudit || *futuresRangeUniverseDiscoveryAudit || *futuresHigherTFNestedRangeRotationAudit || *futuresRangeContextTriageAudit || *futuresRangeStateConstructionLoopAudit || *futuresRangeRouterRotationPremiseAudit || *srAudit || *srBoundaryAudit || *srBoundaryInspect || *srRejectionTimingAudit || *srConfirmationTimingAudit || *srFalseBreakReclaimTimingAudit || *compressionBreakoutAudit || *rangeRegimeDurabilityAudit {
			return fmt.Errorf("-futures-range-context-router-audit cannot be combined with other audit, detector, source-expansion, or diagnostic flags")
		}
	}
	if *futuresRangeRouterRotationPremiseAudit {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-router-rotation-premise-audit requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if tradeProducingFlagSelected {
			return fmt.Errorf("%s", lab.RangeRouterRotationPremiseStopStateRejectedStrategyBacktest)
		}
		if *detector || *detectorSweep || *detectorDurabilitySweep || *detectorContextRefinementAudit || *holdInsideDirectionalEdgeAudit || *holdInsideMidlineTransitionAudit || *holdInsideMidlineReactionAudit || *futuresImpulseAbsorptionAudit || *futuresRangeCandidateDiscoveryAudit || *futuresRangeUniverseDiscoveryAudit || *futuresHigherTFNestedRangeRotationAudit || *futuresRangeContextTriageAudit || *futuresRangeStateConstructionLoopAudit || *futuresRangeContextRouterAudit || *srAudit || *srBoundaryAudit || *srBoundaryInspect || *srRejectionTimingAudit || *srConfirmationTimingAudit || *srFalseBreakReclaimTimingAudit || *compressionBreakoutAudit || *rangeRegimeDurabilityAudit {
			return fmt.Errorf("-futures-range-router-rotation-premise-audit cannot be combined with other audit, detector, source-expansion, or diagnostic flags")
		}
	}

	var strategy lab.Strategy = lab.EmptyStrategy{}
	strategyName := strategy.Name()
	var prototypeStrategy lab.HoldInsideMidlineTouchPrototypeStrategy
	if *holdInsideMidlineTouchPrototype {
		if *futuresCleanBreakoutBaselineBacktest || *futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest || *futuresRangeUniverseStructuredCompressionBaselineBacktest || *futuresRangeUniverseStructuredCompressionOptimization || *futuresRangeUniverseStructuredCompressionStrategyReplay || *futuresRangeUniverseStructuredCompressionWalkForwardRobustness || *futuresRangeFirstOccupancyRotationV1Optimization {
			return fmt.Errorf("-hold-inside-midline-touch-prototype cannot be combined with trade-producing backtest flags")
		}
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-hold-inside-midline-touch-prototype requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if *detectorLookbackDays <= 0 {
			return fmt.Errorf("detector lookback days must be positive")
		}
		detectorCfg := lab.DefaultCompressionRangeDetectorConfig()
		detectorCfg.LookbackDays = *detectorLookbackDays
		var err error
		prototypeStrategy, err = lab.NewHoldInsideMidlineTouchPrototypeStrategy(candles, detectorCfg, lab.HoldInsideMidlineTouchPrototypeConfig{}, lab.DefaultSplits())
		if err != nil {
			return err
		}
		strategy = prototypeStrategy
		strategyName = strategy.Name()
	}
	if *futuresImpulseAbsorptionAudit {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-impulse-absorption-audit requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
	}
	if *futuresRangeCandidateDiscoveryAudit {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-candidate-discovery-audit requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
	}
	if *futuresRangeUniverseDiscoveryAudit {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-universe-discovery-audit requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if *holdInsideMidlineTouchPrototype || *futuresCleanBreakoutBaselineBacktest || *futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest || *futuresRangeUniverseStructuredCompressionBaselineBacktest || *futuresRangeUniverseStructuredCompressionOptimization || *futuresRangeUniverseStructuredCompressionStrategyReplay || *futuresRangeUniverseStructuredCompressionWalkForwardRobustness || *futuresRangeFirstOccupancyRotationV1Optimization {
			return fmt.Errorf("-futures-range-universe-discovery-audit cannot be combined with trade-producing prototype/backtest flags")
		}
	}
	if *futuresCleanBreakoutBaselineBacktest {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-clean-breakout-baseline-backtest requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if *futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest || *futuresRangeUniverseStructuredCompressionBaselineBacktest || *futuresRangeUniverseStructuredCompressionOptimization || *futuresRangeUniverseStructuredCompressionStrategyReplay || *futuresRangeUniverseStructuredCompressionWalkForwardRobustness || *futuresRangeFirstOccupancyRotationV1Optimization {
			return fmt.Errorf("-futures-clean-breakout-baseline-backtest cannot be combined with other trade-producing range-universe flags")
		}
	}
	if *futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-universe-breakout-retest-acceptance-baseline-backtest requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if *futuresRangeUniverseStructuredCompressionBaselineBacktest || *futuresRangeUniverseStructuredCompressionOptimization || *futuresRangeUniverseStructuredCompressionStrategyReplay || *futuresRangeUniverseStructuredCompressionWalkForwardRobustness || *futuresRangeFirstOccupancyRotationV1Optimization {
			return fmt.Errorf("-futures-range-universe-breakout-retest-acceptance-baseline-backtest cannot be combined with structured-compression trade-producing flags")
		}
	}
	if *futuresRangeUniverseStructuredCompressionBaselineBacktest {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-universe-structured-compression-baseline-backtest requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if *futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest || *futuresRangeUniverseStructuredCompressionOptimization || *futuresRangeUniverseStructuredCompressionStrategyReplay || *futuresRangeUniverseStructuredCompressionWalkForwardRobustness || *futuresRangeFirstOccupancyRotationV1Optimization {
			return fmt.Errorf("-futures-range-universe-structured-compression-baseline-backtest cannot be combined with other trade-producing range-universe flags")
		}
	}
	if *futuresRangeUniverseStructuredCompressionOptimization {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-universe-structured-compression-optimization requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if *futuresRangeUniverseStructuredCompressionStrategyReplay || *futuresRangeUniverseStructuredCompressionWalkForwardRobustness || *futuresRangeFirstOccupancyRotationV1Optimization {
			return fmt.Errorf("-futures-range-universe-structured-compression-optimization cannot be combined with structured-compression replay or walk-forward flags")
		}
	}
	if *futuresRangeUniverseStructuredCompressionStrategyReplay {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-universe-structured-compression-strategy-replay requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if *futuresRangeUniverseStructuredCompressionWalkForwardRobustness || *futuresRangeFirstOccupancyRotationV1Optimization {
			return fmt.Errorf("-futures-range-universe-structured-compression-strategy-replay cannot be combined with -futures-range-universe-structured-compression-walk-forward-robustness")
		}
	}
	if *futuresRangeUniverseStructuredCompressionWalkForwardRobustness {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-universe-structured-compression-walk-forward-robustness requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
	}
	if *futuresRangeFirstOccupancyRotationV1Optimization {
		if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
			return fmt.Errorf("-futures-range-first-occupancy-rotation-v1-optimization requires Binance USDT-M futures source; got product=%q comparison_only=%t", sourceManifest.Product, sourceManifest.ComparisonOnly)
		}
		if *holdInsideMidlineTouchPrototype || *futuresCleanBreakoutBaselineBacktest || *futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest || *futuresRangeUniverseStructuredCompressionBaselineBacktest || *futuresRangeUniverseStructuredCompressionOptimization || *futuresRangeUniverseStructuredCompressionStrategyReplay || *futuresRangeUniverseStructuredCompressionWalkForwardRobustness {
			return fmt.Errorf("-futures-range-first-occupancy-rotation-v1-optimization cannot be combined with trade-producing prototype/baseline/optimization/replay/walk-forward flags")
		}
	}

	var cleanBreakoutResult lab.FuturesCleanBreakoutBaselineResult
	var breakoutRetestAcceptanceResult lab.FuturesRangeUniverseBreakoutRetestAcceptanceBaselineResult
	var structuredCompressionResult lab.FuturesRangeUniverseStructuredCompressionBaselineResult
	var structuredCompressionOptimizationResult lab.FuturesRangeUniverseStructuredCompressionOptimizationResult
	var structuredCompressionStrategyReplayResult lab.FuturesRangeUniverseStructuredCompressionStrategyReplayResult
	var structuredCompressionWalkForwardResult lab.FuturesRangeUniverseStructuredCompressionWalkForwardResult
	var nestedRangeRotationResult lab.FuturesHigherTFNestedRangeRotationAuditResult
	var occupancyRotationV1Result lab.FuturesRangeFirstOccupancyRotationV1OptimizationResult
	var rangeContextTriageResult lab.FuturesRangeContextTriageAuditResult
	var rangeStateConstructionLoopResult lab.FuturesRangeStateConstructionLoopAuditResult
	var rangeContextRouterResult lab.FuturesRangeContextRouterAuditResult
	var rangeRouterRotationPremiseResult lab.FuturesRangeRouterRotationPremiseAuditResult
	var result lab.BacktestResult
	if *futuresCleanBreakoutBaselineBacktest {
		var err error
		cleanBreakoutResult, err = lab.RunFuturesCleanBreakoutBaselineBacktest(candles, lab.DefaultFuturesCleanBreakoutBaselineConfig(), cfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		result = lab.BacktestResult{Trades: cleanBreakoutResult.Trades}
		strategyName = lab.FuturesCleanBreakoutBaselineName
	} else if *futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest {
		var err error
		breakoutRetestAcceptanceCfg := futuresRangeUniverseBreakoutRetestAcceptanceBaselineConfigForRun()
		breakoutRetestAcceptanceResult, err = lab.RunFuturesRangeUniverseBreakoutRetestAcceptanceBaselineBacktest(breakoutRetestAcceptanceCfg, cfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		result = lab.BacktestResult{Trades: breakoutRetestAcceptanceResult.Trades}
		strategyName = lab.FuturesRangeUniverseBreakoutRetestAcceptanceBaselineName
	} else if *futuresRangeUniverseStructuredCompressionBaselineBacktest {
		var err error
		structuredCompressionCfg := futuresRangeUniverseStructuredCompressionBaselineConfigForRun()
		structuredCompressionResult, err = lab.RunFuturesRangeUniverseStructuredCompressionBaselineBacktest(structuredCompressionCfg, cfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		result = lab.BacktestResult{Trades: structuredCompressionResult.Trades}
		strategyName = lab.FuturesRangeUniverseStructuredCompressionBaselineName
	} else if *futuresRangeUniverseStructuredCompressionOptimization {
		var err error
		optimizationCfg := futuresRangeUniverseStructuredCompressionOptimizationConfigForRun()
		structuredCompressionOptimizationResult, err = lab.RunFuturesRangeUniverseStructuredCompressionOptimization(optimizationCfg, cfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		result = lab.BacktestResult{Trades: structuredCompressionOptimizationResult.Trades}
		strategyName = lab.FuturesRangeUniverseStructuredCompressionOptimizationName
	} else if *futuresRangeUniverseStructuredCompressionStrategyReplay {
		var err error
		replayCfg := futuresRangeUniverseStructuredCompressionStrategyReplayConfigForRun()
		structuredCompressionStrategyReplayResult, err = lab.RunFuturesRangeUniverseStructuredCompressionStrategyReplay(replayCfg, cfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		result = lab.BacktestResult{Trades: structuredCompressionStrategyReplayResult.Trades}
		strategyName = lab.FuturesRangeUniverseStructuredCompressionStrategyReplayName
	} else if *futuresRangeUniverseStructuredCompressionWalkForwardRobustness {
		var err error
		walkForwardCfg := futuresRangeUniverseStructuredCompressionWalkForwardConfigForRun()
		structuredCompressionWalkForwardResult, err = lab.RunFuturesRangeUniverseStructuredCompressionWalkForwardRobustness(walkForwardCfg, cfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		result = lab.BacktestResult{Trades: structuredCompressionWalkForwardResult.Trades}
		strategyName = lab.FuturesRangeUniverseStructuredCompressionWalkForwardName
	} else if *futuresRangeFirstOccupancyRotationV1Optimization {
		var err error
		occupancyCfg := futuresRangeFirstOccupancyRotationV1OptimizationConfigForRun()
		occupancyRotationV1Result, err = lab.RunFuturesRangeFirstOccupancyRotationV1Optimization(candles, sourceManifest, occupancyCfg, cfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		result = lab.BacktestResult{Trades: occupancyRotationV1Result.Trades}
		strategyName = lab.FuturesRangeFirstOccupancyRotationV1OptimizationName
	} else {
		result = lab.RunBacktest(candles, strategy, cfg)
	}
	summaries := lab.SummarizeSplits(result.Trades, *startBalance, lab.DefaultSplits())

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(*outDir, "source_manifest.json"), sourceManifest); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(*outDir, "trades.json"), result.Trades); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(*outDir, "summary.json"), summaries); err != nil {
		return err
	}
	if err := writeSummaryCSV(filepath.Join(*outDir, "summary.csv"), summaries); err != nil {
		return err
	}
	if *holdInsideMidlineTouchPrototype {
		signalRows := prototypeStrategy.SignalRows()
		tradeRows := prototypeStrategy.TradeRows(result.Trades, lab.DefaultSplits())
		prototypeSummaryRows := lab.SummarizeHoldInsideMidlineTouchPrototype(tradeRows, *startBalance, lab.DefaultSplits())
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_touch_prototype_signals.json"), signalRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineTouchPrototypeSignalsCSV(filepath.Join(*outDir, "hold_inside_midline_touch_prototype_signals.csv"), signalRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_touch_prototype_trades.json"), tradeRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineTouchPrototypeTradesCSV(filepath.Join(*outDir, "hold_inside_midline_touch_prototype_trades.csv"), tradeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_touch_prototype_summary.json"), prototypeSummaryRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineTouchPrototypeSummaryCSV(filepath.Join(*outDir, "hold_inside_midline_touch_prototype_summary.csv"), prototypeSummaryRows); err != nil {
			return err
		}
		fmt.Printf("hold_inside_midline_touch_prototype signal_rows=%d trades=%d summary_rows=%d stop_state=%s\n",
			len(signalRows),
			len(tradeRows),
			len(prototypeSummaryRows),
			lab.HoldInsideMidlineTouchPrototypeStopState(prototypeSummaryRows),
		)
	}
	if *futuresImpulseAbsorptionAudit {
		absorptionCfg := lab.DefaultFuturesImpulseAbsorptionAuditConfig()
		candidateRows, summaryRows, stabilityRows, err := lab.RunFuturesImpulseAbsorptionAudit(candles, absorptionCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_impulse_absorption_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeFuturesImpulseAbsorptionCandidatesCSV(filepath.Join(*outDir, "futures_impulse_absorption_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_impulse_absorption_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeFuturesImpulseAbsorptionSummaryCSV(filepath.Join(*outDir, "futures_impulse_absorption_summary.csv"), summaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_impulse_absorption_stability.json"), stabilityRows); err != nil {
			return err
		}
		if err := writeFuturesImpulseAbsorptionStabilityCSV(filepath.Join(*outDir, "futures_impulse_absorption_stability.csv"), stabilityRows); err != nil {
			return err
		}
		fmt.Printf("futures_impulse_absorption_audit candidate_rows=%d summary_rows=%d stability_rows=%d warmup_bars=%d horizons=%s stop_state=%s\n",
			len(candidateRows),
			len(summaryRows),
			len(stabilityRows),
			absorptionCfg.WarmupBars,
			formatIntSlice(absorptionCfg.HorizonsBars),
			lab.FuturesImpulseAbsorptionReviewStopState(summaryRows, lab.DefaultSplits()),
		)
	}
	if *futuresRangeCandidateDiscoveryAudit {
		discoveryCfg := lab.DefaultFuturesRangeCandidateDiscoveryAuditConfig()
		coverageRows, candidateRows, summaryRows, rankingRows, stabilityRows, err := lab.RunFuturesRangeCandidateDiscoveryAudit(candles, discoveryCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_discovery_coverage.json"), coverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeDiscoveryCoverageCSV(filepath.Join(*outDir, "futures_range_discovery_coverage.csv"), coverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_discovery_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeFuturesRangeDiscoveryCandidatesCSV(filepath.Join(*outDir, "futures_range_discovery_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_discovery_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeDiscoverySummaryCSV(filepath.Join(*outDir, "futures_range_discovery_summary.csv"), summaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_discovery_rankings.json"), rankingRows); err != nil {
			return err
		}
		if err := writeFuturesRangeDiscoveryRankingsCSV(filepath.Join(*outDir, "futures_range_discovery_rankings.csv"), rankingRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_discovery_stability.json"), stabilityRows); err != nil {
			return err
		}
		if err := writeFuturesRangeDiscoveryStabilityCSV(filepath.Join(*outDir, "futures_range_discovery_stability.csv"), stabilityRows); err != nil {
			return err
		}
		fmt.Printf("futures_range_candidate_discovery_audit coverage_rows=%d candidate_rows=%d summary_rows=%d ranking_rows=%d stability_rows=%d horizons=%s stop_state=%s\n",
			len(coverageRows),
			len(candidateRows),
			len(summaryRows),
			len(rankingRows),
			len(stabilityRows),
			formatIntSlice(discoveryCfg.HorizonsBars),
			lab.FuturesRangeDiscoveryReviewStopState(rankingRows),
		)
	}
	if *futuresRangeUniverseDiscoveryAudit {
		universeCfg := futuresRangeUniverseDiscoveryConfigForRun()
		universeResult, err := lab.RunFuturesRangeUniverseDiscoveryAudit(universeCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_sources.json"), universeResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseSourcesCSV(filepath.Join(*outDir, "futures_range_universe_sources.csv"), universeResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_coverage.json"), universeResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseCoverageCSV(filepath.Join(*outDir, "futures_range_universe_coverage.csv"), universeResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_candidates.json"), universeResult.CandidateRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseCandidatesCSV(filepath.Join(*outDir, "futures_range_universe_candidates.csv"), universeResult.CandidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_summary.json"), universeResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseSummaryCSV(filepath.Join(*outDir, "futures_range_universe_summary.csv"), universeResult.SummaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_rankings.json"), universeResult.RankingRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseRankingsCSV(filepath.Join(*outDir, "futures_range_universe_rankings.csv"), universeResult.RankingRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_stability.json"), universeResult.StabilityRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStabilityCSV(filepath.Join(*outDir, "futures_range_universe_stability.csv"), universeResult.StabilityRows); err != nil {
			return err
		}
		fmt.Printf("futures_range_universe_discovery_audit source_rows=%d coverage_rows=%d candidate_rows=%d summary_rows=%d ranking_rows=%d stability_rows=%d horizons=%s stop_state=%s\n",
			len(universeResult.SourceRows),
			len(universeResult.CoverageRows),
			len(universeResult.CandidateRows),
			len(universeResult.SummaryRows),
			len(universeResult.RankingRows),
			len(universeResult.StabilityRows),
			formatIntSlice(universeCfg.Discovery.HorizonsBars),
			lab.FuturesRangeUniverseReviewStopState(universeResult.RankingRows),
		)
	}
	if *futuresCleanBreakoutBaselineBacktest {
		if err := writeJSON(filepath.Join(*outDir, "futures_clean_breakout_baseline_signals.json"), cleanBreakoutResult.SignalRows); err != nil {
			return err
		}
		if err := writeFuturesCleanBreakoutBaselineSignalsCSV(filepath.Join(*outDir, "futures_clean_breakout_baseline_signals.csv"), cleanBreakoutResult.SignalRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_clean_breakout_baseline_trades.json"), cleanBreakoutResult.TradeRows); err != nil {
			return err
		}
		if err := writeFuturesCleanBreakoutBaselineTradesCSV(filepath.Join(*outDir, "futures_clean_breakout_baseline_trades.csv"), cleanBreakoutResult.TradeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_clean_breakout_baseline_summary.json"), cleanBreakoutResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesCleanBreakoutBaselineSummaryCSV(filepath.Join(*outDir, "futures_clean_breakout_baseline_summary.csv"), cleanBreakoutResult.SummaryRows); err != nil {
			return err
		}
		coverage := formatCleanBreakoutCoverage(cleanBreakoutResult.CoverageRows)
		fmt.Printf("futures_clean_breakout_baseline signal_rows=%d trades=%d summary_rows=%d coverage=%s stop_state=%s\n",
			len(cleanBreakoutResult.SignalRows),
			len(cleanBreakoutResult.TradeRows),
			len(cleanBreakoutResult.SummaryRows),
			coverage,
			lab.FuturesCleanBreakoutBaselineStopState(cleanBreakoutResult.SummaryRows, lab.DefaultFuturesCleanBreakoutBaselineConfig(), *startBalance, lab.DefaultSplits()),
		)
	}
	if *futuresRangeUniverseBreakoutRetestAcceptanceBaselineBacktest {
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_sources.json"), breakoutRetestAcceptanceResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineSourcesCSV(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_sources.csv"), breakoutRetestAcceptanceResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_coverage.json"), breakoutRetestAcceptanceResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineCoverageCSV(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_coverage.csv"), breakoutRetestAcceptanceResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_selection.json"), breakoutRetestAcceptanceResult.SelectionRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineSelectionCSV(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_selection.csv"), breakoutRetestAcceptanceResult.SelectionRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_signals.json"), breakoutRetestAcceptanceResult.SignalRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineSignalsCSV(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_signals.csv"), breakoutRetestAcceptanceResult.SignalRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_trades.json"), breakoutRetestAcceptanceResult.TradeRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineTradesCSV(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_trades.csv"), breakoutRetestAcceptanceResult.TradeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_summary.json"), breakoutRetestAcceptanceResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineSummaryCSV(filepath.Join(*outDir, "futures_range_universe_breakout_retest_acceptance_baseline_summary.csv"), breakoutRetestAcceptanceResult.SummaryRows); err != nil {
			return err
		}
		breakoutRetestAcceptanceCfg := futuresRangeUniverseBreakoutRetestAcceptanceBaselineConfigForRun()
		fmt.Printf("futures_range_universe_breakout_retest_acceptance_baseline source_rows=%d coverage_rows=%d selection_rows=%d signal_rows=%d trades=%d summary_rows=%d stop_state=%s\n",
			len(breakoutRetestAcceptanceResult.SourceRows),
			len(breakoutRetestAcceptanceResult.CoverageRows),
			len(breakoutRetestAcceptanceResult.SelectionRows),
			len(breakoutRetestAcceptanceResult.SignalRows),
			len(breakoutRetestAcceptanceResult.TradeRows),
			len(breakoutRetestAcceptanceResult.SummaryRows),
			lab.FuturesRangeUniverseBreakoutRetestAcceptanceBaselineStopState(breakoutRetestAcceptanceResult, breakoutRetestAcceptanceCfg, *startBalance, lab.DefaultSplits()),
		)
	}
	if *futuresRangeUniverseStructuredCompressionBaselineBacktest {
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_sources.json"), structuredCompressionResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionBaselineSourcesCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_sources.csv"), structuredCompressionResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_coverage.json"), structuredCompressionResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionBaselineCoverageCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_coverage.csv"), structuredCompressionResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_signals.json"), structuredCompressionResult.SignalRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionBaselineSignalsCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_signals.csv"), structuredCompressionResult.SignalRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_trades.json"), structuredCompressionResult.TradeRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionBaselineTradesCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_trades.csv"), structuredCompressionResult.TradeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_summary.json"), structuredCompressionResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionBaselineSummaryCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_baseline_summary.csv"), structuredCompressionResult.SummaryRows); err != nil {
			return err
		}
		structuredCompressionCfg := futuresRangeUniverseStructuredCompressionBaselineConfigForRun()
		fmt.Printf("futures_range_universe_structured_compression_baseline source_rows=%d coverage_rows=%d signal_rows=%d trades=%d summary_rows=%d stop_state=%s\n",
			len(structuredCompressionResult.SourceRows),
			len(structuredCompressionResult.CoverageRows),
			len(structuredCompressionResult.SignalRows),
			len(structuredCompressionResult.TradeRows),
			len(structuredCompressionResult.SummaryRows),
			lab.FuturesRangeUniverseStructuredCompressionBaselineStopState(structuredCompressionResult.SummaryRows, structuredCompressionCfg, *startBalance, lab.DefaultSplits()),
		)
	}
	if *futuresRangeUniverseStructuredCompressionOptimization {
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_sources.json"), structuredCompressionOptimizationResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionOptimizationSourcesCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_sources.csv"), structuredCompressionOptimizationResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_coverage.json"), structuredCompressionOptimizationResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionOptimizationCoverageCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_coverage.csv"), structuredCompressionOptimizationResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_grid.json"), structuredCompressionOptimizationResult.GridRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionOptimizationGridCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_grid.csv"), structuredCompressionOptimizationResult.GridRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_trades.json"), structuredCompressionOptimizationResult.TradeRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionOptimizationTradesCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_trades.csv"), structuredCompressionOptimizationResult.TradeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_summary.json"), structuredCompressionOptimizationResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionOptimizationSummaryCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_summary.csv"), structuredCompressionOptimizationResult.SummaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_rankings.json"), structuredCompressionOptimizationResult.RankingRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionOptimizationRankingsCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_optimization_rankings.csv"), structuredCompressionOptimizationResult.RankingRows); err != nil {
			return err
		}
		fmt.Printf("futures_range_universe_structured_compression_optimization source_rows=%d coverage_rows=%d grid_rows=%d trades=%d summary_rows=%d ranking_rows=%d selected_config=%s stop_state=%s\n",
			len(structuredCompressionOptimizationResult.SourceRows),
			len(structuredCompressionOptimizationResult.CoverageRows),
			len(structuredCompressionOptimizationResult.GridRows),
			len(structuredCompressionOptimizationResult.TradeRows),
			len(structuredCompressionOptimizationResult.SummaryRows),
			len(structuredCompressionOptimizationResult.RankingRows),
			structuredCompressionOptimizationResult.SelectedConfigID,
			lab.FuturesRangeUniverseStructuredCompressionOptimizationStopState(structuredCompressionOptimizationResult.RankingRows),
		)
	}
	if *futuresRangeUniverseStructuredCompressionStrategyReplay {
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_sources.json"), structuredCompressionStrategyReplayResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionStrategyReplaySourcesCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_sources.csv"), structuredCompressionStrategyReplayResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_coverage.json"), structuredCompressionStrategyReplayResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionStrategyReplayCoverageCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_coverage.csv"), structuredCompressionStrategyReplayResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_signals.json"), structuredCompressionStrategyReplayResult.SignalRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionStrategyReplaySignalsCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_signals.csv"), structuredCompressionStrategyReplayResult.SignalRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_trades.json"), structuredCompressionStrategyReplayResult.TradeRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionStrategyReplayTradesCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_trades.csv"), structuredCompressionStrategyReplayResult.TradeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_summary.json"), structuredCompressionStrategyReplayResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionStrategyReplaySummaryCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_strategy_summary.csv"), structuredCompressionStrategyReplayResult.SummaryRows); err != nil {
			return err
		}
		replayCfg := futuresRangeUniverseStructuredCompressionStrategyReplayConfigForRun()
		fmt.Printf("futures_range_universe_structured_compression_strategy_replay source_rows=%d coverage_rows=%d signal_rows=%d trades=%d summary_rows=%d stop_state=%s\n",
			len(structuredCompressionStrategyReplayResult.SourceRows),
			len(structuredCompressionStrategyReplayResult.CoverageRows),
			len(structuredCompressionStrategyReplayResult.SignalRows),
			len(structuredCompressionStrategyReplayResult.TradeRows),
			len(structuredCompressionStrategyReplayResult.SummaryRows),
			lab.FuturesRangeUniverseStructuredCompressionStrategyReplayStopState(structuredCompressionStrategyReplayResult.SummaryRows, replayCfg, *startBalance, lab.DefaultSplits()),
		)
	}
	if *futuresRangeUniverseStructuredCompressionWalkForwardRobustness {
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_sources.json"), structuredCompressionWalkForwardResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionWalkForwardSourcesCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_sources.csv"), structuredCompressionWalkForwardResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_coverage.json"), structuredCompressionWalkForwardResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionWalkForwardCoverageCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_coverage.csv"), structuredCompressionWalkForwardResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_grid.json"), structuredCompressionWalkForwardResult.GridRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionWalkForwardGridCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_grid.csv"), structuredCompressionWalkForwardResult.GridRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_folds.json"), structuredCompressionWalkForwardResult.FoldRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionWalkForwardFoldsCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_folds.csv"), structuredCompressionWalkForwardResult.FoldRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_trades.json"), structuredCompressionWalkForwardResult.TradeRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionWalkForwardTradesCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_trades.csv"), structuredCompressionWalkForwardResult.TradeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_summary.json"), structuredCompressionWalkForwardResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionWalkForwardSummaryCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_summary.csv"), structuredCompressionWalkForwardResult.SummaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_rankings.json"), structuredCompressionWalkForwardResult.RankingRows); err != nil {
			return err
		}
		if err := writeFuturesRangeUniverseStructuredCompressionWalkForwardRankingsCSV(filepath.Join(*outDir, "futures_range_universe_structured_compression_walk_forward_rankings.csv"), structuredCompressionWalkForwardResult.RankingRows); err != nil {
			return err
		}
		fmt.Printf("futures_range_universe_structured_compression_walk_forward source_rows=%d coverage_rows=%d grid_rows=%d fold_rows=%d trades=%d summary_rows=%d ranking_rows=%d frozen_config=%s stop_state=%s\n",
			len(structuredCompressionWalkForwardResult.SourceRows),
			len(structuredCompressionWalkForwardResult.CoverageRows),
			len(structuredCompressionWalkForwardResult.GridRows),
			len(structuredCompressionWalkForwardResult.FoldRows),
			len(structuredCompressionWalkForwardResult.TradeRows),
			len(structuredCompressionWalkForwardResult.SummaryRows),
			len(structuredCompressionWalkForwardResult.RankingRows),
			structuredCompressionWalkForwardResult.FrozenConfigID,
			lab.FuturesRangeUniverseStructuredCompressionWalkForwardStopState(structuredCompressionWalkForwardResult.FoldRows),
		)
	}
	if *futuresRangeFirstOccupancyRotationV1Optimization {
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_sources.json"), occupancyRotationV1Result.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1SourcesCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_sources.csv"), occupancyRotationV1Result.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_coverage.json"), occupancyRotationV1Result.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1CoverageCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_coverage.csv"), occupancyRotationV1Result.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_grid.json"), occupancyRotationV1Result.GridRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1GridCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_grid.csv"), occupancyRotationV1Result.GridRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_baseline.json"), occupancyRotationV1Result.BaselineRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1BaselineCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_baseline.csv"), occupancyRotationV1Result.BaselineRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_signals.json"), occupancyRotationV1Result.SignalRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1SignalsCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_signals.csv"), occupancyRotationV1Result.SignalRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_trades.json"), occupancyRotationV1Result.TradeRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1TradesCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_trades.csv"), occupancyRotationV1Result.TradeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_summary.json"), occupancyRotationV1Result.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1SummaryCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_summary.csv"), occupancyRotationV1Result.SummaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_rankings.json"), occupancyRotationV1Result.RankingRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1RankingsCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_rankings.csv"), occupancyRotationV1Result.RankingRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_selection.json"), occupancyRotationV1Result.SelectionRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1SelectionCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_selection.csv"), occupancyRotationV1Result.SelectionRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_skips.json"), occupancyRotationV1Result.SkipRows); err != nil {
			return err
		}
		if err := writeFuturesRangeFirstOccupancyRotationV1SkipsCSV(filepath.Join(*outDir, "futures_range_first_occupancy_rotation_v1_skips.csv"), occupancyRotationV1Result.SkipRows); err != nil {
			return err
		}
		fmt.Printf("futures_range_first_occupancy_rotation_v1 source_rows=%d coverage_rows=%d grid_rows=%d baseline_trades=%d selected_config=%s selected_trades=%d summary_rows=%d ranking_rows=%d stop_state=%s\n",
			len(occupancyRotationV1Result.SourceRows),
			len(occupancyRotationV1Result.CoverageRows),
			len(occupancyRotationV1Result.GridRows),
			len(occupancyRotationV1Result.Trades),
			occupancyRotationV1Result.SelectedConfigID,
			len(occupancyRotationV1Result.TradeRows)-len(occupancyRotationV1Result.Trades),
			len(occupancyRotationV1Result.SummaryRows),
			len(occupancyRotationV1Result.RankingRows),
			occupancyRotationV1Result.StopState,
		)
	}
	if *futuresHigherTFNestedRangeRotationAudit {
		nestedCfg := futuresHigherTFNestedRangeRotationAuditConfigForRun()
		var err error
		nestedRangeRotationResult, err = lab.RunFuturesHigherTFNestedRangeRotationAudit(candles, sourceManifest, nestedCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_sources.json"), nestedRangeRotationResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesHigherTFNestedRangeRotationSourcesCSV(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_sources.csv"), nestedRangeRotationResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_coverage.json"), nestedRangeRotationResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesHigherTFNestedRangeRotationCoverageCSV(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_coverage.csv"), nestedRangeRotationResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_parent_ranges.json"), nestedRangeRotationResult.ParentRangeRows); err != nil {
			return err
		}
		if err := writeFuturesHigherTFNestedRangeRotationParentRangesCSV(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_parent_ranges.csv"), nestedRangeRotationResult.ParentRangeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_child_ranges.json"), nestedRangeRotationResult.ChildRangeRows); err != nil {
			return err
		}
		if err := writeFuturesHigherTFNestedRangeRotationChildRangesCSV(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_child_ranges.csv"), nestedRangeRotationResult.ChildRangeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_events.json"), nestedRangeRotationResult.EventRows); err != nil {
			return err
		}
		if err := writeFuturesHigherTFNestedRangeRotationEventsCSV(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_events.csv"), nestedRangeRotationResult.EventRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_summary.json"), nestedRangeRotationResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesHigherTFNestedRangeRotationSummaryCSV(filepath.Join(*outDir, "futures_higher_tf_nested_range_rotation_summary.csv"), nestedRangeRotationResult.SummaryRows); err != nil {
			return err
		}
		fmt.Printf("futures_higher_tf_nested_range_rotation_audit source_rows=%d coverage_rows=%d parent_ranges=%d child_ranges=%d events=%d summary_rows=%d stop_state=%s\n",
			len(nestedRangeRotationResult.SourceRows),
			len(nestedRangeRotationResult.CoverageRows),
			len(nestedRangeRotationResult.ParentRangeRows),
			len(nestedRangeRotationResult.ChildRangeRows),
			len(nestedRangeRotationResult.EventRows),
			len(nestedRangeRotationResult.SummaryRows),
			lab.FuturesHigherTFNestedRangeRotationAuditStopState(nestedRangeRotationResult),
		)
	}
	if *futuresRangeContextTriageAudit {
		triageCfg := futuresRangeContextTriageAuditConfigForRun()
		var err error
		rangeContextTriageResult, err = lab.RunFuturesRangeContextTriageAudit(candles, sourceManifest, triageCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_triage_sources.json"), rangeContextTriageResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextTriageSourcesCSV(filepath.Join(*outDir, "futures_range_context_triage_sources.csv"), rangeContextTriageResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_triage_coverage.json"), rangeContextTriageResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextTriageCoverageCSV(filepath.Join(*outDir, "futures_range_context_triage_coverage.csv"), rangeContextTriageResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_triage_episodes.json"), rangeContextTriageResult.EpisodeRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextTriageEpisodesCSV(filepath.Join(*outDir, "futures_range_context_triage_episodes.csv"), rangeContextTriageResult.EpisodeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_triage_quality.json"), rangeContextTriageResult.QualityRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextTriageQualityCSV(filepath.Join(*outDir, "futures_range_context_triage_quality.csv"), rangeContextTriageResult.QualityRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_triage_sessions.json"), rangeContextTriageResult.SessionRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextTriageSessionsCSV(filepath.Join(*outDir, "futures_range_context_triage_sessions.csv"), rangeContextTriageResult.SessionRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_triage_failure_modes.json"), rangeContextTriageResult.FailureModeRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextTriageFailureModesCSV(filepath.Join(*outDir, "futures_range_context_triage_failure_modes.csv"), rangeContextTriageResult.FailureModeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_triage_cohorts.json"), rangeContextTriageResult.CohortRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextTriageCohortsCSV(filepath.Join(*outDir, "futures_range_context_triage_cohorts.csv"), rangeContextTriageResult.CohortRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_triage_rankings.json"), rangeContextTriageResult.RankingRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextTriageRankingsCSV(filepath.Join(*outDir, "futures_range_context_triage_rankings.csv"), rangeContextTriageResult.RankingRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_triage_summary.json"), rangeContextTriageResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextTriageSummaryCSV(filepath.Join(*outDir, "futures_range_context_triage_summary.csv"), rangeContextTriageResult.SummaryRows); err != nil {
			return err
		}
		fmt.Printf("futures_range_context_triage_audit source_rows=%d coverage_rows=%d episode_rows=%d quality_rows=%d session_rows=%d failure_mode_rows=%d cohort_rows=%d ranking_rows=%d passing_cohorts=%d stop_state=%s\n",
			len(rangeContextTriageResult.SourceRows),
			len(rangeContextTriageResult.CoverageRows),
			len(rangeContextTriageResult.EpisodeRows),
			len(rangeContextTriageResult.QualityRows),
			len(rangeContextTriageResult.SessionRows),
			len(rangeContextTriageResult.FailureModeRows),
			len(rangeContextTriageResult.CohortRows),
			len(rangeContextTriageResult.RankingRows),
			rangeContextTriageResult.PassingCohorts,
			rangeContextTriageResult.StopState,
		)
	}
	if *futuresRangeStateConstructionLoopAudit {
		stateCfg := futuresRangeStateConstructionLoopAuditConfigForRun()
		var err error
		rangeStateConstructionLoopResult, err = lab.RunFuturesRangeStateConstructionLoopAudit(candles, sourceManifest, stateCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_state_construction_loop_sources.json"), rangeStateConstructionLoopResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeStateConstructionLoopSourcesCSV(filepath.Join(*outDir, "futures_range_state_construction_loop_sources.csv"), rangeStateConstructionLoopResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_state_construction_loop_coverage.json"), rangeStateConstructionLoopResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeStateConstructionLoopCoverageCSV(filepath.Join(*outDir, "futures_range_state_construction_loop_coverage.csv"), rangeStateConstructionLoopResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_state_construction_loop_feature_windows.json"), rangeStateConstructionLoopResult.FeatureWindowRows); err != nil {
			return err
		}
		if err := writeFuturesRangeStateConstructionLoopFeatureWindowsCSV(filepath.Join(*outDir, "futures_range_state_construction_loop_feature_windows.csv"), rangeStateConstructionLoopResult.FeatureWindowRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_state_construction_loop_states.json"), rangeStateConstructionLoopResult.StateRows); err != nil {
			return err
		}
		if err := writeFuturesRangeStateConstructionLoopStatesCSV(filepath.Join(*outDir, "futures_range_state_construction_loop_states.csv"), rangeStateConstructionLoopResult.StateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_state_construction_loop_labels.json"), rangeStateConstructionLoopResult.LabelRows); err != nil {
			return err
		}
		if err := writeFuturesRangeStateConstructionLoopLabelsCSV(filepath.Join(*outDir, "futures_range_state_construction_loop_labels.csv"), rangeStateConstructionLoopResult.LabelRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_state_construction_loop_cohorts.json"), rangeStateConstructionLoopResult.CohortRows); err != nil {
			return err
		}
		if err := writeFuturesRangeStateConstructionLoopCohortsCSV(filepath.Join(*outDir, "futures_range_state_construction_loop_cohorts.csv"), rangeStateConstructionLoopResult.CohortRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_state_construction_loop_rankings.json"), rangeStateConstructionLoopResult.RankingRows); err != nil {
			return err
		}
		if err := writeFuturesRangeStateConstructionLoopRankingsCSV(filepath.Join(*outDir, "futures_range_state_construction_loop_rankings.csv"), rangeStateConstructionLoopResult.RankingRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_state_construction_loop_summary.json"), rangeStateConstructionLoopResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeStateConstructionLoopSummaryCSV(filepath.Join(*outDir, "futures_range_state_construction_loop_summary.csv"), rangeStateConstructionLoopResult.SummaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_state_construction_loop_skips.json"), rangeStateConstructionLoopResult.SkipRows); err != nil {
			return err
		}
		if err := writeFuturesRangeStateConstructionLoopSkipsCSV(filepath.Join(*outDir, "futures_range_state_construction_loop_skips.csv"), rangeStateConstructionLoopResult.SkipRows); err != nil {
			return err
		}
		fmt.Printf("futures_range_state_construction_loop_audit source_rows=%d coverage_rows=%d feature_window_rows=%d state_rows=%d label_rows=%d cohort_rows=%d ranking_rows=%d passing_cohorts=%d stop_state=%s\n",
			len(rangeStateConstructionLoopResult.SourceRows),
			len(rangeStateConstructionLoopResult.CoverageRows),
			len(rangeStateConstructionLoopResult.FeatureWindowRows),
			len(rangeStateConstructionLoopResult.StateRows),
			len(rangeStateConstructionLoopResult.LabelRows),
			len(rangeStateConstructionLoopResult.CohortRows),
			len(rangeStateConstructionLoopResult.RankingRows),
			rangeStateConstructionLoopResult.PassingCohorts,
			rangeStateConstructionLoopResult.StopState,
		)
	}
	if *futuresRangeContextRouterAudit {
		routerCfg := futuresRangeContextRouterAuditConfigForRun()
		var err error
		rangeContextRouterResult, err = lab.RunFuturesRangeContextRouterAudit(candles, sourceManifest, routerCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_router_sources.json"), rangeContextRouterResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextRouterSourcesCSV(filepath.Join(*outDir, "futures_range_context_router_sources.csv"), rangeContextRouterResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_router_coverage.json"), rangeContextRouterResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextRouterCoverageCSV(filepath.Join(*outDir, "futures_range_context_router_coverage.csv"), rangeContextRouterResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_router_rules.json"), rangeContextRouterResult.RuleRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextRouterRulesCSV(filepath.Join(*outDir, "futures_range_context_router_rules.csv"), rangeContextRouterResult.RuleRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_router_rows.json"), rangeContextRouterResult.Rows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextRouterRowsCSV(filepath.Join(*outDir, "futures_range_context_router_rows.csv"), rangeContextRouterResult.Rows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_router_cohorts.json"), rangeContextRouterResult.CohortRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextRouterCohortsCSV(filepath.Join(*outDir, "futures_range_context_router_cohorts.csv"), rangeContextRouterResult.CohortRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_router_rankings.json"), rangeContextRouterResult.RankingRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextRouterRankingsCSV(filepath.Join(*outDir, "futures_range_context_router_rankings.csv"), rangeContextRouterResult.RankingRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_router_summary.json"), rangeContextRouterResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextRouterSummaryCSV(filepath.Join(*outDir, "futures_range_context_router_summary.csv"), rangeContextRouterResult.SummaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_context_router_skips.json"), rangeContextRouterResult.SkipRows); err != nil {
			return err
		}
		if err := writeFuturesRangeContextRouterSkipsCSV(filepath.Join(*outDir, "futures_range_context_router_skips.csv"), rangeContextRouterResult.SkipRows); err != nil {
			return err
		}
		fmt.Printf("futures_range_context_router_audit source_rows=%d coverage_rows=%d rule_rows=%d router_rows=%d cohort_rows=%d ranking_rows=%d passing_cohorts=%d stop_state=%s\n",
			len(rangeContextRouterResult.SourceRows),
			len(rangeContextRouterResult.CoverageRows),
			len(rangeContextRouterResult.RuleRows),
			len(rangeContextRouterResult.Rows),
			len(rangeContextRouterResult.CohortRows),
			len(rangeContextRouterResult.RankingRows),
			rangeContextRouterResult.PassingCohorts,
			rangeContextRouterResult.StopState,
		)
	}
	if *futuresRangeRouterRotationPremiseAudit {
		premiseCfg := futuresRangeRouterRotationPremiseAuditConfigForRun()
		var err error
		rangeRouterRotationPremiseResult, err = lab.RunFuturesRangeRouterRotationPremiseAudit(candles, sourceManifest, premiseCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_sources.json"), rangeRouterRotationPremiseResult.SourceRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseSourcesCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_sources.csv"), rangeRouterRotationPremiseResult.SourceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_coverage.json"), rangeRouterRotationPremiseResult.CoverageRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseCoverageCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_coverage.csv"), rangeRouterRotationPremiseResult.CoverageRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_router_dependency.json"), rangeRouterRotationPremiseResult.DependencyRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseRouterDependencyCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_router_dependency.csv"), rangeRouterRotationPremiseResult.DependencyRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_context_segments.json"), rangeRouterRotationPremiseResult.SegmentRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseContextSegmentsCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_context_segments.csv"), rangeRouterRotationPremiseResult.SegmentRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_events.json"), rangeRouterRotationPremiseResult.EventRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseEventsCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_events.csv"), rangeRouterRotationPremiseResult.EventRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_outcomes.json"), rangeRouterRotationPremiseResult.OutcomeRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseOutcomesCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_outcomes.csv"), rangeRouterRotationPremiseResult.OutcomeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_cohorts.json"), rangeRouterRotationPremiseResult.CohortRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseCohortsCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_cohorts.csv"), rangeRouterRotationPremiseResult.CohortRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_rankings.json"), rangeRouterRotationPremiseResult.RankingRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseRankingsCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_rankings.csv"), rangeRouterRotationPremiseResult.RankingRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_summary.json"), rangeRouterRotationPremiseResult.SummaryRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseSummaryCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_summary.csv"), rangeRouterRotationPremiseResult.SummaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "futures_range_router_rotation_premise_skips.json"), rangeRouterRotationPremiseResult.SkipRows); err != nil {
			return err
		}
		if err := writeFuturesRangeRouterRotationPremiseSkipsCSV(filepath.Join(*outDir, "futures_range_router_rotation_premise_skips.csv"), rangeRouterRotationPremiseResult.SkipRows); err != nil {
			return err
		}
		fmt.Printf("futures_range_router_rotation_premise_audit source_rows=%d coverage_rows=%d router_dependency_rows=%d context_segments=%d events=%d outcomes=%d cohort_rows=%d ranking_rows=%d passing_cohorts=%d stop_state=%s\n",
			len(rangeRouterRotationPremiseResult.SourceRows),
			len(rangeRouterRotationPremiseResult.CoverageRows),
			len(rangeRouterRotationPremiseResult.DependencyRows),
			len(rangeRouterRotationPremiseResult.SegmentRows),
			len(rangeRouterRotationPremiseResult.EventRows),
			len(rangeRouterRotationPremiseResult.OutcomeRows),
			len(rangeRouterRotationPremiseResult.CohortRows),
			len(rangeRouterRotationPremiseResult.RankingRows),
			rangeRouterRotationPremiseResult.PassingCohorts,
			rangeRouterRotationPremiseResult.StopState,
		)
	}
	var srRows []lab.SRAuditRow
	srCfg := lab.DefaultSRAuditConfig()
	if *srAudit || *srBoundaryAudit || *srBoundaryInspect || *srRejectionTimingAudit || *srConfirmationTimingAudit || *srFalseBreakReclaimTimingAudit {
		var err error
		srRows, err = lab.RunSRAudit(candles, srCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
	}
	if *srAudit {
		if err := writeJSON(filepath.Join(*outDir, "sr_touch_audit.json"), srRows); err != nil {
			return err
		}
		if err := writeSRTouchAuditCSV(filepath.Join(*outDir, "sr_touch_audit.csv"), srRows); err != nil {
			return err
		}
		nearSupportRows := 0
		nearResistanceRows := 0
		for _, row := range srRows {
			if row.NearSupport {
				nearSupportRows++
			}
			if row.NearResistance {
				nearResistanceRows++
			}
		}
		warmupBars := 0
		if len(srRows) > 0 {
			warmupBars = srRows[0].WarmupBars
		}
		fmt.Printf("sr_audit rows=%d lookback_bars=%d warmup_bars=%d near_support_rows=%d near_resistance_rows=%d\n",
			len(srRows),
			srCfg.LookbackBars,
			warmupBars,
			nearSupportRows,
			nearResistanceRows,
		)
	}
	if *srBoundaryAudit || *srBoundaryInspect {
		boundaryCfg := lab.DefaultSRBoundaryAuditConfig()
		events, qualityRows, err := lab.RunSRBoundaryAudit(candles, srRows, boundaryCfg)
		if err != nil {
			return err
		}
		if *srBoundaryAudit {
			if err := writeJSON(filepath.Join(*outDir, "sr_boundary_events.json"), events); err != nil {
				return err
			}
			if err := writeSRBoundaryEventsCSV(filepath.Join(*outDir, "sr_boundary_events.csv"), events); err != nil {
				return err
			}
			if err := writeJSON(filepath.Join(*outDir, "sr_boundary_quality.json"), qualityRows); err != nil {
				return err
			}
			if err := writeSRBoundaryQualityCSV(filepath.Join(*outDir, "sr_boundary_quality.csv"), qualityRows); err != nil {
				return err
			}
			fmt.Printf("sr_boundary_audit events=%d summary_rows=%d horizons=%s detector_active_only=%t\n",
				len(events),
				len(qualityRows),
				formatIntSlice(boundaryCfg.HorizonsBars),
				boundaryCfg.DetectorActiveOnly,
			)
		}
		if *srBoundaryInspect {
			comparisonRows := lab.SummarizeSRBoundaryCandidateComparison(events)
			if err := writeJSON(filepath.Join(*outDir, "sr_boundary_candidate_comparison.json"), comparisonRows); err != nil {
				return err
			}
			if err := writeSRBoundaryCandidateComparisonCSV(filepath.Join(*outDir, "sr_boundary_candidate_comparison.csv"), comparisonRows); err != nil {
				return err
			}
			fmt.Printf("sr_boundary_inspect events=%d comparison_rows=%d horizons=%s detector_active_only=%t\n",
				len(events),
				len(comparisonRows),
				formatIntSlice(boundaryCfg.HorizonsBars),
				boundaryCfg.DetectorActiveOnly,
			)
		}
	}
	if *srRejectionTimingAudit {
		timingCfg := lab.DefaultSRRejectionTimingAuditConfig()
		candidateRows, summaryRows, err := lab.RunSRRejectionTimingAudit(candles, srRows, timingCfg)
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_rejection_timing_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeSRRejectionTimingCandidatesCSV(filepath.Join(*outDir, "sr_rejection_timing_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_rejection_timing_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeSRRejectionTimingSummaryCSV(filepath.Join(*outDir, "sr_rejection_timing_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("sr_rejection_timing_audit candidate_rows=%d summary_rows=%d horizons=%s detector_active_only=%t\n",
			len(candidateRows),
			len(summaryRows),
			formatIntSlice(timingCfg.HorizonsBars),
			timingCfg.DetectorActiveOnly,
		)
	}
	if *srConfirmationTimingAudit {
		confirmationCfg := lab.DefaultSRConfirmationTimingAuditConfig()
		candidateRows, summaryRows, err := lab.RunSRConfirmationTimingAudit(candles, srRows, confirmationCfg)
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_confirmation_timing_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeSRConfirmationTimingCandidatesCSV(filepath.Join(*outDir, "sr_confirmation_timing_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_confirmation_timing_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeSRConfirmationTimingSummaryCSV(filepath.Join(*outDir, "sr_confirmation_timing_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("sr_confirmation_timing_audit candidate_rows=%d summary_rows=%d delays=%s horizons=%s detector_active_only=%t\n",
			len(candidateRows),
			len(summaryRows),
			formatIntSlice(confirmationCfg.ConfirmationDelayBars),
			formatIntSlice(confirmationCfg.HorizonsBars),
			confirmationCfg.DetectorActiveOnly,
		)
	}
	if *srFalseBreakReclaimTimingAudit {
		falseBreakCfg := lab.DefaultSRFalseBreakReclaimTimingAuditConfig()
		candidateRows, summaryRows, err := lab.RunSRFalseBreakReclaimTimingAudit(candles, srRows, falseBreakCfg)
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_false_break_reclaim_timing_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeSRFalseBreakReclaimTimingCandidatesCSV(filepath.Join(*outDir, "sr_false_break_reclaim_timing_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_false_break_reclaim_timing_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeSRFalseBreakReclaimTimingSummaryCSV(filepath.Join(*outDir, "sr_false_break_reclaim_timing_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("sr_false_break_reclaim_timing_audit candidate_rows=%d summary_rows=%d max_break_delay=%d max_reclaim_delay=%d horizons=%s detector_active_only=%t\n",
			len(candidateRows),
			len(summaryRows),
			falseBreakCfg.MaxBreakDelayBars,
			falseBreakCfg.MaxReclaimDelayBars,
			formatIntSlice(falseBreakCfg.HorizonsBars),
			falseBreakCfg.DetectorActiveOnly,
		)
	}
	if *compressionBreakoutAudit {
		breakoutCfg := lab.DefaultCompressionBreakoutAuditConfig()
		candidateRows, summaryRows, err := lab.RunCompressionBreakoutAudit(candles, breakoutCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "compression_breakout_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeCompressionBreakoutCandidatesCSV(filepath.Join(*outDir, "compression_breakout_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "compression_breakout_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeCompressionBreakoutSummaryCSV(filepath.Join(*outDir, "compression_breakout_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("compression_breakout_audit candidate_rows=%d summary_rows=%d max_breakout_delay=%d horizons=%s detector_profile_id=%s\n",
			len(candidateRows),
			len(summaryRows),
			breakoutCfg.MaxBreakoutDelayBars,
			formatIntSlice(breakoutCfg.HorizonsBars),
			breakoutCfg.DetectorProfileID,
		)
	}
	if *rangeRegimeDurabilityAudit {
		durabilityCfg := lab.DefaultRangeRegimeDurabilityAuditConfig()
		episodeRows, summaryRows, err := lab.RunRangeRegimeDurabilityAudit(candles, durabilityCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "range_regime_durability_episodes.json"), episodeRows); err != nil {
			return err
		}
		if err := writeRangeRegimeDurabilityEpisodesCSV(filepath.Join(*outDir, "range_regime_durability_episodes.csv"), episodeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "range_regime_durability_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeRangeRegimeDurabilitySummaryCSV(filepath.Join(*outDir, "range_regime_durability_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("range_regime_durability_audit episode_rows=%d summary_rows=%d quick_invalidation_bars=%d horizons=%s detector_profile_id=%s\n",
			len(episodeRows),
			len(summaryRows),
			durabilityCfg.QuickInvalidationBars,
			formatIntSlice(durabilityCfg.HorizonsBars),
			durabilityCfg.DetectorProfileID,
		)
	}
	if *detector || *detectorSweep || *detectorDurabilitySweep || *detectorContextRefinementAudit || *holdInsideDirectionalEdgeAudit || *holdInsideMidlineTransitionAudit || *holdInsideMidlineReactionAudit || *holdInsideMidlineTouchPrototype {
		if *detectorLookbackDays <= 0 {
			return fmt.Errorf("detector lookback days must be positive")
		}
	}
	if *detector {
		if *detectorPercentile <= 0 || *detectorPercentile >= 1 {
			return fmt.Errorf("detector percentile must be between 0 and 1")
		}
		if *detectorMinConsecutiveBars <= 0 {
			return fmt.Errorf("detector min consecutive bars must be positive")
		}

		cfg := lab.DefaultCompressionRangeDetectorConfig()
		cfg.LookbackDays = *detectorLookbackDays
		cfg.Percentile = *detectorPercentile
		cfg.MinConsecutiveBars = *detectorMinConsecutiveBars
		cfg.UseBollinger = *detectorUseBollinger
		cfg.UseADX = *detectorUseADX

		rangeDetector := lab.CompressionRangeDetector{Config: cfg}
		classifications, err := rangeDetector.Classify(candles)
		if err != nil {
			return err
		}
		dutyRows, episodes := lab.SummarizeDetectorSplits(candles, classifications, lab.DefaultSplits())
		if err := writeJSON(filepath.Join(*outDir, "detector_duty_cycle.json"), dutyRows); err != nil {
			return err
		}
		if err := writeDetectorDutyCycleCSV(filepath.Join(*outDir, "detector_duty_cycle.csv"), dutyRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "range_episodes.json"), episodes); err != nil {
			return err
		}
		if err := writeRangeEpisodesCSV(filepath.Join(*outDir, "range_episodes.csv"), episodes); err != nil {
			return err
		}
		for _, row := range dutyRows {
			if row.Split == "full_2021_2026" {
				fmt.Printf("detector=%s active_bars=%d total_bars=%d duty_cycle=%.4f episodes=%d\n",
					rangeDetector.Name(), row.ActiveBars, row.TotalBars, row.DutyCycle, row.Episodes)
				break
			}
		}
	}
	if *detectorSweep {
		cfg := lab.DefaultCompressionRangeDetectorConfig()
		cfg.LookbackDays = *detectorLookbackDays

		sweepRows, err := lab.RunDetectorSweep(candles, cfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_sweep.json"), sweepRows); err != nil {
			return err
		}
		if err := writeDetectorSweepCSV(filepath.Join(*outDir, "detector_sweep.csv"), sweepRows); err != nil {
			return err
		}
		for _, row := range sweepRows {
			if row.IsBalancedBaseline && row.Split == "full_2021_2026" {
				fmt.Printf("detector_sweep profiles=%d rows=%d baseline_active_bars=%d baseline_total_bars=%d baseline_duty_cycle=%.4f baseline_episodes=%d\n",
					len(lab.DefaultDetectorSweepProfiles(*detectorLookbackDays)),
					len(sweepRows),
					row.ActiveBars,
					row.TotalBars,
					row.DutyCycle,
					row.Episodes,
				)
				break
			}
		}
	}
	if *detectorDurabilitySweep {
		detectorCfg := lab.DefaultCompressionRangeDetectorConfig()
		detectorCfg.LookbackDays = *detectorLookbackDays
		durabilityCfg := lab.DefaultRangeRegimeDurabilityAuditConfig()

		broadRows, sliceRows, stabilityRows, err := lab.RunDetectorDurabilitySweep(candles, detectorCfg, durabilityCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_durability_sweep.json"), broadRows); err != nil {
			return err
		}
		if err := writeDetectorDurabilitySweepCSV(filepath.Join(*outDir, "detector_durability_sweep.csv"), broadRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_durability_slices.json"), sliceRows); err != nil {
			return err
		}
		if err := writeDetectorDurabilitySlicesCSV(filepath.Join(*outDir, "detector_durability_slices.csv"), sliceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_durability_stability.json"), stabilityRows); err != nil {
			return err
		}
		if err := writeDetectorDurabilityStabilityCSV(filepath.Join(*outDir, "detector_durability_stability.csv"), stabilityRows); err != nil {
			return err
		}
		fmt.Printf("detector_durability_sweep profiles=%d broad_rows=%d slice_rows=%d stability_rows=%d quick_invalidation_bars=%d horizons=%s\n",
			len(lab.DefaultDetectorSweepProfiles(*detectorLookbackDays)),
			len(broadRows),
			len(sliceRows),
			len(stabilityRows),
			durabilityCfg.QuickInvalidationBars,
			formatIntSlice(durabilityCfg.HorizonsBars),
		)
	}
	if *detectorContextRefinementAudit {
		detectorCfg := lab.DefaultCompressionRangeDetectorConfig()
		detectorCfg.LookbackDays = *detectorLookbackDays
		refinementCfg := lab.DefaultDetectorContextRefinementAuditConfig()
		refinementCfg.Profiles = lab.DefaultDetectorContextRefinementProfiles(*detectorLookbackDays)

		candidateRows, summaryRows, stabilityRows, err := lab.RunDetectorContextRefinementAudit(candles, detectorCfg, refinementCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_context_refinement_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeDetectorContextRefinementCandidatesCSV(filepath.Join(*outDir, "detector_context_refinement_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_context_refinement_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeDetectorContextRefinementSummaryCSV(filepath.Join(*outDir, "detector_context_refinement_summary.csv"), summaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_context_refinement_stability.json"), stabilityRows); err != nil {
			return err
		}
		if err := writeDetectorContextRefinementStabilityCSV(filepath.Join(*outDir, "detector_context_refinement_stability.csv"), stabilityRows); err != nil {
			return err
		}
		fmt.Printf("detector_context_refinement_audit profiles=%d rules=%d candidate_rows=%d summary_rows=%d stability_rows=%d quick_invalidation_bars=%d horizons=%s\n",
			len(refinementCfg.Profiles),
			len(refinementCfg.ContextRules),
			len(candidateRows),
			len(summaryRows),
			len(stabilityRows),
			refinementCfg.QuickInvalidationBars,
			formatIntSlice(refinementCfg.HorizonsBars),
		)
	}
	if *holdInsideDirectionalEdgeAudit {
		detectorCfg := lab.DefaultCompressionRangeDetectorConfig()
		detectorCfg.LookbackDays = *detectorLookbackDays
		edgeDefaults := lab.DefaultHoldInsideDirectionalEdgeAuditConfig()

		candidateRows, summaryRows, stabilityRows, err := lab.RunHoldInsideDirectionalEdgeAudit(candles, detectorCfg, lab.HoldInsideDirectionalEdgeAuditConfig{}, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_directional_edge_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeHoldInsideDirectionalEdgeCandidatesCSV(filepath.Join(*outDir, "hold_inside_directional_edge_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_directional_edge_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeHoldInsideDirectionalEdgeSummaryCSV(filepath.Join(*outDir, "hold_inside_directional_edge_summary.csv"), summaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_directional_edge_stability.json"), stabilityRows); err != nil {
			return err
		}
		if err := writeHoldInsideDirectionalEdgeStabilityCSV(filepath.Join(*outDir, "hold_inside_directional_edge_stability.csv"), stabilityRows); err != nil {
			return err
		}
		fmt.Printf("hold_inside_directional_edge_audit profiles=%d rules=%d paper_sides=%d candidate_rows=%d summary_rows=%d stability_rows=%d quick_invalidation_bars=%d horizons=%s\n",
			len(edgeDefaults.Profiles),
			len(edgeDefaults.ContextRules),
			2,
			len(candidateRows),
			len(summaryRows),
			len(stabilityRows),
			edgeDefaults.QuickInvalidationBars,
			formatIntSlice(edgeDefaults.HorizonsBars),
		)
	}
	if *holdInsideMidlineTransitionAudit {
		detectorCfg := lab.DefaultCompressionRangeDetectorConfig()
		detectorCfg.LookbackDays = *detectorLookbackDays
		midlineDefaults := lab.DefaultHoldInsideMidlineTransitionAuditConfig()

		candidateRows, summaryRows, stabilityRows, err := lab.RunHoldInsideMidlineTransitionAudit(candles, detectorCfg, lab.HoldInsideMidlineTransitionAuditConfig{}, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_transition_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineTransitionCandidatesCSV(filepath.Join(*outDir, "hold_inside_midline_transition_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_transition_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineTransitionSummaryCSV(filepath.Join(*outDir, "hold_inside_midline_transition_summary.csv"), summaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_transition_stability.json"), stabilityRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineTransitionStabilityCSV(filepath.Join(*outDir, "hold_inside_midline_transition_stability.csv"), stabilityRows); err != nil {
			return err
		}
		fmt.Printf("hold_inside_midline_transition_audit profiles=%d rules=%d candidate_rows=%d summary_rows=%d stability_rows=%d quick_invalidation_bars=%d horizons=%s\n",
			len(midlineDefaults.Profiles),
			len(midlineDefaults.ContextRules),
			len(candidateRows),
			len(summaryRows),
			len(stabilityRows),
			midlineDefaults.QuickInvalidationBars,
			formatIntSlice(midlineDefaults.HorizonsBars),
		)
	}
	if *holdInsideMidlineReactionAudit {
		detectorCfg := lab.DefaultCompressionRangeDetectorConfig()
		detectorCfg.LookbackDays = *detectorLookbackDays
		reactionDefaults := lab.DefaultHoldInsideMidlineReactionAuditConfig()

		candidateRows, funnelRows, summaryRows, stabilityRows, err := lab.RunHoldInsideMidlineReactionAudit(candles, detectorCfg, lab.HoldInsideMidlineReactionAuditConfig{}, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_reaction_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineReactionCandidatesCSV(filepath.Join(*outDir, "hold_inside_midline_reaction_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_reaction_funnel_summary.json"), funnelRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineReactionFunnelSummaryCSV(filepath.Join(*outDir, "hold_inside_midline_reaction_funnel_summary.csv"), funnelRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_reaction_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineReactionSummaryCSV(filepath.Join(*outDir, "hold_inside_midline_reaction_summary.csv"), summaryRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "hold_inside_midline_reaction_stability.json"), stabilityRows); err != nil {
			return err
		}
		if err := writeHoldInsideMidlineReactionStabilityCSV(filepath.Join(*outDir, "hold_inside_midline_reaction_stability.csv"), stabilityRows); err != nil {
			return err
		}
		fmt.Printf("hold_inside_midline_reaction_audit profiles=%d rules=%d event_types=%d candidate_rows=%d funnel_rows=%d summary_rows=%d stability_rows=%d max_midline_event_delay_bars=%d quick_invalidation_bars=%d horizons=%s\n",
			len(reactionDefaults.Profiles),
			len(reactionDefaults.ContextRules),
			2,
			len(candidateRows),
			len(funnelRows),
			len(summaryRows),
			len(stabilityRows),
			reactionDefaults.MaxMidlineEventDelayBars,
			reactionDefaults.QuickInvalidationBars,
			formatIntSlice(reactionDefaults.HorizonsBars),
		)
	}

	first := candles[0].OpenTime.Format(time.RFC3339)
	last := candles[len(candles)-1].CloseTime.Format(time.RFC3339)
	fmt.Printf("loaded %d candles from %s to %s\n", len(candles), first, last)
	fmt.Printf("strategy=%s trades=%d output=%s\n", strategyName, len(result.Trades), *outDir)
	return nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func writeJSONTaggedCSV[T any](path string, rows []T) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	rowType := reflect.TypeOf((*T)(nil)).Elem()
	if rowType.Kind() == reflect.Pointer {
		rowType = rowType.Elem()
	}
	fields := csvTaggedFields(rowType, nil)
	headers := make([]string, 0, len(fields))
	for _, field := range fields {
		headers = append(headers, field.header)
	}
	if err := w.Write(headers); err != nil {
		return err
	}
	for _, row := range rows {
		value := reflect.ValueOf(row)
		if value.Kind() == reflect.Pointer {
			if value.IsNil() {
				continue
			}
			value = value.Elem()
		}
		record := make([]string, 0, len(fields))
		for _, field := range fields {
			record = append(record, csvScalar(csvFieldByIndex(value, field.index)))
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return w.Error()
}

type csvTaggedField struct {
	header string
	index  []int
}

func csvTaggedFields(rowType reflect.Type, prefix []int) []csvTaggedField {
	fields := []csvTaggedField{}
	if rowType.Kind() == reflect.Pointer {
		rowType = rowType.Elem()
	}
	if rowType.Kind() != reflect.Struct {
		return fields
	}
	for i := 0; i < rowType.NumField(); i++ {
		field := rowType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		index := append(append([]int(nil), prefix...), i)
		header := strings.Split(field.Tag.Get("json"), ",")[0]
		if header == "-" {
			continue
		}
		fieldType := field.Type
		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}
		if field.Anonymous && fieldType.Kind() == reflect.Struct && header == "" {
			fields = append(fields, csvTaggedFields(fieldType, index)...)
			continue
		}
		if header == "" {
			header = field.Name
		}
		fields = append(fields, csvTaggedField{header: header, index: index})
	}
	return fields
}

func csvFieldByIndex(value reflect.Value, index []int) reflect.Value {
	for _, part := range index {
		if value.Kind() == reflect.Pointer {
			if value.IsNil() {
				return reflect.Value{}
			}
			value = value.Elem()
		}
		if value.Kind() != reflect.Struct || part >= value.NumField() {
			return reflect.Value{}
		}
		value = value.Field(part)
	}
	return value
}

func csvScalar(value reflect.Value) string {
	if !value.IsValid() {
		return ""
	}
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return formatFloat(value.Float())
	default:
		if value.CanInterface() {
			return fmt.Sprint(value.Interface())
		}
		return ""
	}
}

func writeSummaryCSV(path string, rows []lab.SummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"total_trades",
		"wins",
		"losses",
		"win_rate",
		"gross_pnl",
		"net_pnl",
		"total_costs",
		"profit_factor",
		"gross_profit_factor",
		"max_drawdown",
		"expectancy",
		"avg_hold_bars",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.TotalTrades),
			strconv.Itoa(row.Wins),
			strconv.Itoa(row.Losses),
			formatFloat(row.WinRate),
			formatFloat(row.GrossPnL),
			formatFloat(row.NetPnL),
			formatFloat(row.TotalCosts),
			formatFloat(row.ProfitFactor),
			formatFloat(row.GrossProfitFactor),
			formatFloat(row.MaxDrawdown),
			formatFloat(row.Expectancy),
			formatFloat(row.AvgHoldBars),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorDutyCycleCSV(path string, rows []lab.DetectorDutyCycleRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"active_bars",
		"total_bars",
		"duty_cycle",
		"episodes",
		"avg_episode_length",
		"median_episode_length",
		"longest_episode_length",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			strconv.Itoa(row.ActiveBars),
			strconv.Itoa(row.TotalBars),
			formatFloat(row.DutyCycle),
			strconv.Itoa(row.Episodes),
			formatFloat(row.AvgEpisodeLength),
			formatFloat(row.MedianEpisodeLength),
			strconv.Itoa(row.LongestEpisodeLength),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeRangeEpisodesCSV(path string, episodes []lab.RangeEpisode) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"start_index",
		"end_index",
		"start_time",
		"end_time",
		"length_bars",
	}); err != nil {
		return err
	}
	for _, episode := range episodes {
		if err := w.Write([]string{
			episode.Split,
			strconv.Itoa(episode.StartIndex),
			strconv.Itoa(episode.EndIndex),
			episode.StartTime,
			episode.EndTime,
			strconv.Itoa(episode.LengthBars),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorSweepCSV(path string, rows []lab.DetectorSweepRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"split",
		"active_bars",
		"total_bars",
		"duty_cycle",
		"episodes",
		"avg_episode_length",
		"median_episode_length",
		"longest_episode_length",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.Split,
			strconv.Itoa(row.ActiveBars),
			strconv.Itoa(row.TotalBars),
			formatFloat(row.DutyCycle),
			strconv.Itoa(row.Episodes),
			formatFloat(row.AvgEpisodeLength),
			formatFloat(row.MedianEpisodeLength),
			strconv.Itoa(row.LongestEpisodeLength),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorDurabilitySweepCSV(path string, rows []lab.DetectorDurabilitySweepRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"split",
		"horizon_bars",
		"active_bars",
		"total_bars",
		"duty_cycle",
		"detector_episodes",
		"avg_detector_episode_length",
		"median_detector_episode_length",
		"longest_detector_episode_length",
		"durability_episode_count",
		"avg_raw_length_bars",
		"avg_active_length_bars",
		"avg_episode_width_pct",
		"avg_normalized_atr",
		"avg_end_normalized_atr",
		"avg_width_to_atr_ratio",
		"label_reentered_range_count",
		"label_persisted_inside_range_count",
		"label_quick_invalidated_count",
		"label_invalidated_up_count",
		"label_invalidated_down_count",
		"label_chopped_count",
		"label_trended_up_count",
		"label_trended_down_count",
		"label_reentered_range_rate",
		"label_persisted_inside_range_rate",
		"label_quick_invalidated_rate",
		"label_invalidated_up_rate",
		"label_invalidated_down_rate",
		"label_chopped_rate",
		"label_trended_up_rate",
		"label_trended_down_rate",
		"label_avg_close_drift_pct",
		"label_avg_max_up_move_pct",
		"label_avg_max_down_move_pct",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.Split,
			strconv.Itoa(row.HorizonBars),
			strconv.Itoa(row.ActiveBars),
			strconv.Itoa(row.TotalBars),
			formatFloat(row.DutyCycle),
			strconv.Itoa(row.DetectorEpisodes),
			formatFloat(row.AvgDetectorEpisodeLength),
			formatFloat(row.MedianDetectorEpisodeLength),
			strconv.Itoa(row.LongestDetectorEpisodeLength),
			strconv.Itoa(row.DurabilityEpisodeCount),
			formatFloat(row.AvgRawLengthBars),
			formatFloat(row.AvgActiveLengthBars),
			formatFloat(row.AvgEpisodeWidthPct),
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.AvgEndNormalizedATR),
			formatFloat(row.AvgWidthToATRRatio),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelPersistedInsideRangeCount),
			strconv.Itoa(row.LabelQuickInvalidatedCount),
			strconv.Itoa(row.LabelInvalidatedUpCount),
			strconv.Itoa(row.LabelInvalidatedDownCount),
			strconv.Itoa(row.LabelChoppedCount),
			strconv.Itoa(row.LabelTrendedUpCount),
			strconv.Itoa(row.LabelTrendedDownCount),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelPersistedInsideRangeRate),
			formatFloat(row.LabelQuickInvalidatedRate),
			formatFloat(row.LabelInvalidatedUpRate),
			formatFloat(row.LabelInvalidatedDownRate),
			formatFloat(row.LabelChoppedRate),
			formatFloat(row.LabelTrendedUpRate),
			formatFloat(row.LabelTrendedDownRate),
			formatFloat(row.LabelAvgCloseDriftPct),
			formatFloat(row.LabelAvgMaxUpMovePct),
			formatFloat(row.LabelAvgMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorDurabilitySlicesCSV(path string, rows []lab.DetectorDurabilitySliceRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"split",
		"horizon_bars",
		"raw_length_bucket",
		"active_length_bucket",
		"episode_width_bucket",
		"width_to_atr_bucket",
		"episode_count",
		"avg_raw_length_bars",
		"avg_active_length_bars",
		"avg_episode_width_pct",
		"avg_normalized_atr",
		"avg_end_normalized_atr",
		"avg_width_to_atr_ratio",
		"label_reentered_range_count",
		"label_persisted_inside_range_count",
		"label_quick_invalidated_count",
		"label_invalidated_up_count",
		"label_invalidated_down_count",
		"label_chopped_count",
		"label_trended_up_count",
		"label_trended_down_count",
		"label_reentered_range_rate",
		"label_persisted_inside_range_rate",
		"label_quick_invalidated_rate",
		"label_invalidated_up_rate",
		"label_invalidated_down_rate",
		"label_chopped_rate",
		"label_trended_up_rate",
		"label_trended_down_rate",
		"label_avg_close_drift_pct",
		"label_avg_max_up_move_pct",
		"label_avg_max_down_move_pct",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.Split,
			strconv.Itoa(row.HorizonBars),
			row.RawLengthBucket,
			row.ActiveLengthBucket,
			row.EpisodeWidthBucket,
			row.WidthToATRBucket,
			strconv.Itoa(row.EpisodeCount),
			formatFloat(row.AvgRawLengthBars),
			formatFloat(row.AvgActiveLengthBars),
			formatFloat(row.AvgEpisodeWidthPct),
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.AvgEndNormalizedATR),
			formatFloat(row.AvgWidthToATRRatio),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelPersistedInsideRangeCount),
			strconv.Itoa(row.LabelQuickInvalidatedCount),
			strconv.Itoa(row.LabelInvalidatedUpCount),
			strconv.Itoa(row.LabelInvalidatedDownCount),
			strconv.Itoa(row.LabelChoppedCount),
			strconv.Itoa(row.LabelTrendedUpCount),
			strconv.Itoa(row.LabelTrendedDownCount),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelPersistedInsideRangeRate),
			formatFloat(row.LabelQuickInvalidatedRate),
			formatFloat(row.LabelInvalidatedUpRate),
			formatFloat(row.LabelInvalidatedDownRate),
			formatFloat(row.LabelChoppedRate),
			formatFloat(row.LabelTrendedUpRate),
			formatFloat(row.LabelTrendedDownRate),
			formatFloat(row.LabelAvgCloseDriftPct),
			formatFloat(row.LabelAvgMaxUpMovePct),
			formatFloat(row.LabelAvgMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorDurabilityStabilityCSV(path string, rows []lab.DetectorDurabilityStabilityRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"horizon_bars",
		"period_splits",
		"period_episode_count",
		"episode_count_min",
		"episode_count_max",
		"episode_count_delta",
		"duty_cycle_min",
		"duty_cycle_max",
		"duty_cycle_delta",
		"label_reentered_range_rate_min",
		"label_reentered_range_rate_max",
		"label_reentered_range_rate_delta",
		"label_persisted_inside_range_rate_min",
		"label_persisted_inside_range_rate_max",
		"label_persisted_inside_range_rate_delta",
		"label_quick_invalidated_rate_min",
		"label_quick_invalidated_rate_max",
		"label_quick_invalidated_rate_delta",
		"label_chopped_rate_min",
		"label_chopped_rate_max",
		"label_chopped_rate_delta",
		"label_trended_rate_min",
		"label_trended_rate_max",
		"label_trended_rate_delta",
		"label_avg_close_drift_pct_min",
		"label_avg_close_drift_pct_max",
		"label_avg_close_drift_pct_delta",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			strconv.Itoa(row.HorizonBars),
			strconv.Itoa(row.PeriodSplits),
			strconv.Itoa(row.PeriodEpisodeCount),
			strconv.Itoa(row.EpisodeCountMin),
			strconv.Itoa(row.EpisodeCountMax),
			strconv.Itoa(row.EpisodeCountDelta),
			formatFloat(row.DutyCycleMin),
			formatFloat(row.DutyCycleMax),
			formatFloat(row.DutyCycleDelta),
			formatFloat(row.LabelReenteredRangeRateMin),
			formatFloat(row.LabelReenteredRangeRateMax),
			formatFloat(row.LabelReenteredRangeRateDelta),
			formatFloat(row.LabelPersistedInsideRangeRateMin),
			formatFloat(row.LabelPersistedInsideRangeRateMax),
			formatFloat(row.LabelPersistedInsideRangeRateDelta),
			formatFloat(row.LabelQuickInvalidatedRateMin),
			formatFloat(row.LabelQuickInvalidatedRateMax),
			formatFloat(row.LabelQuickInvalidatedRateDelta),
			formatFloat(row.LabelChoppedRateMin),
			formatFloat(row.LabelChoppedRateMax),
			formatFloat(row.LabelChoppedRateDelta),
			formatFloat(row.LabelTrendedRateMin),
			formatFloat(row.LabelTrendedRateMax),
			formatFloat(row.LabelTrendedRateDelta),
			formatFloat(row.LabelAvgCloseDriftPctMin),
			formatFloat(row.LabelAvgCloseDriftPctMax),
			formatFloat(row.LabelAvgCloseDriftPctDelta),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorContextRefinementCandidatesCSV(path string, rows []lab.DetectorContextRefinementCandidateRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"context_rule",
		"hold_bars",
		"require_mid_50",
		"split",
		"source_episode_id",
		"episode_start_index",
		"episode_end_index",
		"episode_start_time",
		"episode_end_time",
		"raw_length_bars",
		"active_length_bars",
		"raw_length_bucket",
		"active_length_bucket",
		"episode_high",
		"episode_low",
		"episode_end_close",
		"episode_width_pct",
		"episode_width_bucket",
		"avg_normalized_atr",
		"end_normalized_atr",
		"width_to_atr_ratio",
		"width_to_atr_bucket",
		"decision_index",
		"decision_time",
		"decision_close",
		"decision_close_position",
		"decision_close_position_bucket",
		"horizon_bars",
		"label_window_start_index",
		"label_window_end_index",
		"label_window_start_time",
		"label_window_end_time",
		"label_reentered_range",
		"label_persisted_inside_range",
		"label_quick_invalidated",
		"label_invalidated_up",
		"label_invalidated_down",
		"label_chopped",
		"label_trended_up",
		"label_trended_down",
		"label_close_drift_pct",
		"label_max_up_move_pct",
		"label_max_down_move_pct",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.ContextRule,
			strconv.Itoa(row.HoldBars),
			strconv.FormatBool(row.RequireMid50),
			row.Split,
			strconv.Itoa(row.SourceEpisodeID),
			strconv.Itoa(row.EpisodeStartIndex),
			strconv.Itoa(row.EpisodeEndIndex),
			row.EpisodeStartTime,
			row.EpisodeEndTime,
			strconv.Itoa(row.RawLengthBars),
			strconv.Itoa(row.ActiveLengthBars),
			row.RawLengthBucket,
			row.ActiveLengthBucket,
			formatFloat(row.EpisodeHigh),
			formatFloat(row.EpisodeLow),
			formatFloat(row.EpisodeEndClose),
			formatFloat(row.EpisodeWidthPct),
			row.EpisodeWidthBucket,
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.EndNormalizedATR),
			formatFloat(row.WidthToATRRatio),
			row.WidthToATRBucket,
			strconv.Itoa(row.DecisionIndex),
			row.DecisionTime,
			formatFloat(row.DecisionClose),
			formatFloat(row.DecisionClosePosition),
			row.DecisionClosePositionBucket,
			strconv.Itoa(row.HorizonBars),
			strconv.Itoa(row.LabelWindowStartIndex),
			strconv.Itoa(row.LabelWindowEndIndex),
			row.LabelWindowStartTime,
			row.LabelWindowEndTime,
			strconv.FormatBool(row.LabelReenteredRange),
			strconv.FormatBool(row.LabelPersistedInsideRange),
			strconv.FormatBool(row.LabelQuickInvalidated),
			strconv.FormatBool(row.LabelInvalidatedUp),
			strconv.FormatBool(row.LabelInvalidatedDown),
			strconv.FormatBool(row.LabelChopped),
			strconv.FormatBool(row.LabelTrendedUp),
			strconv.FormatBool(row.LabelTrendedDown),
			formatFloat(row.LabelCloseDriftPct),
			formatFloat(row.LabelMaxUpMovePct),
			formatFloat(row.LabelMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorContextRefinementSummaryCSV(path string, rows []lab.DetectorContextRefinementSummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"context_rule",
		"hold_bars",
		"require_mid_50",
		"split",
		"horizon_bars",
		"source_episode_count",
		"candidate_count",
		"candidate_rate",
		"avg_raw_length_bars",
		"avg_active_length_bars",
		"avg_episode_width_pct",
		"avg_normalized_atr",
		"avg_end_normalized_atr",
		"avg_width_to_atr_ratio",
		"avg_decision_close_position",
		"label_reentered_range_count",
		"label_persisted_inside_range_count",
		"label_quick_invalidated_count",
		"label_invalidated_up_count",
		"label_invalidated_down_count",
		"label_chopped_count",
		"label_trended_up_count",
		"label_trended_down_count",
		"label_reentered_range_rate",
		"label_persisted_inside_range_rate",
		"label_quick_invalidated_rate",
		"label_invalidated_up_rate",
		"label_invalidated_down_rate",
		"label_chopped_rate",
		"label_trended_up_rate",
		"label_trended_down_rate",
		"label_avg_close_drift_pct",
		"label_avg_max_up_move_pct",
		"label_avg_max_down_move_pct",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.ContextRule,
			strconv.Itoa(row.HoldBars),
			strconv.FormatBool(row.RequireMid50),
			row.Split,
			strconv.Itoa(row.HorizonBars),
			strconv.Itoa(row.SourceEpisodeCount),
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.CandidateRate),
			formatFloat(row.AvgRawLengthBars),
			formatFloat(row.AvgActiveLengthBars),
			formatFloat(row.AvgEpisodeWidthPct),
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.AvgEndNormalizedATR),
			formatFloat(row.AvgWidthToATRRatio),
			formatFloat(row.AvgDecisionClosePosition),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelPersistedInsideRangeCount),
			strconv.Itoa(row.LabelQuickInvalidatedCount),
			strconv.Itoa(row.LabelInvalidatedUpCount),
			strconv.Itoa(row.LabelInvalidatedDownCount),
			strconv.Itoa(row.LabelChoppedCount),
			strconv.Itoa(row.LabelTrendedUpCount),
			strconv.Itoa(row.LabelTrendedDownCount),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelPersistedInsideRangeRate),
			formatFloat(row.LabelQuickInvalidatedRate),
			formatFloat(row.LabelInvalidatedUpRate),
			formatFloat(row.LabelInvalidatedDownRate),
			formatFloat(row.LabelChoppedRate),
			formatFloat(row.LabelTrendedUpRate),
			formatFloat(row.LabelTrendedDownRate),
			formatFloat(row.LabelAvgCloseDriftPct),
			formatFloat(row.LabelAvgMaxUpMovePct),
			formatFloat(row.LabelAvgMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorContextRefinementStabilityCSV(path string, rows []lab.DetectorContextRefinementStabilityRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"context_rule",
		"hold_bars",
		"require_mid_50",
		"horizon_bars",
		"period_splits",
		"source_episode_count",
		"source_episode_count_min",
		"source_episode_count_max",
		"source_episode_count_delta",
		"candidate_count",
		"candidate_count_min",
		"candidate_count_max",
		"candidate_count_delta",
		"candidate_rate_min",
		"candidate_rate_max",
		"candidate_rate_delta",
		"label_persisted_inside_range_rate_min",
		"label_persisted_inside_range_rate_max",
		"label_persisted_inside_range_rate_delta",
		"label_quick_invalidated_rate_min",
		"label_quick_invalidated_rate_max",
		"label_quick_invalidated_rate_delta",
		"label_chopped_rate_min",
		"label_chopped_rate_max",
		"label_chopped_rate_delta",
		"label_trended_rate_min",
		"label_trended_rate_max",
		"label_trended_rate_delta",
		"label_avg_close_drift_pct_min",
		"label_avg_close_drift_pct_max",
		"label_avg_close_drift_pct_delta",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.ContextRule,
			strconv.Itoa(row.HoldBars),
			strconv.FormatBool(row.RequireMid50),
			strconv.Itoa(row.HorizonBars),
			strconv.Itoa(row.PeriodSplits),
			strconv.Itoa(row.SourceEpisodeCount),
			strconv.Itoa(row.SourceEpisodeCountMin),
			strconv.Itoa(row.SourceEpisodeCountMax),
			strconv.Itoa(row.SourceEpisodeCountDelta),
			strconv.Itoa(row.CandidateCount),
			strconv.Itoa(row.CandidateCountMin),
			strconv.Itoa(row.CandidateCountMax),
			strconv.Itoa(row.CandidateCountDelta),
			formatFloat(row.CandidateRateMin),
			formatFloat(row.CandidateRateMax),
			formatFloat(row.CandidateRateDelta),
			formatFloat(row.LabelPersistedInsideRangeRateMin),
			formatFloat(row.LabelPersistedInsideRangeRateMax),
			formatFloat(row.LabelPersistedInsideRangeRateDelta),
			formatFloat(row.LabelQuickInvalidatedRateMin),
			formatFloat(row.LabelQuickInvalidatedRateMax),
			formatFloat(row.LabelQuickInvalidatedRateDelta),
			formatFloat(row.LabelChoppedRateMin),
			formatFloat(row.LabelChoppedRateMax),
			formatFloat(row.LabelChoppedRateDelta),
			formatFloat(row.LabelTrendedRateMin),
			formatFloat(row.LabelTrendedRateMax),
			formatFloat(row.LabelTrendedRateDelta),
			formatFloat(row.LabelAvgCloseDriftPctMin),
			formatFloat(row.LabelAvgCloseDriftPctMax),
			formatFloat(row.LabelAvgCloseDriftPctDelta),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeHoldInsideDirectionalEdgeCandidatesCSV(path string, rows []lab.HoldInsideDirectionalEdgeCandidateRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"context_rule",
		"hold_bars",
		"require_mid_50",
		"split",
		"paper_side",
		"source_episode_id",
		"episode_start_index",
		"episode_end_index",
		"episode_start_time",
		"episode_end_time",
		"raw_length_bars",
		"active_length_bars",
		"raw_length_bucket",
		"active_length_bucket",
		"episode_high",
		"episode_low",
		"episode_mid",
		"episode_end_close",
		"episode_width_pct",
		"episode_width_bucket",
		"avg_normalized_atr",
		"end_normalized_atr",
		"width_to_atr_ratio",
		"width_to_atr_bucket",
		"decision_index",
		"decision_time",
		"decision_close",
		"decision_close_position",
		"decision_close_position_bucket",
		"decision_mid_side",
		"decision_distance_to_high_pct",
		"decision_distance_to_low_pct",
		"decision_distance_to_mid_pct",
		"horizon_bars",
		"label_window_start_index",
		"label_window_end_index",
		"label_window_start_time",
		"label_window_end_time",
		"label_favorable_move_pct",
		"label_adverse_move_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse",
		"label_touched_mid",
		"label_closed_across_mid",
		"label_side_boundary_touch",
		"label_opposite_close_break",
		"label_reentered_range",
		"label_persisted_inside_range",
		"label_quick_invalidated",
		"label_invalidated_up",
		"label_invalidated_down",
		"label_trended_up",
		"label_trended_down",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.ContextRule,
			strconv.Itoa(row.HoldBars),
			strconv.FormatBool(row.RequireMid50),
			row.Split,
			row.PaperSide,
			strconv.Itoa(row.SourceEpisodeID),
			strconv.Itoa(row.EpisodeStartIndex),
			strconv.Itoa(row.EpisodeEndIndex),
			row.EpisodeStartTime,
			row.EpisodeEndTime,
			strconv.Itoa(row.RawLengthBars),
			strconv.Itoa(row.ActiveLengthBars),
			row.RawLengthBucket,
			row.ActiveLengthBucket,
			formatFloat(row.EpisodeHigh),
			formatFloat(row.EpisodeLow),
			formatFloat(row.EpisodeMid),
			formatFloat(row.EpisodeEndClose),
			formatFloat(row.EpisodeWidthPct),
			row.EpisodeWidthBucket,
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.EndNormalizedATR),
			formatFloat(row.WidthToATRRatio),
			row.WidthToATRBucket,
			strconv.Itoa(row.DecisionIndex),
			row.DecisionTime,
			formatFloat(row.DecisionClose),
			formatFloat(row.DecisionClosePosition),
			row.DecisionClosePositionBucket,
			row.DecisionMidSide,
			formatFloat(row.DecisionDistanceToHighPct),
			formatFloat(row.DecisionDistanceToLowPct),
			formatFloat(row.DecisionDistanceToMidPct),
			strconv.Itoa(row.HorizonBars),
			strconv.Itoa(row.LabelWindowStartIndex),
			strconv.Itoa(row.LabelWindowEndIndex),
			row.LabelWindowStartTime,
			row.LabelWindowEndTime,
			formatFloat(row.LabelFavorableMovePct),
			formatFloat(row.LabelAdverseMovePct),
			formatFloat(row.LabelFavorableMinusAdverse),
			strconv.FormatBool(row.LabelFavorableGTAdverse),
			strconv.FormatBool(row.LabelTouchedMid),
			strconv.FormatBool(row.LabelClosedAcrossMid),
			strconv.FormatBool(row.LabelSideBoundaryTouch),
			strconv.FormatBool(row.LabelOppositeCloseBreak),
			strconv.FormatBool(row.LabelReenteredRange),
			strconv.FormatBool(row.LabelPersistedInsideRange),
			strconv.FormatBool(row.LabelQuickInvalidated),
			strconv.FormatBool(row.LabelInvalidatedUp),
			strconv.FormatBool(row.LabelInvalidatedDown),
			strconv.FormatBool(row.LabelTrendedUp),
			strconv.FormatBool(row.LabelTrendedDown),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeHoldInsideDirectionalEdgeSummaryCSV(path string, rows []lab.HoldInsideDirectionalEdgeSummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"context_rule",
		"hold_bars",
		"require_mid_50",
		"split",
		"horizon_bars",
		"paper_side",
		"decision_close_position_bucket",
		"source_episode_count",
		"candidate_count",
		"candidate_rate",
		"avg_raw_length_bars",
		"avg_active_length_bars",
		"avg_episode_width_pct",
		"avg_normalized_atr",
		"avg_end_normalized_atr",
		"avg_width_to_atr_ratio",
		"avg_decision_close_position",
		"avg_decision_distance_to_high_pct",
		"avg_decision_distance_to_low_pct",
		"avg_decision_distance_to_mid_pct",
		"label_favorable_greater_than_adverse_count",
		"label_touched_mid_count",
		"label_closed_across_mid_count",
		"label_side_boundary_touch_count",
		"label_opposite_close_break_count",
		"label_quick_invalidated_count",
		"label_invalidated_up_count",
		"label_invalidated_down_count",
		"label_trended_up_count",
		"label_trended_down_count",
		"label_avg_favorable_move_pct",
		"label_avg_adverse_move_pct",
		"label_avg_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
		"label_touched_mid_rate",
		"label_closed_across_mid_rate",
		"label_side_boundary_touch_rate",
		"label_opposite_close_break_rate",
		"label_quick_invalidated_rate",
		"label_invalidated_up_rate",
		"label_invalidated_down_rate",
		"label_trended_up_rate",
		"label_trended_down_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.ContextRule,
			strconv.Itoa(row.HoldBars),
			strconv.FormatBool(row.RequireMid50),
			row.Split,
			strconv.Itoa(row.HorizonBars),
			row.PaperSide,
			row.DecisionClosePositionBucket,
			strconv.Itoa(row.SourceEpisodeCount),
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.CandidateRate),
			formatFloat(row.AvgRawLengthBars),
			formatFloat(row.AvgActiveLengthBars),
			formatFloat(row.AvgEpisodeWidthPct),
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.AvgEndNormalizedATR),
			formatFloat(row.AvgWidthToATRRatio),
			formatFloat(row.AvgDecisionClosePosition),
			formatFloat(row.AvgDecisionDistanceToHighPct),
			formatFloat(row.AvgDecisionDistanceToLowPct),
			formatFloat(row.AvgDecisionDistanceToMidPct),
			strconv.Itoa(row.LabelFavorableGTAdverseCount),
			strconv.Itoa(row.LabelTouchedMidCount),
			strconv.Itoa(row.LabelClosedAcrossMidCount),
			strconv.Itoa(row.LabelSideBoundaryTouchCount),
			strconv.Itoa(row.LabelOppositeCloseBreakCount),
			strconv.Itoa(row.LabelQuickInvalidatedCount),
			strconv.Itoa(row.LabelInvalidatedUpCount),
			strconv.Itoa(row.LabelInvalidatedDownCount),
			strconv.Itoa(row.LabelTrendedUpCount),
			strconv.Itoa(row.LabelTrendedDownCount),
			formatFloat(row.LabelAvgFavorableMovePct),
			formatFloat(row.LabelAvgAdverseMovePct),
			formatFloat(row.LabelAvgFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGTAdverseRate),
			formatFloat(row.LabelTouchedMidRate),
			formatFloat(row.LabelClosedAcrossMidRate),
			formatFloat(row.LabelSideBoundaryTouchRate),
			formatFloat(row.LabelOppositeCloseBreakRate),
			formatFloat(row.LabelQuickInvalidatedRate),
			formatFloat(row.LabelInvalidatedUpRate),
			formatFloat(row.LabelInvalidatedDownRate),
			formatFloat(row.LabelTrendedUpRate),
			formatFloat(row.LabelTrendedDownRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeHoldInsideDirectionalEdgeStabilityCSV(path string, rows []lab.HoldInsideDirectionalEdgeStabilityRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"context_rule",
		"hold_bars",
		"require_mid_50",
		"horizon_bars",
		"paper_side",
		"decision_close_position_bucket",
		"period_splits",
		"source_episode_count",
		"source_episode_count_min",
		"source_episode_count_max",
		"source_episode_count_delta",
		"candidate_count",
		"candidate_count_min",
		"candidate_count_max",
		"candidate_count_delta",
		"candidate_rate_min",
		"candidate_rate_max",
		"candidate_rate_delta",
		"label_favorable_greater_than_adverse_rate_min",
		"label_favorable_greater_than_adverse_rate_max",
		"label_favorable_greater_than_adverse_rate_delta",
		"label_avg_favorable_minus_adverse_pct_min",
		"label_avg_favorable_minus_adverse_pct_max",
		"label_avg_favorable_minus_adverse_pct_delta",
		"label_avg_favorable_move_pct_min",
		"label_avg_favorable_move_pct_max",
		"label_avg_favorable_move_pct_delta",
		"label_avg_adverse_move_pct_min",
		"label_avg_adverse_move_pct_max",
		"label_avg_adverse_move_pct_delta",
		"label_touched_mid_rate_min",
		"label_touched_mid_rate_max",
		"label_touched_mid_rate_delta",
		"label_closed_across_mid_rate_min",
		"label_closed_across_mid_rate_max",
		"label_closed_across_mid_rate_delta",
		"label_side_boundary_touch_rate_min",
		"label_side_boundary_touch_rate_max",
		"label_side_boundary_touch_rate_delta",
		"label_opposite_close_break_rate_min",
		"label_opposite_close_break_rate_max",
		"label_opposite_close_break_rate_delta",
		"label_quick_invalidated_rate_min",
		"label_quick_invalidated_rate_max",
		"label_quick_invalidated_rate_delta",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.ContextRule,
			strconv.Itoa(row.HoldBars),
			strconv.FormatBool(row.RequireMid50),
			strconv.Itoa(row.HorizonBars),
			row.PaperSide,
			row.DecisionClosePositionBucket,
			strconv.Itoa(row.PeriodSplits),
			strconv.Itoa(row.SourceEpisodeCount),
			strconv.Itoa(row.SourceEpisodeCountMin),
			strconv.Itoa(row.SourceEpisodeCountMax),
			strconv.Itoa(row.SourceEpisodeCountDelta),
			strconv.Itoa(row.CandidateCount),
			strconv.Itoa(row.CandidateCountMin),
			strconv.Itoa(row.CandidateCountMax),
			strconv.Itoa(row.CandidateCountDelta),
			formatFloat(row.CandidateRateMin),
			formatFloat(row.CandidateRateMax),
			formatFloat(row.CandidateRateDelta),
			formatFloat(row.LabelFavorableGTAdverseRateMin),
			formatFloat(row.LabelFavorableGTAdverseRateMax),
			formatFloat(row.LabelFavorableGTAdverseRateDelta),
			formatFloat(row.LabelAvgFavorableMinusAdversePctMin),
			formatFloat(row.LabelAvgFavorableMinusAdversePctMax),
			formatFloat(row.LabelAvgFavorableMinusAdversePctDelta),
			formatFloat(row.LabelAvgFavorableMovePctMin),
			formatFloat(row.LabelAvgFavorableMovePctMax),
			formatFloat(row.LabelAvgFavorableMovePctDelta),
			formatFloat(row.LabelAvgAdverseMovePctMin),
			formatFloat(row.LabelAvgAdverseMovePctMax),
			formatFloat(row.LabelAvgAdverseMovePctDelta),
			formatFloat(row.LabelTouchedMidRateMin),
			formatFloat(row.LabelTouchedMidRateMax),
			formatFloat(row.LabelTouchedMidRateDelta),
			formatFloat(row.LabelClosedAcrossMidRateMin),
			formatFloat(row.LabelClosedAcrossMidRateMax),
			formatFloat(row.LabelClosedAcrossMidRateDelta),
			formatFloat(row.LabelSideBoundaryTouchRateMin),
			formatFloat(row.LabelSideBoundaryTouchRateMax),
			formatFloat(row.LabelSideBoundaryTouchRateDelta),
			formatFloat(row.LabelOppositeCloseBreakRateMin),
			formatFloat(row.LabelOppositeCloseBreakRateMax),
			formatFloat(row.LabelOppositeCloseBreakRateDelta),
			formatFloat(row.LabelQuickInvalidatedRateMin),
			formatFloat(row.LabelQuickInvalidatedRateMax),
			formatFloat(row.LabelQuickInvalidatedRateDelta),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeHoldInsideMidlineReactionCandidatesCSV(path string, rows []lab.HoldInsideMidlineReactionCandidateRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeHoldInsideMidlineReactionFunnelSummaryCSV(path string, rows []lab.HoldInsideMidlineReactionFunnelSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeHoldInsideMidlineReactionSummaryCSV(path string, rows []lab.HoldInsideMidlineReactionSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeHoldInsideMidlineReactionStabilityCSV(path string, rows []lab.HoldInsideMidlineReactionStabilityRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeHoldInsideMidlineTouchPrototypeSignalsCSV(path string, rows []lab.HoldInsideMidlineTouchPrototypeSignalRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeHoldInsideMidlineTouchPrototypeTradesCSV(path string, rows []lab.HoldInsideMidlineTouchPrototypeTradeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeHoldInsideMidlineTouchPrototypeSummaryCSV(path string, rows []lab.HoldInsideMidlineTouchPrototypeSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesImpulseAbsorptionCandidatesCSV(path string, rows []lab.FuturesImpulseAbsorptionCandidateRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesImpulseAbsorptionSummaryCSV(path string, rows []lab.FuturesImpulseAbsorptionSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesImpulseAbsorptionStabilityCSV(path string, rows []lab.FuturesImpulseAbsorptionStabilityRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeDiscoveryCoverageCSV(path string, rows []lab.FuturesRangeDiscoveryCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeDiscoveryCandidatesCSV(path string, rows []lab.FuturesRangeDiscoveryCandidateRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeDiscoverySummaryCSV(path string, rows []lab.FuturesRangeDiscoverySummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeDiscoveryRankingsCSV(path string, rows []lab.FuturesRangeDiscoveryRankingRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeDiscoveryStabilityCSV(path string, rows []lab.FuturesRangeDiscoveryStabilityRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseSourcesCSV(path string, rows []lab.FuturesRangeUniverseSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseCoverageCSV(path string, rows []lab.FuturesRangeUniverseCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseCandidatesCSV(path string, rows []lab.FuturesRangeUniverseCandidateRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseSummaryCSV(path string, rows []lab.FuturesRangeUniverseSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseRankingsCSV(path string, rows []lab.FuturesRangeUniverseRankingRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStabilityCSV(path string, rows []lab.FuturesRangeUniverseStabilityRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesCleanBreakoutBaselineSignalsCSV(path string, rows []lab.FuturesCleanBreakoutBaselineSignalRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesCleanBreakoutBaselineTradesCSV(path string, rows []lab.FuturesCleanBreakoutBaselineTradeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesCleanBreakoutBaselineSummaryCSV(path string, rows []lab.FuturesCleanBreakoutBaselineSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineSourcesCSV(path string, rows []lab.FuturesRangeUniverseSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineCoverageCSV(path string, rows []lab.FuturesRangeUniverseCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineSelectionCSV(path string, rows []lab.FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineSignalsCSV(path string, rows []lab.FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineTradesCSV(path string, rows []lab.FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseBreakoutRetestAcceptanceBaselineSummaryCSV(path string, rows []lab.FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionBaselineSourcesCSV(path string, rows []lab.FuturesRangeUniverseSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionBaselineCoverageCSV(path string, rows []lab.FuturesRangeUniverseCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionBaselineSignalsCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionSignalRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionBaselineTradesCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionTradeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionBaselineSummaryCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionOptimizationSourcesCSV(path string, rows []lab.FuturesRangeUniverseSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionOptimizationCoverageCSV(path string, rows []lab.FuturesRangeUniverseCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionOptimizationGridCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionOptimizationGridRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionOptimizationTradesCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionOptimizationTradeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionOptimizationSummaryCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionOptimizationRankingsCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionOptimizationRankingRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionStrategyReplaySourcesCSV(path string, rows []lab.FuturesRangeUniverseSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionStrategyReplayCoverageCSV(path string, rows []lab.FuturesRangeUniverseCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionStrategyReplaySignalsCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionStrategyReplayTradesCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionStrategyReplayTradeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionStrategyReplaySummaryCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionWalkForwardSourcesCSV(path string, rows []lab.FuturesRangeUniverseSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionWalkForwardCoverageCSV(path string, rows []lab.FuturesRangeUniverseCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionWalkForwardGridCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionOptimizationGridRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionWalkForwardFoldsCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionWalkForwardFoldRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionWalkForwardTradesCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionOptimizationTradeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionWalkForwardSummaryCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeUniverseStructuredCompressionWalkForwardRankingsCSV(path string, rows []lab.FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1SourcesCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1SourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1CoverageCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1CoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1GridCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1GridRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1BaselineCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1BaselineRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1SignalsCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1SignalRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1TradesCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1TradeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1SummaryCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1SummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1RankingsCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1RankingRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1SelectionCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1SelectionRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeFirstOccupancyRotationV1SkipsCSV(path string, rows []lab.FuturesRangeFirstOccupancyRotationV1SkipRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesHigherTFNestedRangeRotationSourcesCSV(path string, rows []lab.FuturesHigherTFNestedRangeRotationSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesHigherTFNestedRangeRotationCoverageCSV(path string, rows []lab.FuturesHigherTFNestedRangeRotationCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesHigherTFNestedRangeRotationParentRangesCSV(path string, rows []lab.FuturesHigherTFNestedRangeRotationParentRangeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesHigherTFNestedRangeRotationChildRangesCSV(path string, rows []lab.FuturesHigherTFNestedRangeRotationChildRangeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesHigherTFNestedRangeRotationEventsCSV(path string, rows []lab.FuturesHigherTFNestedRangeRotationEventRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesHigherTFNestedRangeRotationSummaryCSV(path string, rows []lab.FuturesHigherTFNestedRangeRotationSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextTriageSourcesCSV(path string, rows []lab.FuturesRangeContextTriageSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextTriageCoverageCSV(path string, rows []lab.FuturesRangeContextTriageCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextTriageEpisodesCSV(path string, rows []lab.FuturesRangeContextTriageEpisodeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextTriageQualityCSV(path string, rows []lab.FuturesRangeContextTriageQualityRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextTriageSessionsCSV(path string, rows []lab.FuturesRangeContextTriageSessionRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextTriageFailureModesCSV(path string, rows []lab.FuturesRangeContextTriageFailureModeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextTriageCohortsCSV(path string, rows []lab.FuturesRangeContextTriageCohortRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextTriageRankingsCSV(path string, rows []lab.FuturesRangeContextTriageRankingRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextTriageSummaryCSV(path string, rows []lab.FuturesRangeContextTriageSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeStateConstructionLoopSourcesCSV(path string, rows []lab.FuturesRangeStateConstructionLoopSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeStateConstructionLoopCoverageCSV(path string, rows []lab.FuturesRangeStateConstructionLoopCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeStateConstructionLoopFeatureWindowsCSV(path string, rows []lab.FuturesRangeStateConstructionLoopFeatureWindowRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeStateConstructionLoopStatesCSV(path string, rows []lab.FuturesRangeStateConstructionLoopStateRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeStateConstructionLoopLabelsCSV(path string, rows []lab.FuturesRangeStateConstructionLoopLabelRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeStateConstructionLoopCohortsCSV(path string, rows []lab.FuturesRangeStateConstructionLoopCohortRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeStateConstructionLoopRankingsCSV(path string, rows []lab.FuturesRangeStateConstructionLoopRankingRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeStateConstructionLoopSummaryCSV(path string, rows []lab.FuturesRangeStateConstructionLoopSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeStateConstructionLoopSkipsCSV(path string, rows []lab.FuturesRangeStateConstructionLoopSkipRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextRouterSourcesCSV(path string, rows []lab.FuturesRangeContextRouterSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextRouterCoverageCSV(path string, rows []lab.FuturesRangeContextRouterCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextRouterRulesCSV(path string, rows []lab.FuturesRangeContextRouterRuleRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextRouterRowsCSV(path string, rows []lab.FuturesRangeContextRouterRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextRouterCohortsCSV(path string, rows []lab.FuturesRangeContextRouterCohortRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextRouterRankingsCSV(path string, rows []lab.FuturesRangeContextRouterRankingRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextRouterSummaryCSV(path string, rows []lab.FuturesRangeContextRouterSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeContextRouterSkipsCSV(path string, rows []lab.FuturesRangeContextRouterSkipRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseSourcesCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseSourceRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseCoverageCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseCoverageRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseRouterDependencyCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseRouterDependencyRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseContextSegmentsCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseContextSegmentRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseEventsCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseEventRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseOutcomesCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseOutcomeRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseCohortsCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseCohortRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseRankingsCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseRankingRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseSummaryCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseSummaryRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func writeFuturesRangeRouterRotationPremiseSkipsCSV(path string, rows []lab.FuturesRangeRouterRotationPremiseSkipRow) error {
	return writeJSONTaggedCSV(path, rows)
}

func formatCleanBreakoutCoverage(rows []lab.FuturesRangeDiscoveryCoverageRow) string {
	parts := make([]string, 0, len(rows))
	for _, row := range rows {
		parts = append(parts, fmt.Sprintf("%s:%d", row.Timeframe, row.RowCount))
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, ",")
}

func writeHoldInsideMidlineTransitionCandidatesCSV(path string, rows []lab.HoldInsideMidlineTransitionCandidateRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"context_rule",
		"hold_bars",
		"require_mid_50",
		"split",
		"source_episode_id",
		"episode_start_index",
		"episode_end_index",
		"episode_start_time",
		"episode_end_time",
		"raw_length_bars",
		"active_length_bars",
		"raw_length_bucket",
		"active_length_bucket",
		"episode_high",
		"episode_low",
		"episode_mid",
		"episode_end_close",
		"episode_width_pct",
		"episode_width_bucket",
		"avg_normalized_atr",
		"end_normalized_atr",
		"width_to_atr_ratio",
		"width_to_atr_bucket",
		"decision_index",
		"decision_time",
		"decision_close",
		"decision_close_position",
		"decision_close_position_bucket",
		"decision_mid_side",
		"decision_distance_to_high_pct",
		"decision_distance_to_low_pct",
		"decision_distance_to_mid_pct",
		"horizon_bars",
		"label_window_start_index",
		"label_window_end_index",
		"label_window_start_time",
		"label_window_end_time",
		"label_touched_mid",
		"label_closed_across_mid",
		"label_first_mid_touch_delay_bars",
		"label_first_mid_close_across_delay_bars",
		"label_mid_touch_before_boundary_touch",
		"label_mid_cross_before_boundary_close_break",
		"label_reentered_range",
		"label_persisted_inside_range",
		"label_quick_invalidated",
		"label_invalidated_up",
		"label_invalidated_down",
		"label_trended_up",
		"label_trended_down",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.ContextRule,
			strconv.Itoa(row.HoldBars),
			strconv.FormatBool(row.RequireMid50),
			row.Split,
			strconv.Itoa(row.SourceEpisodeID),
			strconv.Itoa(row.EpisodeStartIndex),
			strconv.Itoa(row.EpisodeEndIndex),
			row.EpisodeStartTime,
			row.EpisodeEndTime,
			strconv.Itoa(row.RawLengthBars),
			strconv.Itoa(row.ActiveLengthBars),
			row.RawLengthBucket,
			row.ActiveLengthBucket,
			formatFloat(row.EpisodeHigh),
			formatFloat(row.EpisodeLow),
			formatFloat(row.EpisodeMid),
			formatFloat(row.EpisodeEndClose),
			formatFloat(row.EpisodeWidthPct),
			row.EpisodeWidthBucket,
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.EndNormalizedATR),
			formatFloat(row.WidthToATRRatio),
			row.WidthToATRBucket,
			strconv.Itoa(row.DecisionIndex),
			row.DecisionTime,
			formatFloat(row.DecisionClose),
			formatFloat(row.DecisionClosePosition),
			row.DecisionClosePositionBucket,
			row.DecisionMidSide,
			formatFloat(row.DecisionDistanceToHighPct),
			formatFloat(row.DecisionDistanceToLowPct),
			formatFloat(row.DecisionDistanceToMidPct),
			strconv.Itoa(row.HorizonBars),
			strconv.Itoa(row.LabelWindowStartIndex),
			strconv.Itoa(row.LabelWindowEndIndex),
			row.LabelWindowStartTime,
			row.LabelWindowEndTime,
			strconv.FormatBool(row.LabelTouchedMid),
			strconv.FormatBool(row.LabelClosedAcrossMid),
			strconv.Itoa(row.LabelFirstMidTouchDelayBars),
			strconv.Itoa(row.LabelFirstMidCloseAcrossDelayBars),
			strconv.FormatBool(row.LabelMidTouchBeforeBoundaryTouch),
			strconv.FormatBool(row.LabelMidCrossBeforeBoundaryBreak),
			strconv.FormatBool(row.LabelReenteredRange),
			strconv.FormatBool(row.LabelPersistedInsideRange),
			strconv.FormatBool(row.LabelQuickInvalidated),
			strconv.FormatBool(row.LabelInvalidatedUp),
			strconv.FormatBool(row.LabelInvalidatedDown),
			strconv.FormatBool(row.LabelTrendedUp),
			strconv.FormatBool(row.LabelTrendedDown),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeHoldInsideMidlineTransitionSummaryCSV(path string, rows []lab.HoldInsideMidlineTransitionSummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"context_rule",
		"hold_bars",
		"require_mid_50",
		"split",
		"horizon_bars",
		"decision_mid_side",
		"decision_close_position_bucket",
		"source_episode_count",
		"candidate_count",
		"candidate_rate",
		"avg_raw_length_bars",
		"avg_active_length_bars",
		"avg_episode_width_pct",
		"avg_normalized_atr",
		"avg_end_normalized_atr",
		"avg_width_to_atr_ratio",
		"avg_decision_close_position",
		"avg_decision_distance_to_high_pct",
		"avg_decision_distance_to_low_pct",
		"avg_decision_distance_to_mid_pct",
		"label_touched_mid_count",
		"label_closed_across_mid_count",
		"label_mid_touch_before_boundary_touch_count",
		"label_mid_cross_before_boundary_close_break_count",
		"label_reentered_range_count",
		"label_persisted_inside_range_count",
		"label_quick_invalidated_count",
		"label_invalidated_up_count",
		"label_invalidated_down_count",
		"label_trended_up_count",
		"label_trended_down_count",
		"label_touched_mid_rate",
		"label_closed_across_mid_rate",
		"label_mid_touch_before_boundary_touch_rate",
		"label_mid_cross_before_boundary_close_break_rate",
		"label_reentered_range_rate",
		"label_persisted_inside_range_rate",
		"label_quick_invalidated_rate",
		"label_invalidated_up_rate",
		"label_invalidated_down_rate",
		"label_trended_up_rate",
		"label_trended_down_rate",
		"label_avg_first_mid_touch_delay_bars",
		"label_avg_first_mid_close_across_delay_bars",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.ContextRule,
			strconv.Itoa(row.HoldBars),
			strconv.FormatBool(row.RequireMid50),
			row.Split,
			strconv.Itoa(row.HorizonBars),
			row.DecisionMidSide,
			row.DecisionClosePositionBucket,
			strconv.Itoa(row.SourceEpisodeCount),
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.CandidateRate),
			formatFloat(row.AvgRawLengthBars),
			formatFloat(row.AvgActiveLengthBars),
			formatFloat(row.AvgEpisodeWidthPct),
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.AvgEndNormalizedATR),
			formatFloat(row.AvgWidthToATRRatio),
			formatFloat(row.AvgDecisionClosePosition),
			formatFloat(row.AvgDecisionDistanceToHighPct),
			formatFloat(row.AvgDecisionDistanceToLowPct),
			formatFloat(row.AvgDecisionDistanceToMidPct),
			strconv.Itoa(row.LabelTouchedMidCount),
			strconv.Itoa(row.LabelClosedAcrossMidCount),
			strconv.Itoa(row.LabelMidTouchBeforeBoundaryTouchCount),
			strconv.Itoa(row.LabelMidCrossBeforeBoundaryBreakCount),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelPersistedInsideRangeCount),
			strconv.Itoa(row.LabelQuickInvalidatedCount),
			strconv.Itoa(row.LabelInvalidatedUpCount),
			strconv.Itoa(row.LabelInvalidatedDownCount),
			strconv.Itoa(row.LabelTrendedUpCount),
			strconv.Itoa(row.LabelTrendedDownCount),
			formatFloat(row.LabelTouchedMidRate),
			formatFloat(row.LabelClosedAcrossMidRate),
			formatFloat(row.LabelMidTouchBeforeBoundaryTouchRate),
			formatFloat(row.LabelMidCrossBeforeBoundaryBreakRate),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelPersistedInsideRangeRate),
			formatFloat(row.LabelQuickInvalidatedRate),
			formatFloat(row.LabelInvalidatedUpRate),
			formatFloat(row.LabelInvalidatedDownRate),
			formatFloat(row.LabelTrendedUpRate),
			formatFloat(row.LabelTrendedDownRate),
			formatFloat(row.LabelAvgFirstMidTouchDelayBars),
			formatFloat(row.LabelAvgFirstMidCloseAcrossDelayBars),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeHoldInsideMidlineTransitionStabilityCSV(path string, rows []lab.HoldInsideMidlineTransitionStabilityRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"context_rule",
		"hold_bars",
		"require_mid_50",
		"horizon_bars",
		"decision_mid_side",
		"decision_close_position_bucket",
		"period_splits",
		"source_episode_count",
		"source_episode_count_min",
		"source_episode_count_max",
		"source_episode_count_delta",
		"candidate_count",
		"candidate_count_min",
		"candidate_count_max",
		"candidate_count_delta",
		"candidate_rate_min",
		"candidate_rate_max",
		"candidate_rate_delta",
		"label_touched_mid_rate_min",
		"label_touched_mid_rate_max",
		"label_touched_mid_rate_delta",
		"label_closed_across_mid_rate_min",
		"label_closed_across_mid_rate_max",
		"label_closed_across_mid_rate_delta",
		"label_mid_touch_before_boundary_touch_rate_min",
		"label_mid_touch_before_boundary_touch_rate_max",
		"label_mid_touch_before_boundary_touch_rate_delta",
		"label_mid_cross_before_boundary_close_break_rate_min",
		"label_mid_cross_before_boundary_close_break_rate_max",
		"label_mid_cross_before_boundary_close_break_rate_delta",
		"label_reentered_range_rate_min",
		"label_reentered_range_rate_max",
		"label_reentered_range_rate_delta",
		"label_persisted_inside_range_rate_min",
		"label_persisted_inside_range_rate_max",
		"label_persisted_inside_range_rate_delta",
		"label_quick_invalidated_rate_min",
		"label_quick_invalidated_rate_max",
		"label_quick_invalidated_rate_delta",
		"label_invalidated_up_rate_min",
		"label_invalidated_up_rate_max",
		"label_invalidated_up_rate_delta",
		"label_invalidated_down_rate_min",
		"label_invalidated_down_rate_max",
		"label_invalidated_down_rate_delta",
		"label_trended_up_rate_min",
		"label_trended_up_rate_max",
		"label_trended_up_rate_delta",
		"label_trended_down_rate_min",
		"label_trended_down_rate_max",
		"label_trended_down_rate_delta",
		"label_trended_rate_min",
		"label_trended_rate_max",
		"label_trended_rate_delta",
		"label_avg_first_mid_touch_delay_bars_min",
		"label_avg_first_mid_touch_delay_bars_max",
		"label_avg_first_mid_touch_delay_bars_delta",
		"label_avg_first_mid_close_across_delay_bars_min",
		"label_avg_first_mid_close_across_delay_bars_max",
		"label_avg_first_mid_close_across_delay_bars_delta",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.ContextRule,
			strconv.Itoa(row.HoldBars),
			strconv.FormatBool(row.RequireMid50),
			strconv.Itoa(row.HorizonBars),
			row.DecisionMidSide,
			row.DecisionClosePositionBucket,
			strconv.Itoa(row.PeriodSplits),
			strconv.Itoa(row.SourceEpisodeCount),
			strconv.Itoa(row.SourceEpisodeCountMin),
			strconv.Itoa(row.SourceEpisodeCountMax),
			strconv.Itoa(row.SourceEpisodeCountDelta),
			strconv.Itoa(row.CandidateCount),
			strconv.Itoa(row.CandidateCountMin),
			strconv.Itoa(row.CandidateCountMax),
			strconv.Itoa(row.CandidateCountDelta),
			formatFloat(row.CandidateRateMin),
			formatFloat(row.CandidateRateMax),
			formatFloat(row.CandidateRateDelta),
			formatFloat(row.LabelTouchedMidRateMin),
			formatFloat(row.LabelTouchedMidRateMax),
			formatFloat(row.LabelTouchedMidRateDelta),
			formatFloat(row.LabelClosedAcrossMidRateMin),
			formatFloat(row.LabelClosedAcrossMidRateMax),
			formatFloat(row.LabelClosedAcrossMidRateDelta),
			formatFloat(row.LabelMidTouchBeforeBoundaryTouchRateMin),
			formatFloat(row.LabelMidTouchBeforeBoundaryTouchRateMax),
			formatFloat(row.LabelMidTouchBeforeBoundaryTouchRateDelta),
			formatFloat(row.LabelMidCrossBeforeBoundaryBreakRateMin),
			formatFloat(row.LabelMidCrossBeforeBoundaryBreakRateMax),
			formatFloat(row.LabelMidCrossBeforeBoundaryBreakRateDelta),
			formatFloat(row.LabelReenteredRangeRateMin),
			formatFloat(row.LabelReenteredRangeRateMax),
			formatFloat(row.LabelReenteredRangeRateDelta),
			formatFloat(row.LabelPersistedInsideRangeRateMin),
			formatFloat(row.LabelPersistedInsideRangeRateMax),
			formatFloat(row.LabelPersistedInsideRangeRateDelta),
			formatFloat(row.LabelQuickInvalidatedRateMin),
			formatFloat(row.LabelQuickInvalidatedRateMax),
			formatFloat(row.LabelQuickInvalidatedRateDelta),
			formatFloat(row.LabelInvalidatedUpRateMin),
			formatFloat(row.LabelInvalidatedUpRateMax),
			formatFloat(row.LabelInvalidatedUpRateDelta),
			formatFloat(row.LabelInvalidatedDownRateMin),
			formatFloat(row.LabelInvalidatedDownRateMax),
			formatFloat(row.LabelInvalidatedDownRateDelta),
			formatFloat(row.LabelTrendedUpRateMin),
			formatFloat(row.LabelTrendedUpRateMax),
			formatFloat(row.LabelTrendedUpRateDelta),
			formatFloat(row.LabelTrendedDownRateMin),
			formatFloat(row.LabelTrendedDownRateMax),
			formatFloat(row.LabelTrendedDownRateDelta),
			formatFloat(row.LabelTrendedRateMin),
			formatFloat(row.LabelTrendedRateMax),
			formatFloat(row.LabelTrendedRateDelta),
			formatFloat(row.LabelAvgFirstMidTouchDelayBarsMin),
			formatFloat(row.LabelAvgFirstMidTouchDelayBarsMax),
			formatFloat(row.LabelAvgFirstMidTouchDelayBarsDelta),
			formatFloat(row.LabelAvgFirstMidCloseAcrossDelayBarsMin),
			formatFloat(row.LabelAvgFirstMidCloseAcrossDelayBarsMax),
			formatFloat(row.LabelAvgFirstMidCloseAcrossDelayBarsDelta),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeCompressionBreakoutCandidatesCSV(path string, rows []lab.CompressionBreakoutCandidateRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"breakout_delay_bars",
		"horizon_bars",
		"episode_raw_length_bucket",
		"episode_active_length_bucket",
		"episode_range_width_bucket",
		"breakout_move_bucket",
		"decision_true_range_expansion_bucket",
		"detector_profile_id",
		"candidate_count",
		"avg_episode_raw_length_bars",
		"avg_episode_active_length_bars",
		"avg_episode_range_width_pct",
		"avg_breakout_move_pct",
		"avg_decision_true_range_atr",
		"label_reentered_range_count",
		"label_opposite_close_break_count",
		"label_favorable_greater_than_adverse_count",
		"label_reentered_range_rate",
		"label_opposite_close_break_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.BreakoutDelayBars),
			strconv.Itoa(row.HorizonBars),
			row.EpisodeRawLengthBucket,
			row.EpisodeActiveLengthBucket,
			row.EpisodeRangeWidthBucket,
			row.BreakoutMoveBucket,
			row.DecisionTrueRangeExpansionBucket,
			row.DetectorProfileID,
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgEpisodeRawLengthBars),
			formatFloat(row.AvgEpisodeActiveLengthBars),
			formatFloat(row.AvgEpisodeRangeWidthPct),
			formatFloat(row.AvgBreakoutMovePct),
			formatFloat(row.AvgDecisionTrueRangeATR),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelOppositeCloseBreakCount),
			strconv.Itoa(row.LabelFavorableGreaterThanAdverseCount),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelOppositeCloseBreakRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeCompressionBreakoutSummaryCSV(path string, rows []lab.CompressionBreakoutSummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"horizon_bars",
		"detector_profile_id",
		"candidate_count",
		"avg_breakout_delay_bars",
		"avg_episode_raw_length_bars",
		"avg_episode_active_length_bars",
		"avg_episode_range_width_pct",
		"avg_breakout_move_pct",
		"avg_decision_true_range_atr",
		"label_reentered_range_rate",
		"label_opposite_close_break_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgBreakoutDelayBars),
			formatFloat(row.AvgEpisodeRawLengthBars),
			formatFloat(row.AvgEpisodeActiveLengthBars),
			formatFloat(row.AvgEpisodeRangeWidthPct),
			formatFloat(row.AvgBreakoutMovePct),
			formatFloat(row.AvgDecisionTrueRangeATR),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelOppositeCloseBreakRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeRangeRegimeDurabilityEpisodesCSV(path string, rows []lab.RangeRegimeDurabilityEpisodeRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"episode_id",
		"start_index",
		"end_index",
		"start_time",
		"end_time",
		"horizon_bars",
		"detector_profile_id",
		"raw_length_bars",
		"active_length_bars",
		"raw_length_bucket",
		"active_length_bucket",
		"episode_high",
		"episode_low",
		"episode_end_close",
		"episode_width_pct",
		"episode_width_bucket",
		"avg_normalized_atr",
		"end_normalized_atr",
		"width_to_atr_ratio",
		"width_to_atr_bucket",
		"label_window_start_index",
		"label_window_end_index",
		"label_window_start_time",
		"label_window_end_time",
		"label_reentered_range",
		"label_persisted_inside_range",
		"label_quick_invalidated",
		"label_invalidated_up",
		"label_invalidated_down",
		"label_chopped",
		"label_trended_up",
		"label_trended_down",
		"label_close_drift_pct",
		"label_max_up_move_pct",
		"label_max_down_move_pct",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			strconv.Itoa(row.EpisodeID),
			strconv.Itoa(row.StartIndex),
			strconv.Itoa(row.EndIndex),
			row.StartTime,
			row.EndTime,
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.Itoa(row.RawLengthBars),
			strconv.Itoa(row.ActiveLengthBars),
			row.RawLengthBucket,
			row.ActiveLengthBucket,
			formatFloat(row.EpisodeHigh),
			formatFloat(row.EpisodeLow),
			formatFloat(row.EpisodeEndClose),
			formatFloat(row.EpisodeWidthPct),
			row.EpisodeWidthBucket,
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.EndNormalizedATR),
			formatFloat(row.WidthToATRRatio),
			row.WidthToATRBucket,
			strconv.Itoa(row.LabelWindowStartIndex),
			strconv.Itoa(row.LabelWindowEndIndex),
			row.LabelWindowStartTime,
			row.LabelWindowEndTime,
			strconv.FormatBool(row.LabelReenteredRange),
			strconv.FormatBool(row.LabelPersistedInsideRange),
			strconv.FormatBool(row.LabelQuickInvalidated),
			strconv.FormatBool(row.LabelInvalidatedUp),
			strconv.FormatBool(row.LabelInvalidatedDown),
			strconv.FormatBool(row.LabelChopped),
			strconv.FormatBool(row.LabelTrendedUp),
			strconv.FormatBool(row.LabelTrendedDown),
			formatFloat(row.LabelCloseDriftPct),
			formatFloat(row.LabelMaxUpMovePct),
			formatFloat(row.LabelMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeRangeRegimeDurabilitySummaryCSV(path string, rows []lab.RangeRegimeDurabilitySummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"horizon_bars",
		"raw_length_bucket",
		"active_length_bucket",
		"episode_width_bucket",
		"width_to_atr_bucket",
		"detector_profile_id",
		"episode_count",
		"avg_raw_length_bars",
		"avg_active_length_bars",
		"avg_episode_width_pct",
		"avg_normalized_atr",
		"avg_end_normalized_atr",
		"avg_width_to_atr_ratio",
		"label_reentered_range_count",
		"label_persisted_inside_range_count",
		"label_quick_invalidated_count",
		"label_invalidated_up_count",
		"label_invalidated_down_count",
		"label_chopped_count",
		"label_trended_up_count",
		"label_trended_down_count",
		"label_reentered_range_rate",
		"label_persisted_inside_range_rate",
		"label_quick_invalidated_rate",
		"label_invalidated_up_rate",
		"label_invalidated_down_rate",
		"label_chopped_rate",
		"label_trended_up_rate",
		"label_trended_down_rate",
		"label_avg_close_drift_pct",
		"label_avg_max_up_move_pct",
		"label_avg_max_down_move_pct",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			strconv.Itoa(row.HorizonBars),
			row.RawLengthBucket,
			row.ActiveLengthBucket,
			row.EpisodeWidthBucket,
			row.WidthToATRBucket,
			row.DetectorProfileID,
			strconv.Itoa(row.EpisodeCount),
			formatFloat(row.AvgRawLengthBars),
			formatFloat(row.AvgActiveLengthBars),
			formatFloat(row.AvgEpisodeWidthPct),
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.AvgEndNormalizedATR),
			formatFloat(row.AvgWidthToATRRatio),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelPersistedInsideRangeCount),
			strconv.Itoa(row.LabelQuickInvalidatedCount),
			strconv.Itoa(row.LabelInvalidatedUpCount),
			strconv.Itoa(row.LabelInvalidatedDownCount),
			strconv.Itoa(row.LabelChoppedCount),
			strconv.Itoa(row.LabelTrendedUpCount),
			strconv.Itoa(row.LabelTrendedDownCount),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelPersistedInsideRangeRate),
			formatFloat(row.LabelQuickInvalidatedRate),
			formatFloat(row.LabelInvalidatedUpRate),
			formatFloat(row.LabelInvalidatedDownRate),
			formatFloat(row.LabelChoppedRate),
			formatFloat(row.LabelTrendedUpRate),
			formatFloat(row.LabelTrendedDownRate),
			formatFloat(row.LabelAvgCloseDriftPct),
			formatFloat(row.LabelAvgMaxUpMovePct),
			formatFloat(row.LabelAvgMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRTouchAuditCSV(path string, rows []lab.SRAuditRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"index",
		"open_time",
		"close_time",
		"split",
		"close",
		"timeframe",
		"mode",
		"lookback_bars",
		"warmup_bars",
		"min_strength",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"qualified_zone_count",
		"raw_zone_count",
		"has_support",
		"near_support",
		"nearest_support",
		"nearest_support_distance",
		"nearest_support_distance_pct",
		"nearest_support_strength",
		"nearest_support_score",
		"nearest_support_top",
		"nearest_support_bottom",
		"nearest_support_last_touch_index",
		"nearest_support_source_pivots",
		"has_resistance",
		"near_resistance",
		"nearest_resistance",
		"nearest_resistance_distance",
		"nearest_resistance_distance_pct",
		"nearest_resistance_strength",
		"nearest_resistance_score",
		"nearest_resistance_top",
		"nearest_resistance_bottom",
		"nearest_resistance_last_touch_index",
		"nearest_resistance_source_pivots",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			strconv.Itoa(row.Index),
			row.OpenTime,
			row.CloseTime,
			row.Split,
			formatFloat(row.Close),
			row.Timeframe,
			row.Mode,
			strconv.Itoa(row.LookbackBars),
			strconv.Itoa(row.WarmupBars),
			strconv.Itoa(row.MinStrength),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.QualifiedZoneCount),
			strconv.Itoa(row.RawZoneCount),
			strconv.FormatBool(row.HasSupport),
			strconv.FormatBool(row.NearSupport),
			formatFloat(row.NearestSupport),
			formatFloat(row.NearestSupportDistance),
			formatFloat(row.NearestSupportDistancePct),
			strconv.Itoa(row.NearestSupportStrength),
			formatFloat(row.NearestSupportScore),
			formatFloat(row.NearestSupportTop),
			formatFloat(row.NearestSupportBottom),
			strconv.Itoa(row.NearestSupportLastTouchIndex),
			formatIntSlice(row.NearestSupportSourcePivots),
			strconv.FormatBool(row.HasResistance),
			strconv.FormatBool(row.NearResistance),
			formatFloat(row.NearestResistance),
			formatFloat(row.NearestResistanceDistance),
			formatFloat(row.NearestResistanceDistancePct),
			strconv.Itoa(row.NearestResistanceStrength),
			formatFloat(row.NearestResistanceScore),
			formatFloat(row.NearestResistanceTop),
			formatFloat(row.NearestResistanceBottom),
			strconv.Itoa(row.NearestResistanceLastTouchIndex),
			formatIntSlice(row.NearestResistanceSourcePivots),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRBoundaryEventsCSV(path string, rows []lab.SRBoundaryEventRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"index",
		"open_time",
		"close_time",
		"split",
		"side",
		"close",
		"boundary_price",
		"zone_top",
		"zone_bottom",
		"zone_width",
		"rejection_threshold",
		"distance_pct",
		"strength",
		"strength_bucket",
		"score",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"horizon_bars",
		"future_max_high",
		"future_min_low",
		"future_close",
		"favorable_move",
		"adverse_move",
		"favorable_move_pct",
		"adverse_move_pct",
		"distance_bucket",
		"wick_break",
		"close_break",
		"reclaimed_after_break",
		"rejected",
		"favorable_greater_than_adverse",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			strconv.Itoa(row.Index),
			row.OpenTime,
			row.CloseTime,
			row.Split,
			row.Side,
			formatFloat(row.Close),
			formatFloat(row.BoundaryPrice),
			formatFloat(row.ZoneTop),
			formatFloat(row.ZoneBottom),
			formatFloat(row.ZoneWidth),
			formatFloat(row.RejectionThreshold),
			formatFloat(row.DistancePct),
			strconv.Itoa(row.Strength),
			row.StrengthBucket,
			formatFloat(row.Score),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.HorizonBars),
			formatFloat(row.FutureMaxHigh),
			formatFloat(row.FutureMinLow),
			formatFloat(row.FutureClose),
			formatFloat(row.FavorableMove),
			formatFloat(row.AdverseMove),
			formatFloat(row.FavorableMovePct),
			formatFloat(row.AdverseMovePct),
			row.DistanceBucket,
			strconv.FormatBool(row.WickBreak),
			strconv.FormatBool(row.CloseBreak),
			strconv.FormatBool(row.ReclaimedAfterBreak),
			strconv.FormatBool(row.Rejected),
			strconv.FormatBool(row.FavorableGreaterThanAdverse),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRBoundaryQualityCSV(path string, rows []lab.SRBoundaryQualityRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"horizon_bars",
		"strength_bucket",
		"distance_bucket",
		"event_count",
		"avg_score",
		"avg_distance_pct",
		"avg_favorable_pct",
		"median_favorable_pct",
		"avg_adverse_pct",
		"median_adverse_pct",
		"close_break_rate",
		"wick_break_rate",
		"reclaim_after_break_rate",
		"rejection_rate",
		"favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.StrengthBucket,
			row.DistanceBucket,
			strconv.Itoa(row.EventCount),
			formatFloat(row.AvgScore),
			formatFloat(row.AvgDistancePct),
			formatFloat(row.AvgFavorablePct),
			formatFloat(row.MedianFavorablePct),
			formatFloat(row.AvgAdversePct),
			formatFloat(row.MedianAdversePct),
			formatFloat(row.CloseBreakRate),
			formatFloat(row.WickBreakRate),
			formatFloat(row.ReclaimAfterBreakRate),
			formatFloat(row.RejectionRate),
			formatFloat(row.FavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRBoundaryCandidateComparisonCSV(path string, rows []lab.SRBoundaryCandidateComparisonRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"horizon_bars",
		"strength_bucket",
		"distance_bucket",
		"event_count",
		"close_break_count",
		"rejected_count",
		"reclaimed_after_break_count",
		"close_break_rate",
		"rejection_rate",
		"reclaim_event_rate",
		"reclaim_given_close_break_rate",
		"all_avg_favorable_pct",
		"all_avg_adverse_pct",
		"all_favorable_minus_adverse_pct",
		"all_favorable_greater_than_adverse_rate",
		"rejected_avg_favorable_pct",
		"rejected_avg_adverse_pct",
		"rejected_favorable_minus_adverse_pct",
		"rejected_favorable_greater_than_adverse_rate",
		"reclaimed_avg_favorable_pct",
		"reclaimed_avg_adverse_pct",
		"reclaimed_favorable_minus_adverse_pct",
		"reclaimed_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.StrengthBucket,
			row.DistanceBucket,
			strconv.Itoa(row.EventCount),
			strconv.Itoa(row.CloseBreakCount),
			strconv.Itoa(row.RejectedCount),
			strconv.Itoa(row.ReclaimedAfterBreakCount),
			formatFloat(row.CloseBreakRate),
			formatFloat(row.RejectionRate),
			formatFloat(row.ReclaimEventRate),
			formatFloat(row.ReclaimGivenCloseBreakRate),
			formatFloat(row.AllAvgFavorablePct),
			formatFloat(row.AllAvgAdversePct),
			formatFloat(row.AllFavorableMinusAdversePct),
			formatFloat(row.AllFavorableGreaterThanAdverseRate),
			formatFloat(row.RejectedAvgFavorablePct),
			formatFloat(row.RejectedAvgAdversePct),
			formatFloat(row.RejectedFavorableMinusAdversePct),
			formatFloat(row.RejectedFavorableGreaterThanAdverseRate),
			formatFloat(row.ReclaimedAvgFavorablePct),
			formatFloat(row.ReclaimedAvgAdversePct),
			formatFloat(row.ReclaimedFavorableMinusAdversePct),
			formatFloat(row.ReclaimedFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRRejectionTimingCandidatesCSV(path string, rows []lab.SRRejectionTimingCandidateRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"horizon_bars",
		"close_location",
		"touched_zone",
		"pierced_zone",
		"closed_back",
		"decision_rejection_candidate",
		"wick_beyond_bucket",
		"strength_bucket",
		"distance_bucket",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"avg_score",
		"avg_distance_pct",
		"avg_wick_beyond_pct",
		"label_close_break_count",
		"label_wick_break_count",
		"label_reclaimed_after_break_count",
		"label_rejected_count",
		"label_favorable_greater_than_adverse_count",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.CloseLocation,
			strconv.FormatBool(row.TouchedZone),
			strconv.FormatBool(row.PiercedZone),
			strconv.FormatBool(row.ClosedBack),
			strconv.FormatBool(row.DecisionRejectionCandidate),
			row.WickBeyondBucket,
			row.StrengthBucket,
			row.DistanceBucket,
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgScore),
			formatFloat(row.AvgDistancePct),
			formatFloat(row.AvgWickBeyondPct),
			strconv.Itoa(row.LabelCloseBreakCount),
			strconv.Itoa(row.LabelWickBreakCount),
			strconv.Itoa(row.LabelReclaimedAfterBreakCount),
			strconv.Itoa(row.LabelRejectedCount),
			strconv.Itoa(row.LabelFavorableGreaterThanAdverseCount),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRRejectionTimingSummaryCSV(path string, rows []lab.SRRejectionTimingSummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"horizon_bars",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"touched_count",
		"pierced_count",
		"closed_back_count",
		"decision_rejection_candidate_count",
		"touched_rate",
		"pierced_rate",
		"closed_back_rate",
		"decision_rejection_candidate_rate",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
		"decision_candidate_label_close_break_rate",
		"decision_candidate_label_wick_break_rate",
		"decision_candidate_label_reclaim_event_rate",
		"decision_candidate_label_reclaim_given_close_break_rate",
		"decision_candidate_label_rejection_rate",
		"decision_candidate_label_avg_favorable_pct",
		"decision_candidate_label_avg_adverse_pct",
		"decision_candidate_label_favorable_minus_adverse_pct",
		"decision_candidate_label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			strconv.Itoa(row.TouchedCount),
			strconv.Itoa(row.PiercedCount),
			strconv.Itoa(row.ClosedBackCount),
			strconv.Itoa(row.DecisionRejectionCandidateCount),
			formatFloat(row.TouchedRate),
			formatFloat(row.PiercedRate),
			formatFloat(row.ClosedBackRate),
			formatFloat(row.DecisionRejectionCandidateRate),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
			formatFloat(row.DecisionCandidateLabelCloseBreakRate),
			formatFloat(row.DecisionCandidateLabelWickBreakRate),
			formatFloat(row.DecisionCandidateLabelReclaimEventRate),
			formatFloat(row.DecisionCandidateLabelReclaimGivenCloseBreakRate),
			formatFloat(row.DecisionCandidateLabelRejectionRate),
			formatFloat(row.DecisionCandidateLabelAvgFavorablePct),
			formatFloat(row.DecisionCandidateLabelAvgAdversePct),
			formatFloat(row.DecisionCandidateLabelFavorableMinusAdversePct),
			formatFloat(row.DecisionCandidateLabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRConfirmationTimingCandidatesCSV(path string, rows []lab.SRConfirmationTimingCandidateRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"confirmation_delay_bars",
		"horizon_bars",
		"seed_close_location",
		"seed_pierced_zone",
		"seed_wick_beyond_bucket",
		"confirmation_close_location",
		"confirmation_favorable_close",
		"confirmation_wrong_side_close",
		"decision_confirmation_candidate",
		"strength_bucket",
		"distance_bucket",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"avg_score",
		"avg_distance_pct",
		"avg_seed_wick_beyond_pct",
		"avg_confirmation_move_pct",
		"label_close_break_count",
		"label_wick_break_count",
		"label_reclaimed_after_break_count",
		"label_rejected_count",
		"label_favorable_greater_than_adverse_count",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.ConfirmationDelayBars),
			strconv.Itoa(row.HorizonBars),
			row.SeedCloseLocation,
			strconv.FormatBool(row.SeedPiercedZone),
			row.SeedWickBeyondBucket,
			row.ConfirmationCloseLocation,
			strconv.FormatBool(row.ConfirmationFavorableClose),
			strconv.FormatBool(row.ConfirmationWrongSideClose),
			strconv.FormatBool(row.DecisionConfirmationCandidate),
			row.StrengthBucket,
			row.DistanceBucket,
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgScore),
			formatFloat(row.AvgDistancePct),
			formatFloat(row.AvgSeedWickBeyondPct),
			formatFloat(row.AvgConfirmationMovePct),
			strconv.Itoa(row.LabelCloseBreakCount),
			strconv.Itoa(row.LabelWickBreakCount),
			strconv.Itoa(row.LabelReclaimedAfterBreakCount),
			strconv.Itoa(row.LabelRejectedCount),
			strconv.Itoa(row.LabelFavorableGreaterThanAdverseCount),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRConfirmationTimingSummaryCSV(path string, rows []lab.SRConfirmationTimingSummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"confirmation_delay_bars",
		"horizon_bars",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"confirmation_favorable_close_count",
		"confirmation_wrong_side_close_count",
		"decision_confirmation_candidate_count",
		"confirmation_favorable_close_rate",
		"confirmation_wrong_side_close_rate",
		"decision_confirmation_candidate_rate",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
		"decision_candidate_label_close_break_rate",
		"decision_candidate_label_wick_break_rate",
		"decision_candidate_label_reclaim_event_rate",
		"decision_candidate_label_reclaim_given_close_break_rate",
		"decision_candidate_label_rejection_rate",
		"decision_candidate_label_avg_favorable_pct",
		"decision_candidate_label_avg_adverse_pct",
		"decision_candidate_label_favorable_minus_adverse_pct",
		"decision_candidate_label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.ConfirmationDelayBars),
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			strconv.Itoa(row.ConfirmationFavorableCloseCount),
			strconv.Itoa(row.ConfirmationWrongSideCloseCount),
			strconv.Itoa(row.DecisionConfirmationCandidateCount),
			formatFloat(row.ConfirmationFavorableCloseRate),
			formatFloat(row.ConfirmationWrongSideCloseRate),
			formatFloat(row.DecisionConfirmationCandidateRate),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
			formatFloat(row.DecisionCandidateLabelCloseBreakRate),
			formatFloat(row.DecisionCandidateLabelWickBreakRate),
			formatFloat(row.DecisionCandidateLabelReclaimEventRate),
			formatFloat(row.DecisionCandidateLabelReclaimGivenCloseBreakRate),
			formatFloat(row.DecisionCandidateLabelRejectionRate),
			formatFloat(row.DecisionCandidateLabelAvgFavorablePct),
			formatFloat(row.DecisionCandidateLabelAvgAdversePct),
			formatFloat(row.DecisionCandidateLabelFavorableMinusAdversePct),
			formatFloat(row.DecisionCandidateLabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRFalseBreakReclaimTimingCandidatesCSV(path string, rows []lab.SRFalseBreakReclaimTimingCandidateRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"break_delay_bars",
		"reclaim_delay_bars",
		"total_delay_bars",
		"horizon_bars",
		"anchor_close_location",
		"break_close_location",
		"reclaim_close_location",
		"break_move_bucket",
		"reclaim_move_bucket",
		"decision_false_break_reclaim_candidate",
		"strength_bucket",
		"distance_bucket",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"avg_score",
		"avg_distance_pct",
		"avg_break_move_pct",
		"avg_reclaim_move_pct",
		"label_close_break_count",
		"label_wick_break_count",
		"label_reclaimed_after_break_count",
		"label_rejected_count",
		"label_favorable_greater_than_adverse_count",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.BreakDelayBars),
			strconv.Itoa(row.ReclaimDelayBars),
			strconv.Itoa(row.TotalDelayBars),
			strconv.Itoa(row.HorizonBars),
			row.AnchorCloseLocation,
			row.BreakCloseLocation,
			row.ReclaimCloseLocation,
			row.BreakMoveBucket,
			row.ReclaimMoveBucket,
			strconv.FormatBool(row.DecisionFalseBreakReclaimCandidate),
			row.StrengthBucket,
			row.DistanceBucket,
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgScore),
			formatFloat(row.AvgDistancePct),
			formatFloat(row.AvgBreakMovePct),
			formatFloat(row.AvgReclaimMovePct),
			strconv.Itoa(row.LabelCloseBreakCount),
			strconv.Itoa(row.LabelWickBreakCount),
			strconv.Itoa(row.LabelReclaimedAfterBreakCount),
			strconv.Itoa(row.LabelRejectedCount),
			strconv.Itoa(row.LabelFavorableGreaterThanAdverseCount),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRFalseBreakReclaimTimingSummaryCSV(path string, rows []lab.SRFalseBreakReclaimTimingSummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"horizon_bars",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"decision_false_break_reclaim_candidate_count",
		"decision_false_break_reclaim_candidate_rate",
		"avg_break_delay_bars",
		"avg_reclaim_delay_bars",
		"avg_total_delay_bars",
		"avg_break_move_pct",
		"avg_reclaim_move_pct",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
		"label_decision_candidate_close_break_rate",
		"label_decision_candidate_wick_break_rate",
		"label_decision_candidate_reclaim_event_rate",
		"label_decision_candidate_reclaim_given_close_break_rate",
		"label_decision_candidate_rejection_rate",
		"label_decision_candidate_avg_favorable_pct",
		"label_decision_candidate_avg_adverse_pct",
		"label_decision_candidate_favorable_minus_adverse_pct",
		"label_decision_candidate_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			strconv.Itoa(row.DecisionFalseBreakReclaimCandidateCount),
			formatFloat(row.DecisionFalseBreakReclaimCandidateRate),
			formatFloat(row.AvgBreakDelayBars),
			formatFloat(row.AvgReclaimDelayBars),
			formatFloat(row.AvgTotalDelayBars),
			formatFloat(row.AvgBreakMovePct),
			formatFloat(row.AvgReclaimMovePct),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
			formatFloat(row.LabelDecisionCandidateCloseBreakRate),
			formatFloat(row.LabelDecisionCandidateWickBreakRate),
			formatFloat(row.LabelDecisionCandidateReclaimEventRate),
			formatFloat(row.LabelDecisionCandidateReclaimGivenCloseBreakRate),
			formatFloat(row.LabelDecisionCandidateRejectionRate),
			formatFloat(row.LabelDecisionCandidateAvgFavorablePct),
			formatFloat(row.LabelDecisionCandidateAvgAdversePct),
			formatFloat(row.LabelDecisionCandidateFavorableMinusAdversePct),
			formatFloat(row.LabelDecisionCandidateFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func formatIntSlice(values []int) string {
	if len(values) == 0 {
		return ""
	}
	out := strconv.Itoa(values[0])
	for _, value := range values[1:] {
		out += ";" + strconv.Itoa(value)
	}
	return out
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
