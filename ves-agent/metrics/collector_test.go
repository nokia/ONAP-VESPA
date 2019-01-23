package metrics

import (
	"github.com/nokia/onap-vespa/govel"
	"context"
	"errors"
	"testing"
	"time"
	"github.com/nokia/onap-vespa/ves-agent/config"

	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/suite"
)

type APIMock struct {
	mock.Mock
}

func (m *APIMock) Query(ctx context.Context, query string, ts time.Time) (model.Value, error) {
	args := m.MethodCalled("Query", ctx, query, ts)
	return args.Get(0).(model.Value), args.Error(1)
}

func (m *APIMock) QueryRange(ctx context.Context, query string, r v1.Range) (model.Value, error) {
	args := m.MethodCalled("QueryRange", ctx, query, r)
	return args.Get(0).(model.Value), args.Error(1)
}

func (m *APIMock) LabelValues(ctx context.Context, label string) (model.LabelValues, error) {
	args := m.MethodCalled("LabelValues", ctx, label)
	return args.Get(0).(model.LabelValues), args.Error(1)
}

func (m *APIMock) Series(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) ([]model.LabelSet, error) {
	args := m.MethodCalled("Series", ctx, matches, startTime, endTime)
	return args.Get(0).([]model.LabelSet), args.Error(1)
}

func (m *APIMock) Snapshot(ctx context.Context, skipHead bool) (v1.SnapshotResult, error) {
	args := m.MethodCalled("Snapshot", ctx, skipHead)
	return args.Get(0).(v1.SnapshotResult), args.Error(1)
}

func (m *APIMock) Targets(ctx context.Context) (v1.TargetsResult, error) {
	args := m.MethodCalled("Targets", ctx)
	return args.Get(0).(v1.TargetsResult), args.Error(1)
}

func (m *APIMock) AlertManagers(ctx context.Context) (v1.AlertManagersResult, error) {
	args := m.MethodCalled("AlertManagers", ctx)
	return args.Get(0).(v1.AlertManagersResult), args.Error(1)
}

func (m *APIMock) CleanTombstones(ctx context.Context) error {
	args := m.MethodCalled("CleanTombstones", ctx)
	return args.Error(0)
}

func (m *APIMock) Config(ctx context.Context) (v1.ConfigResult, error) {
	args := m.MethodCalled("Config", ctx)
	return args.Get(0).(v1.ConfigResult), args.Error(1)
}

func (m *APIMock) DeleteSeries(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) error {
	args := m.MethodCalled("DeleteSeries", ctx, matches, startTime, endTime)
	return args.Error(0)
}

func (m *APIMock) Flags(ctx context.Context) (v1.FlagsResult, error) {
	args := m.MethodCalled("Flags", ctx)
	return args.Get(0).(v1.FlagsResult), args.Error(1)
}

type CollectorTestSuite struct {
	suite.Suite
	confEvent   govel.EventConfiguration
	namingCodes map[string]string
}

func (s *CollectorTestSuite) SetupSuite() {
	s.confEvent = govel.EventConfiguration{
		VNFName:      "VNFName",
		NfNamingCode: "hsxp",
	}
	s.namingCodes = make(map[string]string)
	s.namingCodes["ope-1"] = "oam"
	s.namingCodes["ope-2"] = "pro"
}

func TestCollector(t *testing.T) {
	suite.Run(t, new(CollectorTestSuite))
}

func (s *CollectorTestSuite) TestNew() {
	col, err := NewCollector(&config.MeasurementConfiguration{
		Prometheus: config.PrometheusConfig{
			Address: "http://127.0.0.1:9090",
		}},
		&s.confEvent,
		s.namingCodes,
	)
	s.NoError(err)
	s.NotNil(col)

	col, err = NewCollector(&config.MeasurementConfiguration{
		Prometheus: config.PrometheusConfig{
			Address: "127.0.0.1:9090",
		}},
		&s.confEvent,
		s.namingCodes,
	)
	s.Error(err)
	s.Nil(col)
}

func (s *CollectorTestSuite) TestAdjustCollectionStartTime() {
	col := Collector{max: 20 * time.Minute}
	start := time.Now()
	end := start.Add(10*time.Minute + 5*time.Second)
	interval := 1 * time.Minute
	res := col.adjustCollectionStartTime(start, end, interval)
	s.Equal(start, res)

	end = start.Add(30*time.Minute + 5*time.Second)
	res = col.adjustCollectionStartTime(start, end, interval)
	s.Equal(end.Add(-col.max).Truncate(interval), res)
}

