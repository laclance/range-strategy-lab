package lab

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Zero-trade derivatives context source audit.
//
// This audit validates the durable, already-materialized Binance USDT-M futures
// mark/index/(optional) premium 5m source files, binds their provenance, and
// reports whether they can be anti-lookahead-aligned to the existing candle
// anchors well enough to justify a separate later context-audit brief. It does
// not compute context features, labels, cohorts, rankings, entries, exits, or
// any P&L. It answers a source and alignment question only.

const (
	FuturesDerivativesContextSourceAuditName = "futures_derivatives_context_source_audit"

	DerivativesContextSourceAuditStopStateRejectedSourceGap           = "derivatives_context_zero_trade_source_audit_rejected_source_gap"
	DerivativesContextSourceAuditStopStateRejectedLiveOrPrivateAPI    = "derivatives_context_zero_trade_source_audit_rejected_live_or_private_api_path"
	DerivativesContextSourceAuditStopStateRejectedTimestampOrFinality = "derivatives_context_zero_trade_source_audit_rejected_timestamp_or_finality_gap"
	DerivativesContextSourceAuditStopStateRejectedAlignmentGap        = "derivatives_context_zero_trade_source_audit_rejected_alignment_gap"
	DerivativesContextSourceAuditStopStateRejectedClosedFamilyRescue  = "derivatives_context_zero_trade_source_audit_rejected_closed_family_rescue"
	DerivativesContextSourceAuditStopStatePassedNeedsContextBrief     = "derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief"

	derivativesSourceOwner    = "Binance Data Vision"
	derivativesProduct        = "Binance USDT-M futures"
	derivativesInterval       = "5m"
	derivativesIntervalMs     = 300000
	derivativesTimezone       = "UTC"
	derivativesMissingPolicy  = "no imputation, interpolation, future fill, or symbol substitution; missing context surfaced as missing rows"
	derivativesFinalityRule   = "kline known only after close_time; future alignment uses source_close_time + publication_lag <= decision_candle_close_time"
	derivativesSchemaExpected = "open_time,open,high,low,close,close_time,source_object_id"
)

// FuturesDerivativesContextSourceFileConfig names one durable normalized
// derivatives source CSV produced by the materialization step.
type FuturesDerivativesContextSourceFileConfig struct {
	Symbol        string
	SourceFamily  string
	ArchiveFamily string
	Path          string
	Required      bool
	// AllowNonPositive permits non-positive OHLC (premium-index can be <= 0).
	AllowNonPositive bool
}

type FuturesDerivativesContextSourceAuditConfig struct {
	DerivativeSources          []FuturesDerivativesContextSourceFileConfig
	Anchors                    []FuturesRangeUniverseSourceConfig
	MinAlignedCoverage         float64
	ConservativeLagIntervals   int
	MaxExtraStalenessIntervals int
	// EraStartMs/EraEndMs bound the candle-anchor era (UTC ms). Zero values fall
	// back to the approved 2021-01-01..2026-06-16T23:55 era.
	EraStartMs int64
	EraEndMs   int64
	// RescueClosedFamily must stay false; a true value is rejected as a
	// closed-family rescue attempt.
	RescueClosedFamily bool
}

// DefaultFuturesDerivativesContextSourceAuditConfig wires the nine materialized
// derivatives source files and the three candle anchors.
func DefaultFuturesDerivativesContextSourceAuditConfig() FuturesDerivativesContextSourceAuditConfig {
	root := "../binance-bot/data/derivatives"
	derivFile := func(family, sym string) string {
		return filepath.Join(root, fmt.Sprintf("binance_usdm_%s_5m_%s_2021_2026.csv", family, sym))
	}
	sources := []FuturesDerivativesContextSourceFileConfig{}
	for _, sym := range []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"} {
		sources = append(sources,
			FuturesDerivativesContextSourceFileConfig{Symbol: sym, SourceFamily: "mark_price_klines", ArchiveFamily: "markPriceKlines", Path: derivFile("mark_price_klines", sym), Required: true},
			FuturesDerivativesContextSourceFileConfig{Symbol: sym, SourceFamily: "index_price_klines", ArchiveFamily: "indexPriceKlines", Path: derivFile("index_price_klines", sym), Required: true},
			FuturesDerivativesContextSourceFileConfig{Symbol: sym, SourceFamily: "premium_index_klines", ArchiveFamily: "premiumIndexKlines", Path: derivFile("premium_index_klines", sym), Required: false, AllowNonPositive: true},
		)
	}
	return FuturesDerivativesContextSourceAuditConfig{
		DerivativeSources: sources,
		Anchors: []FuturesRangeUniverseSourceConfig{
			{Symbol: "BTCUSDT", Path: "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv", ApprovedPath: "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv"},
			{Symbol: "ETHUSDT", Path: "../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv", ApprovedPath: "../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv"},
			{Symbol: "SOLUSDT", Path: "../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv", ApprovedPath: "../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv"},
		},
		MinAlignedCoverage:         0.99,
		ConservativeLagIntervals:   1,
		MaxExtraStalenessIntervals: 0,
		EraStartMs:                 time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
		EraEndMs:                   time.Date(2026, 6, 16, 23, 55, 0, 0, time.UTC).UnixMilli(),
	}
}