func (s *CollectorTestSuite) TestGetMatrix() {
	api := APIMock{}
	collector := Collector{
		state: &inMemState{},
		api:   &api,
	}

	matrix := model.Matrix{
		&model.SampleStream{
			Metric: model.Metric{"Label1": model.LabelValue("Value1")},
			Values: []model.SamplePair{
				{Timestamp: model.Time(0), Value: model.SampleValue(12)},
			},
		},
	}
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(matrix, nil)

	m, err := collector.getMatrix("foobar", v1.Range{})

	s.NoError(err)
	s.EqualValues(matrix, m)
	api.AssertExpectations(s.T())
}

func (s *CollectorTestSuite) TestGetMatrixError() {
	api := APIMock{}
	collector := Collector{
		state:  &inMemState{},
		api:    &api,
		evtCfg: &s.confEvent,
	}

	var r *model.Matrix
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(r, errors.New("foobar error"))

	m, err := collector.getMatrix("foobar", v1.Range{})

	s.Error(err)
	s.Nil(m)
	api.AssertExpectations(s.T())

	// Test error when returned value is not a matrix
	var r2 model.Vector
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(&r2, nil)

	m, err = collector.getMatrix("foobar", v1.Range{})

	s.Error(err)
	s.Nil(m)

	api.AssertExpectations(s.T())
}

func (s *CollectorTestSuite) TestCollectCpuMetrics() {
	api := APIMock{}
	collector := Collector{
		state: &inMemState{},
		api:   &api,
		rules: config.MetricRules{
			Metrics: []config.MetricRule{
				{
					Expr:      "foobar",
					Target:    "CPUUsageArray.PercentUsage",
					VMIDLabel: "{{.labels.VNFC}}",
					Labels: []config.Label{
						{Name: "CPUIdentifier", Expr: "{{.labels.VCID}}"},
					},
				},
			},
		},
		evtCfg:      &s.confEvent,
		namingCodes: s.namingCodes,
	}

	cpuMatrix := model.Matrix{
		&model.SampleStream{
			Metric: model.Metric{"VNFC": model.LabelValue("ope-1"), "VCID": "1"},
			Values: []model.SamplePair{
				{Timestamp: model.TimeFromUnix(10), Value: model.SampleValue(12)},
				{Timestamp: model.TimeFromUnix(11), Value: model.SampleValue(13)},
				{Timestamp: model.TimeFromUnix(12), Value: model.SampleValue(14)},
			},
		},
		&model.SampleStream{
			Metric: model.Metric{"VNFC": model.LabelValue("ope-1"), "VCID": "2"},
			Values: []model.SamplePair{
				{Timestamp: model.TimeFromUnix(10), Value: model.SampleValue(42)},
				{Timestamp: model.TimeFromUnix(11), Value: model.SampleValue(43)},
				{Timestamp: model.TimeFromUnix(12), Value: model.SampleValue(44)},
			},
		},
	}
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(cpuMatrix, nil)
	measSet, err := collector.CollectMetrics(time.Unix(0, 0), time.Now(), 1*time.Second)
	s.NoError(err)
	s.Len(measSet, 3)
	s.EqualValues(10000000, measSet[0].LastEpochMicrosec)
	s.EqualValues(9000000, measSet[0].StartEpochMicrosec)
	s.EqualValues(11000000, measSet[1].LastEpochMicrosec)
	s.EqualValues(10000000, measSet[1].StartEpochMicrosec)
	s.EqualValues("oam", measSet[0].NfcNamingCode)
	s.EqualValues("oam", measSet[1].NfcNamingCode)

	s.Equal("ope-1", measSet[0].SourceName)
	s.Len(measSet[0].CPUUsageArray, 2)
	s.EqualValues("1", measSet[0].CPUUsageArray[0].CPUIdentifier)
	s.EqualValues("2", measSet[0].CPUUsageArray[1].CPUIdentifier)

	s.EqualValues(12, measSet[0].CPUUsageArray[0].PercentUsage)
	s.EqualValues(42, measSet[0].CPUUsageArray[1].PercentUsage)
	api.AssertExpectations(s.T())

	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(cpuMatrix, nil)
	collector.state = &inMemState{} // Reset state
	meas2, err := collector.Run(time.Unix(0, 0), time.Now(), 1*time.Second)
	s.NoError(err)
	s.EqualValues(measSet, meas2)

	api.AssertExpectations(s.T())
}