// ---- artifact row types ----

type FuturesDerivativesContextSourceRow struct {
	AuditName                 string `json:"audit_name"`
	Symbol                    string `json:"symbol"`
	SourceFamily              string `json:"source_family"`
	SourceOwner               string `json:"source_owner"`
	ArchiveFamily             string `json:"archive_family"`
	ProductScope              string `json:"product_scope"`
	NativeInterval            string `json:"native_interval"`
	DurablePath               string `json:"durable_path"`
	Required                  bool   `json:"required"`
	Schema                    string `json:"schema"`
	SchemaUnambiguous         bool   `json:"schema_unambiguous"`
	TimestampSemantics        string `json:"timestamp_semantics"`
	TimezoneAssumption        string `json:"timezone_assumption"`
	FinalityRule              string `json:"finality_rule"`
	MaxStalenessIntervals     int    `json:"max_staleness_intervals"`
	MissingDataPolicy         string `json:"missing_data_policy"`
	RowCount                  int    `json:"row_count"`
	ExpectedRowCount          int    `json:"expected_row_count"`
	FirstOpenTime             string `json:"first_open_time"`
	LastOpenTime              string `json:"last_open_time"`
	DuplicateCount            int    `json:"duplicate_count"`
	GapCount                  int    `json:"gap_count"`
	MissingIntervalCount      int    `json:"missing_interval_count"`
	PhysicalNonMonotonicCount int    `json:"physical_non_monotonic_count"`
	CloseTimeViolationCount   int    `json:"close_time_violation_count"`
	NonFiniteCount            int    `json:"non_finite_count"`
	NonPositiveCount          int    `json:"non_positive_count"`
	FileSHA256                string `json:"file_sha256"`
	FirstSourceObjectID       string `json:"first_source_object_id"`
	ComparisonOnly            bool   `json:"comparison_only"`
	ValidationStatus          string `json:"validation_status"`
	ValidationError           string `json:"validation_error,omitempty"`
}

type FuturesDerivativesContextCandleAnchorRow struct {
	AuditName                 string `json:"audit_name"`
	Symbol                    string `json:"symbol"`
	Role                      string `json:"role"`
	Path                      string `json:"path"`
	RowCount                  int    `json:"row_count"`
	FirstOpenTime             string `json:"first_open_time"`
	LastOpenTime              string `json:"last_open_time"`
	GapCount                  int    `json:"gap_count"`
	DuplicateCount            int    `json:"duplicate_count"`
	ZeroVolumeCount           int    `json:"zero_volume_count"`
	PhysicalNonMonotonicCount int    `json:"physical_non_monotonic_count"`
	ExpectedRowCount          int    `json:"expected_row_count"`
	ExpectedFactsPass         bool   `json:"expected_facts_pass"`
	ValidationStatus          string `json:"validation_status"`
}

type FuturesDerivativesContextCoverageRow struct {
	AuditName            string  `json:"audit_name"`
	Symbol               string  `json:"symbol"`
	SourceFamily         string  `json:"source_family"`
	Required             bool    `json:"required"`
	RowCount             int     `json:"row_count"`
	ExpectedRowCount     int     `json:"expected_row_count"`
	UniqueInEraRows      int     `json:"unique_in_era_rows"`
	OutOfEraRows         int     `json:"out_of_era_rows"`
	EraCoveragePct       float64 `json:"era_coverage_pct"`
	FrontMissingRows     int     `json:"front_missing_rows"`
	TailMissingRows      int     `json:"tail_missing_rows"`
	GapCount             int     `json:"gap_count"`
	MissingIntervalCount int     `json:"missing_interval_count"`
}