func (s *CollectorTestSuite) TestCollectCpuMetricsFailed() {
	api := APIMock{}
	collector := Collector{
		state: &inMemState{},
		api:   &api,
		rules: config.MetricRules{
			Metrics: []config.MetricRule{
				{
					Expr:      "foobar",
					Target:    "CPUUsageArray.PercentUsage",
					VMIDLabel: "{{.labels.VNFC}}",
					Labels: []config.Label{
						{Name: "CPUIdentifier", Expr: "{{.labels.VCID}}"},
					},
				},
			},
		},
		evtCfg:      &s.confEvent,
		namingCodes: s.namingCodes,
	}

	cpuMatrix := model.Matrix{}
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(cpuMatrix, errors.New("foobar"))
	_, err := collector.CollectMetrics(time.Unix(0, 0), time.Now(), 1*time.Second)
	s.Error(err)

	api.AssertExpectations(s.T())
}

func (s *CollectorTestSuite) TestCollectMemMetrics() {
	api := APIMock{}
	collector := Collector{
		state: &inMemState{},
		api:   &api,
		rules: config.MetricRules{
			Metrics: []config.MetricRule{
				{
					Expr:      "foobar",
					Target:    "MemoryUsageArray.MemoryUsed",
					VMIDLabel: "{{.labels.VNFC}}",
					Labels: []config.Label{
						{Name: "VMIdentifier", Expr: "{{.labels.VNFC}}"},
					},
				}, {
					Expr:      "foobar",
					Target:    "MemoryUsageArray.MemoryFree",
					VMIDLabel: "{{.labels.VNFC}}",
					Labels: []config.Label{
						{Name: "VMIdentifier", Expr: "{{.labels.VNFC}}"},
					},
				},
			},
		},
		evtCfg:      &s.confEvent,
		namingCodes: s.namingCodes,
	}

	memMatrix := model.Matrix{
		&model.SampleStream{
			Metric: model.Metric{"VNFC": model.LabelValue("ope-1")},
			Values: []model.SamplePair{
				{Timestamp: model.TimeFromUnix(10), Value: model.SampleValue(12)},
				{Timestamp: model.TimeFromUnix(11), Value: model.SampleValue(13)},
				{Timestamp: model.TimeFromUnix(12), Value: model.SampleValue(14)},
			},
		},
		&model.SampleStream{
			Metric: model.Metric{"VNFC": model.LabelValue("ope-2")},
			Values: []model.SamplePair{
				{Timestamp: model.TimeFromUnix(10), Value: model.SampleValue(42)},
				{Timestamp: model.TimeFromUnix(11), Value: model.SampleValue(43)},
				{Timestamp: model.TimeFromUnix(12), Value: model.SampleValue(44)},
			},
		},
	}
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Twice().Return(memMatrix, nil)
	measSet, err := collector.CollectMetrics(time.Unix(0, 0), time.Now(), 1*time.Second)
	s.NoError(err)
	s.Len(measSet, 6)
	s.EqualValues(10000000, measSet[0].LastEpochMicrosec)
	s.EqualValues(9000000, measSet[0].StartEpochMicrosec)
	s.EqualValues(11000000, measSet[1].LastEpochMicrosec)
	s.EqualValues(10000000, measSet[1].StartEpochMicrosec)
	s.EqualValues("oam", measSet[0].NfcNamingCode)
	s.EqualValues("pro", measSet[3].NfcNamingCode)

	s.Len(measSet[0].MemoryUsageArray, 1)
	s.EqualValues("ope-1", measSet[0].MemoryUsageArray[0].VMIdentifier)
	s.EqualValues("ope-2", measSet[3].MemoryUsageArray[0].VMIdentifier)

	s.EqualValues(12, measSet[0].MemoryUsageArray[0].MemoryFree)
	s.EqualValues(42, measSet[3].MemoryUsageArray[0].MemoryFree)

	s.EqualValues(12, measSet[0].MemoryUsageArray[0].MemoryUsed)
	s.EqualValues(42, measSet[3].MemoryUsageArray[0].MemoryUsed)
	api.AssertExpectations(s.T())

	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Twice().Return(memMatrix, nil)
	collector.state = &inMemState{} // Reset state
	meas2, err := collector.Run(time.Unix(0, 0), time.Now(), 1*time.Second)
	s.NoError(err)
	s.EqualValues(measSet, meas2)

	api.AssertExpectations(s.T())
}

func (s *CollectorTestSuite) TestCollectMemMetricsFailed() {
	api := APIMock{}
	collector := Collector{
		state: &inMemState{},
		api:   &api,
		rules: config.MetricRules{
			Metrics: []config.MetricRule{
				{
					Expr:      "foobar",
					Target:    "MemoryUsageArray.MemoryUsed",
					VMIDLabel: "{{.labels.VNFC}}",
					Labels: []config.Label{
						{Name: "VMIdentifier", Expr: "{{.labels.VNFC}}"},
					},
				}, {
					Expr:      "foobar",
					Target:    "MemoryUsageArray.MemoryFree",
					VMIDLabel: "{{.labels.VNFC}}",
					Labels: []config.Label{
						{Name: "VMIdentifier", Expr: "{{.labels.VNFC}}"},
					},
				},
			},
		},
		evtCfg:      &s.confEvent,
		namingCodes: s.namingCodes,
	}

	memMatrix := model.Matrix{}
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(memMatrix, errors.New("foobar"))
	_, err := collector.CollectMetrics(time.Unix(0, 0), time.Now(), 1*time.Second)
	s.Error(err)

	api.AssertExpectations(s.T())
}