type FuturesDerivativesContextTimestampAlignmentRow struct {
	AuditName                string  `json:"audit_name"`
	Symbol                   string  `json:"symbol"`
	SourceFamily             string  `json:"source_family"`
	Required                 bool    `json:"required"`
	AnchorDecisionCandles    int     `json:"anchor_decision_candles"`
	ContemporaneousAligned   int     `json:"contemporaneous_aligned_rows"`
	ContemporaneousCoverage  float64 `json:"contemporaneous_coverage_pct"`
	ConservativeLagIntervals int     `json:"conservative_lag_intervals"`
	LagAlignedRows           int     `json:"lag_aligned_rows"`
	LagCoveragePct           float64 `json:"lag_coverage_pct"`
	LagMissingRows           int     `json:"lag_missing_rows"`
	LagWarmupBoundaryRows    int     `json:"lag_warmup_boundary_rows"`
	UsesFutureRows           bool    `json:"uses_future_rows"`
	ExactClosedIntervalJoin  bool    `json:"exact_closed_interval_join"`
	MeetsMinCoverage         bool    `json:"meets_min_coverage"`
}

type FuturesDerivativesContextPublicationLagRow struct {
	AuditName                string  `json:"audit_name"`
	Symbol                   string  `json:"symbol"`
	SourceFamily             string  `json:"source_family"`
	Required                 bool    `json:"required"`
	NativeInterval           string  `json:"native_interval"`
	PublicationLagProven     bool    `json:"publication_lag_proven"`
	ConservativeLagIntervals int     `json:"conservative_lag_intervals"`
	ConservativeLagMs        int64   `json:"conservative_lag_ms"`
	MaxStalenessIntervals    int     `json:"max_staleness_intervals"`
	ForwardFillEnabled       bool    `json:"forward_fill_enabled"`
	AntiLookaheadRule        string  `json:"anti_lookahead_rule"`
	FinalityRule             string  `json:"finality_rule"`
	LagCoveragePct           float64 `json:"lag_coverage_pct"`
}

type FuturesDerivativesContextMissingnessRow struct {
	AuditName             string  `json:"audit_name"`
	Symbol                string  `json:"symbol"`
	SourceFamily          string  `json:"source_family"`
	Required              bool    `json:"required"`
	GapCount              int     `json:"gap_count"`
	MissingIntervalCount  int     `json:"missing_interval_count"`
	MissingIntervalPct    float64 `json:"missing_interval_pct"`
	FrontMissingRows      int     `json:"front_missing_rows"`
	TailMissingRows       int     `json:"tail_missing_rows"`
	LargestGapIntervals   int     `json:"largest_gap_intervals"`
	LagMissingContextRows int     `json:"lag_missing_context_rows"`
	ForwardFilledRows     int     `json:"forward_filled_rows"`
	MissingDataPolicy     string  `json:"missing_data_policy"`
}

type FuturesDerivativesContextProvenanceRow struct {
	AuditName           string `json:"audit_name"`
	Symbol              string `json:"symbol"`
	SourceFamily        string `json:"source_family"`
	SourceOwner         string `json:"source_owner"`
	ArchiveFamily       string `json:"archive_family"`
	ProductScope        string `json:"product_scope"`
	NativeInterval      string `json:"native_interval"`
	DurablePath         string `json:"durable_path"`
	FileSHA256          string `json:"file_sha256"`
	FirstSourceObjectID string `json:"first_source_object_id"`
	ManifestObjectsPath string `json:"manifest_objects_path"`
	ManifestFilesPath   string `json:"manifest_files_path"`
	ComparisonOnly      bool   `json:"comparison_only"`
	ValidationStatus    string `json:"validation_status"`
}

type FuturesDerivativesContextSkipRow struct {
	AuditName     string `json:"audit_name"`
	Symbol        string `json:"symbol"`
	SourceFamily  string `json:"source_family"`
	Reason        string `json:"reason"`
	Count         int    `json:"count"`
	FirstOrPolicy string `json:"first_or_policy"`
}

type FuturesDerivativesContextSourceAuditSummaryRow struct {
	AuditName              string  `json:"audit_name"`
	Scope                  string  `json:"scope"`
	Symbol                 string  `json:"symbol"`
	SourceFamily           string  `json:"source_family"`
	Required               bool    `json:"required"`
	RowCount               int     `json:"row_count"`
	LagCoveragePct         float64 `json:"lag_coverage_pct"`
	MissingIntervalCount   int     `json:"missing_interval_count"`
	ValidationStatus       string  `json:"validation_status"`
	RequiredStreams        int     `json:"required_streams"`
	AlignedRequiredStreams int     `json:"aligned_required_streams"`
	MinRequiredCoverage    float64 `json:"min_required_lag_coverage_pct"`
	Trades                 int     `json:"trades"`
	ZeroTradeCompatible    bool    `json:"zero_trade_compatible"`
	StopState              string  `json:"stop_state"`
}

type FuturesDerivativesContextSourceAuditResult struct {
	SourceRows             []FuturesDerivativesContextSourceRow             `json:"source_rows"`
	CandleAnchorRows       []FuturesDerivativesContextCandleAnchorRow       `json:"candle_anchor_rows"`
	CoverageRows           []FuturesDerivativesContextCoverageRow           `json:"coverage_rows"`
	TimestampAlignmentRows []FuturesDerivativesContextTimestampAlignmentRow `json:"timestamp_alignment_rows"`
	PublicationLagRows     []FuturesDerivativesContextPublicationLagRow     `json:"publication_lag_rows"`
	MissingnessRows        []FuturesDerivativesContextMissingnessRow        `json:"missingness_rows"`
	ProvenanceRows         []FuturesDerivativesContextProvenanceRow         `json:"provenance_rows"`
	SkipRows               []FuturesDerivativesContextSkipRow               `json:"skip_rows"`
	SummaryRows            []FuturesDerivativesContextSourceAuditSummaryRow `json:"summary_rows"`
	AlignedRequiredStreams int                                              `json:"aligned_required_streams"`
	RequiredStreams        int                                              `json:"required_streams"`
	StopState              string                                           `json:"stop_state"`
}

// internal per-stream loaded data
type derivStreamData struct {
	cfg          FuturesDerivativesContextSourceFileConfig
	openClose    map[int64]float64 // in-era open_time -> close price
	rowCount     int
	firstOpen    int64
	lastOpen     int64
	duplicate    int
	gapCount     int
	missingVsEra int
	frontMissing int
	tailMissing  int
	largestGap   int
	physNonMono  int
	closeTimeBad int
	nonFinite    int
	nonPositive  int
	outOfEra     int
	sha256       string
	firstObjID   string
	schema       string
	schemaOK     bool
	loadErr      string
}