func (s *CollectorTestSuite) TestCollectMemMetricInvalidLabel() {
	api := APIMock{}
	collector := Collector{
		state: &inMemState{},
		api:   &api,
		rules: config.MetricRules{
			Metrics: []config.MetricRule{
				{
					Expr:      "foobar",
					Target:    "MemoryUsageArray.MemoryUsed",
					VMIDLabel: "{{.foo.bar}}",
				},
			},
		},
		evtCfg:      &s.confEvent,
		namingCodes: s.namingCodes,
	}

	memMatrix := model.Matrix{
		&model.SampleStream{
			Metric: model.Metric{"VNFC": model.LabelValue("ope-1")},
			Values: []model.SamplePair{
				{Timestamp: model.TimeFromUnix(10), Value: model.SampleValue(12)},
			},
		},
	}
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(memMatrix, nil)
	_, err := collector.CollectMetrics(time.Unix(0, 0), time.Now(), 1*time.Second)
	s.Error(err)

	api.AssertExpectations(s.T())
}

func (s *CollectorTestSuite) TestCollectMemMetricInvalidTargetField() {
	api := APIMock{}
	collector := Collector{
		state: &inMemState{},
		api:   &api,
		rules: config.MetricRules{
			Metrics: []config.MetricRule{
				{
					Expr:      "foobar",
					Target:    "MemoryUsageArray.UnexistingField",
					VMIDLabel: "{{.labels.VNFC}}",
					Labels: []config.Label{
						{Name: "VMIdentifier", Expr: "{{.labels.VNFC}}"},
					},
				},
			},
		},
		evtCfg:      &s.confEvent,
		namingCodes: s.namingCodes,
	}

	memMatrix := model.Matrix{
		&model.SampleStream{
			Metric: model.Metric{"VNFC": model.LabelValue("ope-1")},
			Values: []model.SamplePair{
				{Timestamp: model.TimeFromUnix(10), Value: model.SampleValue(12)},
			},
		},
	}
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(memMatrix, nil)
	_, err := collector.CollectMetrics(time.Unix(0, 0), time.Now(), 1*time.Second)
	s.Error(err)

	api.AssertExpectations(s.T())
}

func (s *CollectorTestSuite) TestCollectAdditionalJSONMetrics() {
	api := APIMock{}
	collector := Collector{
		state: &inMemState{},
		api:   &api,
		rules: config.MetricRules{
			Metrics: []config.MetricRule{
				{
					Expr:           "foobar",
					Target:         "AdditionalObjects",
					VMIDLabel:      "{{.labels.VNFC}}",
					ObjectName:     "oname",
					ObjectInstance: "oinst",
					ObjectKeys: []config.Label{
						{Name: "key", Expr: "{{.labels.KeyLabel}}"},
					},
				},
			},
		},
		evtCfg:      &s.confEvent,
		namingCodes: s.namingCodes,
	}

	matrix := model.Matrix{
		&model.SampleStream{
			Metric: model.Metric{"VNFC": model.LabelValue("ope-1"), "KeyLabel": "1"},
			Values: []model.SamplePair{
				{Timestamp: model.TimeFromUnix(10), Value: model.SampleValue(12)},
			},
		},
	}
	api.On("QueryRange", mock.Anything, "foobar", mock.Anything).Once().Return(matrix, nil)
	meas, err := collector.CollectMetrics(time.Unix(0, 0), time.Now(), 1*time.Second)
	s.NoError(err)
	s.NotNil(meas)
	s.Len(meas, 1)
	s.Len(meas[0].AdditionalObjects, 1)
	s.Equal("oname", meas[0].AdditionalObjects[0].ObjectName)
	s.Len(meas[0].AdditionalObjects[0].ObjectInstances, 1)
	s.Len(meas[0].AdditionalObjects[0].ObjectInstances[0].ObjectInstance, 1)
	v, ok := meas[0].AdditionalObjects[0].ObjectInstances[0].ObjectInstance["oinst"]
	s.True(ok)
	s.Equal(float64(12), v)
	api.AssertExpectations(s.T())
}