func RunFuturesDerivativesContextSourceAudit(cfg FuturesDerivativesContextSourceAuditConfig, splits []Split) (FuturesDerivativesContextSourceAuditResult, error) {
	_ = splits // source/alignment audit does not slice by period
	result := FuturesDerivativesContextSourceAuditResult{}
	if cfg.MinAlignedCoverage <= 0 {
		cfg.MinAlignedCoverage = 0.99
	}
	if cfg.ConservativeLagIntervals <= 0 {
		cfg.ConservativeLagIntervals = 1
	}
	eraStartMs := cfg.EraStartMs
	eraEndMs := cfg.EraEndMs
	if eraStartMs == 0 {
		eraStartMs = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	}
	if eraEndMs == 0 {
		eraEndMs = time.Date(2026, 6, 16, 23, 55, 0, 0, time.UTC).UnixMilli()
	}
	expectedRows := int((eraEndMs-eraStartMs)/derivativesIntervalMs + 1)
	lagMs := int64(cfg.ConservativeLagIntervals) * derivativesIntervalMs

	// Closed-family rescue guard.
	if cfg.RescueClosedFamily {
		result.StopState = DerivativesContextSourceAuditStopStateRejectedClosedFamilyRescue
		return result, fmt.Errorf("%s", result.StopState)
	}

	// Live/private path guard.
	for _, s := range cfg.DerivativeSources {
		if reason := derivativesNonLocalPathReason(s.Path); reason != "" {
			result.StopState = DerivativesContextSourceAuditStopStateRejectedLiveOrPrivateAPI
			return result, fmt.Errorf("%s: %s (%s)", result.StopState, s.Path, reason)
		}
	}

	// Load candle anchors (alignment infra).
	anchorSets := map[string]map[int64]struct{}{}
	for _, a := range cfg.Anchors {
		candles, row, err := LoadFuturesRangeUniverseSource(a, splits)
		anchorRow := FuturesDerivativesContextCandleAnchorRow{
			AuditName: FuturesDerivativesContextSourceAuditName, Symbol: row.Symbol, Role: "candle_alignment_anchor",
			Path: a.Path, ExpectedRowCount: expectedRows, ValidationStatus: row.ValidationStatus,
		}
		if err != nil {
			result.CandleAnchorRows = append(result.CandleAnchorRows, anchorRow)
			result.StopState = DerivativesContextSourceAuditStopStateRejectedSourceGap
			return result, fmt.Errorf("%s: candle anchor %s: %v", result.StopState, a.Symbol, err)
		}
		anchorRow.RowCount = row.RowCount
		anchorRow.FirstOpenTime = row.FirstOpenTime
		anchorRow.LastOpenTime = row.LastOpenTime
		anchorRow.GapCount = row.GapCount
		anchorRow.DuplicateCount = row.DuplicateCount
		anchorRow.ZeroVolumeCount = row.ZeroVolumeCount
		anchorRow.PhysicalNonMonotonicCount = row.PhysicalNonMonotonicCount
		anchorRow.ExpectedFactsPass = row.RowCount == expectedRows &&
			row.FirstOpenTime == derivativesFmtMs(eraStartMs) && row.LastOpenTime == derivativesFmtMs(eraEndMs) &&
			row.GapCount == 0 && row.DuplicateCount == 0
		result.CandleAnchorRows = append(result.CandleAnchorRows, anchorRow)

		set := make(map[int64]struct{}, len(candles))
		for _, c := range candles {
			set[c.OpenTime.UnixMilli()] = struct{}{}
		}
		anchorSets[row.Symbol] = set
	}

	// Load and validate each derivatives source stream.
	var rejectTimestamp, rejectSourceGap, rejectAlignment string

	for _, s := range cfg.DerivativeSources {
		data := loadDerivativesStream(s, eraStartMs, eraEndMs)

		manifestObjects := filepath.Join(filepath.Dir(s.Path), "manifests", "binance_usdm_derivatives_context_source_materialization_objects.csv")
		manifestFiles := filepath.Join(filepath.Dir(s.Path), "manifests", "binance_usdm_derivatives_context_source_materialization_files.csv")

		srcRow := FuturesDerivativesContextSourceRow{
			AuditName: FuturesDerivativesContextSourceAuditName, Symbol: s.Symbol, SourceFamily: s.SourceFamily,
			SourceOwner: derivativesSourceOwner, ArchiveFamily: s.ArchiveFamily, ProductScope: derivativesProduct,
			NativeInterval: derivativesInterval, DurablePath: s.Path, Required: s.Required,
			Schema: data.schema, SchemaUnambiguous: data.schemaOK,
			TimestampSemantics: "open_time/close_time UTC ms", TimezoneAssumption: derivativesTimezone,
			FinalityRule: derivativesFinalityRule, MaxStalenessIntervals: cfg.MaxExtraStalenessIntervals,
			MissingDataPolicy: derivativesMissingPolicy, ExpectedRowCount: expectedRows,
			RowCount: data.rowCount, DuplicateCount: data.duplicate, GapCount: data.gapCount,
			MissingIntervalCount: data.missingVsEra, PhysicalNonMonotonicCount: data.physNonMono,
			CloseTimeViolationCount: data.closeTimeBad, NonFiniteCount: data.nonFinite, NonPositiveCount: data.nonPositive,
			FileSHA256: data.sha256, FirstSourceObjectID: data.firstObjID, ComparisonOnly: !s.Required,
		}
		if data.rowCount > 0 {
			srcRow.FirstOpenTime = derivativesFmtMs(data.firstOpen)
			srcRow.LastOpenTime = derivativesFmtMs(data.lastOpen)
		}

		switch {
		case data.loadErr != "":
			srcRow.ValidationStatus = "rejected"
			srcRow.ValidationError = data.loadErr
			if s.Required && rejectSourceGap == "" {
				rejectSourceGap = fmt.Sprintf("%s/%s: %s", s.SourceFamily, s.Symbol, data.loadErr)
			}
		case !data.schemaOK:
			srcRow.ValidationStatus = "rejected"
			srcRow.ValidationError = "schema is not the expected " + derivativesSchemaExpected
			if s.Required && rejectTimestamp == "" {
				rejectTimestamp = fmt.Sprintf("%s/%s: ambiguous schema %q", s.SourceFamily, s.Symbol, data.schema)
			}
		case data.closeTimeBad > 0 || data.duplicate > 0 || data.physNonMono > 0 || data.nonFinite > 0 || (data.nonPositive > 0 && !s.AllowNonPositive):
			srcRow.ValidationStatus = "rejected"
			srcRow.ValidationError = fmt.Sprintf("integrity faults: close_time_violations=%d duplicates=%d non_monotonic=%d non_finite=%d non_positive=%d",
				data.closeTimeBad, data.duplicate, data.physNonMono, data.nonFinite, data.nonPositive)
			if s.Required && rejectTimestamp == "" {
				rejectTimestamp = fmt.Sprintf("%s/%s: %s", s.SourceFamily, s.Symbol, srcRow.ValidationError)
			}
		default:
			srcRow.ValidationStatus = "accepted"
			if data.gapCount > 0 {
				srcRow.ValidationStatus = "accepted_with_recorded_gaps"
			}
		}
		result.SourceRows = append(result.SourceRows, srcRow)

		// Coverage row.
		uniqueInEra := len(data.openClose)
		eraCov := 0.0
		if expectedRows > 0 {
			eraCov = float64(uniqueInEra) / float64(expectedRows)
		}
		result.CoverageRows = append(result.CoverageRows, FuturesDerivativesContextCoverageRow{
			AuditName: FuturesDerivativesContextSourceAuditName, Symbol: s.Symbol, SourceFamily: s.SourceFamily,
			Required: s.Required, RowCount: data.rowCount, ExpectedRowCount: expectedRows,
			UniqueInEraRows: uniqueInEra, OutOfEraRows: data.outOfEra, EraCoveragePct: eraCov,
			FrontMissingRows: data.frontMissing, TailMissingRows: data.tailMissing,
			GapCount: data.gapCount, MissingIntervalCount: data.missingVsEra,
		})

		// Provenance row.
		result.ProvenanceRows = append(result.ProvenanceRows, FuturesDerivativesContextProvenanceRow{
			AuditName: FuturesDerivativesContextSourceAuditName, Symbol: s.Symbol, SourceFamily: s.SourceFamily,
			SourceOwner: derivativesSourceOwner, ArchiveFamily: s.ArchiveFamily, ProductScope: derivativesProduct,
			NativeInterval: derivativesInterval, DurablePath: s.Path, FileSHA256: data.sha256,
			FirstSourceObjectID: data.firstObjID, ManifestObjectsPath: manifestObjects, ManifestFilesPath: manifestFiles,
			ComparisonOnly: !s.Required, ValidationStatus: srcRow.ValidationStatus,
		})

		// Alignment + missingness only when we have an anchor and usable rows.
		anchor, hasAnchor := anchorSets[s.Symbol]
		if !hasAnchor || data.loadErr != "" {
			continue
		}
		align := derivativesAlignToAnchor(data.openClose, anchor, eraStartMs, lagMs)
		meetsCov := align.lagCoverage >= cfg.MinAlignedCoverage
		result.TimestampAlignmentRows = append(result.TimestampAlignmentRows, FuturesDerivativesContextTimestampAlignmentRow{
			AuditName: FuturesDerivativesContextSourceAuditName, Symbol: s.Symbol, SourceFamily: s.SourceFamily,
			Required: s.Required, AnchorDecisionCandles: align.decisionCandles,
			ContemporaneousAligned: align.contemporaneous, ContemporaneousCoverage: align.contemporaneousCov,
			ConservativeLagIntervals: cfg.ConservativeLagIntervals, LagAlignedRows: align.lagAligned,
			LagCoveragePct: align.lagCoverage, LagMissingRows: align.lagMissing, LagWarmupBoundaryRows: align.warmup,
			UsesFutureRows: false, ExactClosedIntervalJoin: true, MeetsMinCoverage: meetsCov,
		})
		result.PublicationLagRows = append(result.PublicationLagRows, FuturesDerivativesContextPublicationLagRow{
			AuditName: FuturesDerivativesContextSourceAuditName, Symbol: s.Symbol, SourceFamily: s.SourceFamily,
			Required: s.Required, NativeInterval: derivativesInterval, PublicationLagProven: false,
			ConservativeLagIntervals: cfg.ConservativeLagIntervals, ConservativeLagMs: lagMs,
			MaxStalenessIntervals: cfg.MaxExtraStalenessIntervals, ForwardFillEnabled: cfg.MaxExtraStalenessIntervals > 0,
			AntiLookaheadRule: "source_close_time + " + strconv.Itoa(cfg.ConservativeLagIntervals) + "*5m <= decision_candle_close_time",
			FinalityRule:      derivativesFinalityRule, LagCoveragePct: align.lagCoverage,
		})
		missingPct := 0.0
		if expectedRows > 0 {
			missingPct = float64(data.missingVsEra) / float64(expectedRows)
		}
		result.MissingnessRows = append(result.MissingnessRows, FuturesDerivativesContextMissingnessRow{
			AuditName: FuturesDerivativesContextSourceAuditName, Symbol: s.Symbol, SourceFamily: s.SourceFamily,
			Required: s.Required, GapCount: data.gapCount, MissingIntervalCount: data.missingVsEra,
			MissingIntervalPct: missingPct, FrontMissingRows: data.frontMissing, TailMissingRows: data.tailMissing,
			LargestGapIntervals: data.largestGap, LagMissingContextRows: align.lagMissing, ForwardFilledRows: 0,
			MissingDataPolicy: derivativesMissingPolicy,
		})
		if align.lagMissing-align.warmup > 0 {
			result.SkipRows = append(result.SkipRows, FuturesDerivativesContextSkipRow{
				AuditName: FuturesDerivativesContextSourceAuditName, Symbol: s.Symbol, SourceFamily: s.SourceFamily,
				Reason: "missing_context_at_conservative_lag", Count: align.lagMissing - align.warmup,
				FirstOrPolicy: "surfaced_as_missing_no_fill",
			})
		}
		if align.warmup > 0 {
			result.SkipRows = append(result.SkipRows, FuturesDerivativesContextSkipRow{
				AuditName: FuturesDerivativesContextSourceAuditName, Symbol: s.Symbol, SourceFamily: s.SourceFamily,
				Reason: "lag_warmup_boundary", Count: align.warmup, FirstOrPolicy: "era_start_one_interval_warmup",
			})
		}
		if s.Required && srcRow.ValidationStatus != "rejected" && !meetsCov && rejectAlignment == "" {
			rejectAlignment = fmt.Sprintf("%s/%s: lag_coverage=%.5f below min=%.5f", s.SourceFamily, s.Symbol, align.lagCoverage, cfg.MinAlignedCoverage)
		}
	}

	// Basis derivation is intentionally deferred to a later context-audit stage.

	// Decide stop state (priority: source gap > timestamp/finality > alignment).
	stop := DerivativesContextSourceAuditStopStatePassedNeedsContextBrief
	switch {
	case rejectSourceGap != "":
		stop = DerivativesContextSourceAuditStopStateRejectedSourceGap
	case rejectTimestamp != "":
		stop = DerivativesContextSourceAuditStopStateRejectedTimestampOrFinality
	case rejectAlignment != "":
		stop = DerivativesContextSourceAuditStopStateRejectedAlignmentGap
	}

	// Required-stream coverage tally.
	requiredStreams, alignedRequired := 0, 0
	minRequiredCov := math.Inf(1)
	lagByKey := map[string]float64{}
	for _, ar := range result.TimestampAlignmentRows {
		lagByKey[ar.SourceFamily+"|"+ar.Symbol] = ar.LagCoveragePct
	}
	for _, sr := range result.SourceRows {
		if !sr.Required {
			continue
		}
		requiredStreams++
		cov := lagByKey[sr.SourceFamily+"|"+sr.Symbol]
		if sr.ValidationStatus != "rejected" && cov >= cfg.MinAlignedCoverage {
			alignedRequired++
		}
		if cov < minRequiredCov {
			minRequiredCov = cov
		}
	}
	if math.IsInf(minRequiredCov, 1) {
		minRequiredCov = 0
	}
	result.RequiredStreams = requiredStreams
	result.AlignedRequiredStreams = alignedRequired
	result.StopState = stop

	// Summary rows: one per stream + an overall row.
	for _, sr := range result.SourceRows {
		result.SummaryRows = append(result.SummaryRows, FuturesDerivativesContextSourceAuditSummaryRow{
			AuditName: FuturesDerivativesContextSourceAuditName, Scope: "stream", Symbol: sr.Symbol,
			SourceFamily: sr.SourceFamily, Required: sr.Required, RowCount: sr.RowCount,
			LagCoveragePct: lagByKey[sr.SourceFamily+"|"+sr.Symbol], MissingIntervalCount: sr.MissingIntervalCount,
			ValidationStatus: sr.ValidationStatus, Trades: 0, ZeroTradeCompatible: true, StopState: stop,
		})
	}
	result.SummaryRows = append(result.SummaryRows, FuturesDerivativesContextSourceAuditSummaryRow{
		AuditName: FuturesDerivativesContextSourceAuditName, Scope: "overall", RequiredStreams: requiredStreams,
		AlignedRequiredStreams: alignedRequired, MinRequiredCoverage: minRequiredCov, Trades: 0,
		ZeroTradeCompatible: true, StopState: stop,
	})

	// Evaluated stop states (including rejections) are reported via StopState and
	// the artifacts; they are not Go errors, so the CLI still writes every
	// artifact and prints the verdict. Only the pre-evaluation guards
	// (closed-family, live/private, anchor load) return an error.
	return result, nil
}

type derivAlignment struct {
	decisionCandles    int
	contemporaneous    int
	contemporaneousCov float64
	lagAligned         int
	lagMissing         int
	lagCoverage        float64
	warmup             int
}

// derivativesAlignToAnchor computes anti-lookahead alignment: decision candle D
// is served by the source interval D-lagMs (strictly prior, no future rows).
func derivativesAlignToAnchor(openClose map[int64]float64, anchor map[int64]struct{}, eraStartMs, lagMs int64) derivAlignment {
	a := derivAlignment{decisionCandles: len(anchor)}
	for d := range anchor {
		if _, ok := openClose[d]; ok {
			a.contemporaneous++
		}
		lagged := d - lagMs
		if lagged < eraStartMs {
			a.warmup++
			a.lagMissing++
			continue
		}
		if _, ok := openClose[lagged]; ok {
			a.lagAligned++
		} else {
			a.lagMissing++
		}
	}
	if a.decisionCandles > 0 {
		a.contemporaneousCov = float64(a.contemporaneous) / float64(a.decisionCandles)
		a.lagCoverage = float64(a.lagAligned) / float64(a.decisionCandles)
	}
	return a
}

func derivativesNonLocalPathReason(path string) string {
	lower := strings.ToLower(strings.TrimSpace(path))
	switch {
	case strings.HasPrefix(lower, "http://"), strings.HasPrefix(lower, "https://"):
		return "remote URL"
	case strings.HasPrefix(lower, "ws://"), strings.HasPrefix(lower, "wss://"):
		return "websocket stream"
	case strings.HasPrefix(lower, "/tmp/"), strings.Contains(lower, "/tmp/"):
		return "ephemeral /tmp path"
	case strings.Contains(lower, "fapi") || strings.Contains(lower, "api-key") || strings.Contains(lower, "signed"):
		return "live/private API path"
	}
	return ""
}

func derivativesFmtMs(ms int64) string {
	return time.UnixMilli(ms).UTC().Format(time.RFC3339)
}

// loadDerivativesStream reads one normalized derivatives CSV, validates schema
// and closed-interval semantics, and returns in-era rows plus stats.
func loadDerivativesStream(s FuturesDerivativesContextSourceFileConfig, eraStartMs, eraEndMs int64) *derivStreamData {
	data := &derivStreamData{cfg: s, openClose: map[int64]float64{}}
	raw, err := os.ReadFile(s.Path)
	if err != nil {
		data.loadErr = "read: " + err.Error()
		return data
	}
	sum := sha256.Sum256(raw)
	data.sha256 = hex.EncodeToString(sum[:])

	r := csv.NewReader(bytes.NewReader(raw))
	r.FieldsPerRecord = -1
	header, err := r.Read()
	if err != nil {
		data.loadErr = "header: " + err.Error()
		return data
	}
	data.schema = strings.Join(normalizeSchema(header), ",")
	data.schemaOK = data.schema == derivativesSchemaExpected

	var prevOpen int64
	havePrev := false
	var sortedOpens []int64
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			data.loadErr = "row: " + err.Error()
			return data
		}
		if len(rec) < 7 {
			data.loadErr = fmt.Sprintf("row has %d columns, want 7", len(rec))
			return data
		}
		ot, e1 := strconv.ParseInt(strings.TrimSpace(rec[0]), 10, 64)
		ct, e2 := strconv.ParseInt(strings.TrimSpace(rec[5]), 10, 64)
		if e1 != nil || e2 != nil {
			data.loadErr = "non-integer open_time/close_time"
			return data
		}
		o, eo := strconv.ParseFloat(strings.TrimSpace(rec[1]), 64)
		h, eh := strconv.ParseFloat(strings.TrimSpace(rec[2]), 64)
		l, el := strconv.ParseFloat(strings.TrimSpace(rec[3]), 64)
		c, ec := strconv.ParseFloat(strings.TrimSpace(rec[4]), 64)
		if eo != nil || eh != nil || el != nil || ec != nil {
			data.nonFinite++
			continue
		}
		for _, v := range []float64{o, h, l, c} {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				data.nonFinite++
			}
			if v <= 0 {
				data.nonPositive++
			}
		}
		if ct != ot+derivativesIntervalMs-1 {
			data.closeTimeBad++
		}
		if havePrev && ot <= prevOpen {
			data.physNonMono++
		}
		prevOpen = ot
		havePrev = true
		if ot < eraStartMs || ot > eraEndMs {
			data.outOfEra++
			continue
		}
		if _, ok := data.openClose[ot]; ok {
			data.duplicate++
			continue
		}
		data.openClose[ot] = c
		sortedOpens = append(sortedOpens, ot)
		if data.firstObjID == "" {
			data.firstObjID = strings.TrimSpace(rec[6])
		}
	}
	data.rowCount = len(sortedOpens)
	if data.rowCount == 0 {
		if data.loadErr == "" {
			data.loadErr = "no in-era rows"
		}
		return data
	}
	sort.Slice(sortedOpens, func(i, j int) bool { return sortedOpens[i] < sortedOpens[j] })
	data.firstOpen = sortedOpens[0]
	data.lastOpen = sortedOpens[len(sortedOpens)-1]
	data.frontMissing = int((data.firstOpen - eraStartMs) / derivativesIntervalMs)
	data.tailMissing = int((eraEndMs - data.lastOpen) / derivativesIntervalMs)
	for i := 1; i < len(sortedOpens); i++ {
		d := sortedOpens[i] - sortedOpens[i-1]
		if d != derivativesIntervalMs {
			data.gapCount++
			miss := int(d/derivativesIntervalMs) - 1
			data.missingVsEra += miss
			if miss > data.largestGap {
				data.largestGap = miss
			}
		}
	}
	data.missingVsEra += data.frontMissing + data.tailMissing
	return data
}
