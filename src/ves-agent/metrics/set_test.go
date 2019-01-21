package metrics

import (
	"testing"
	"time"
	"ves-agent/evel"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

func measTestFactory(interval time.Duration) MeasurementFactory {
	return MeasurementFactoryFunc(func(vmID string, timestamp time.Time) (*evel.EventMeasurements, error) {
		return evel.NewMeasurements("name", "id", vmID, interval, timestamp.Add(-interval), timestamp), nil
	})
}

type MeasurementSetBuilderSuite struct {
	suite.Suite
}

func TestMeasurementSetBuilder(t *testing.T) {
	suite.Run(t, new(MeasurementSetBuilderSuite))
}

func (s *MeasurementSetBuilderSuite) TestFind() {
	interval := 10 * time.Second
	bld := NewEventMeasurementSetBuilder(measTestFactory(interval))

	s.Len(bld.Measurements(), 0)
	now := time.Now()
	meas, err := bld.find("id", now)
	s.NoError(err)
	s.NotNil(meas)
	s.Equal(now.UnixNano()/1000, meas.LastEpochMicrosec)
	s.Equal(now.Add(-interval).UnixNano()/1000, meas.StartEpochMicrosec)
	s.Equal(interval.Seconds(), meas.MeasurementInterval)
	s.Len(bld.Measurements(), 1)

	meas2, err := bld.find("id", now)
	s.NoError(err)
	s.NotNil(meas2)
	s.Len(bld.Measurements(), 1)
	s.Equal(meas, meas2)
}

func (s *MeasurementSetBuilderSuite) TestSetCPUPercent() {
	interval := 10 * time.Second
	bld := NewEventMeasurementSetBuilder(measTestFactory(interval))

	s.Len(bld.Measurements(), 0)
	now := time.Now()
	err := bld.Set("CPUUsageArray.PercentUsage", "vmid", now, 12, map[string]string{"CPUIdentifier": "cpu-0"})
	s.NoError(err)
	err = bld.Set("CPUUsageArray.PercentUsage", "vmid", now, 13, map[string]string{"CPUIdentifier": "cpu-1"})
	s.NoError(err)

	s.Len(bld.Measurements(), 1)
	s.Equal("vmid", bld.Measurements()[0].SourceName)
	s.Len(bld.Measurements()[0].CPUUsageArray, 2)
	s.Equal("cpu-0", bld.Measurements()[0].CPUUsageArray[0].CPUIdentifier)
	s.EqualValues(12, bld.Measurements()[0].CPUUsageArray[0].PercentUsage)
	s.Equal("cpu-1", bld.Measurements()[0].CPUUsageArray[1].CPUIdentifier)
	s.EqualValues(13, bld.Measurements()[0].CPUUsageArray[1].PercentUsage)
}

func (s *MeasurementSetBuilderSuite) TestSetFreeMem() {
	interval := 10 * time.Second
	bld := NewEventMeasurementSetBuilder(measTestFactory(interval))

	s.Len(bld.Measurements(), 0)
	now := time.Now()
	err := bld.Set("MemoryUsageArray.MemoryFree", "myVM", now, 12, map[string]string{"VMIdentifier": "myVM"})
	s.NoError(err)
	s.Len(bld.Measurements(), 1)
	s.Len(bld.Measurements()[0].MemoryUsageArray, 1)
	s.EqualValues(12, bld.Measurements()[0].MemoryUsageArray[0].MemoryFree)
	s.EqualValues("myVM", bld.Measurements()[0].MemoryUsageArray[0].VMIdentifier)

	err = bld.Set("MemoryUsageArray.MemoryFree", "myVM", now, 13, map[string]string{"VMIdentifier": "myVM"})
	s.NoError(err)
	s.Len(bld.Measurements(), 1)
	s.Len(bld.Measurements()[0].MemoryUsageArray, 1)
	s.EqualValues(13, bld.Measurements()[0].MemoryUsageArray[0].MemoryFree)
	s.EqualValues("myVM", bld.Measurements()[0].MemoryUsageArray[0].VMIdentifier)

	err = bld.Set("MemoryUsageArray.MemoryFree", "yourVM", now, 14, map[string]string{"VMIdentifier": "yourVM"})
	s.NoError(err)
	s.Len(bld.Measurements(), 2)
	s.Len(bld.Measurements()[1].MemoryUsageArray, 1)
}

func (s *MeasurementSetBuilderSuite) TestSetUsedMem() {
	interval := 10 * time.Second
	bld := NewEventMeasurementSetBuilder(measTestFactory(interval))

	s.Len(bld.Measurements(), 0)
	now := time.Now()
	err := bld.Set("MemoryUsageArray.MemoryUsed", "myVM", now, 12, map[string]string{"VMIdentifier": "myVM"})
	s.NoError(err)
	s.Len(bld.Measurements(), 1)
	s.Len(bld.Measurements()[0].MemoryUsageArray, 1)
	s.EqualValues(12, bld.Measurements()[0].MemoryUsageArray[0].MemoryUsed)
	s.EqualValues("myVM", bld.Measurements()[0].MemoryUsageArray[0].VMIdentifier)

	err = bld.Set("MemoryUsageArray.MemoryUsed", "myVM", now, 13, map[string]string{"VMIdentifier": "myVM"})
	s.NoError(err)
	s.Len(bld.Measurements(), 1)
	s.Len(bld.Measurements()[0].MemoryUsageArray, 1)
	s.EqualValues(13, bld.Measurements()[0].MemoryUsageArray[0].MemoryUsed)
	s.EqualValues("myVM", bld.Measurements()[0].MemoryUsageArray[0].VMIdentifier)

	err = bld.Set("MemoryUsageArray.MemoryUsed", "yourVM", now, 14, map[string]string{"VMIdentifier": "yourVM"})
	s.NoError(err)
	s.Len(bld.Measurements(), 2)
	s.Len(bld.Measurements()[1].MemoryUsageArray, 1)
}

func (s *MeasurementSetBuilderSuite) TestSetInvalidArgs() {
	interval := 10 * time.Second
	bld := NewEventMeasurementSetBuilder(measTestFactory(interval))

	s.Len(bld.Measurements(), 0)
	now := time.Now()
	err := bld.Set("MemoryUsageArray.MemoryUsed", "", now, 12, map[string]string{"VMIdentifier": "id"})
	s.Error(err)
	s.Len(bld.Measurements(), 0)
	err = bld.Set("", "id", now, 12, map[string]string{"VMIdentifier": "id"})
	s.Error(err)
	s.Len(bld.Measurements(), 0)
}

func (s *MeasurementSetBuilderSuite) TestSetUnexistingField() {
	interval := 10 * time.Second
	bld := NewEventMeasurementSetBuilder(measTestFactory(interval))

	s.Len(bld.Measurements(), 0)
	now := time.Now()
	err := bld.Set("UnexistingField", "id", now, 12, map[string]string{"VMIdentifier": "id"})
	s.Error(err)
	s.Panics(assert.PanicTestFunc(func() {
		bld.Set("MemoryUsageArray.MemoryUsed", "id", now, 12, map[string]string{"VMIdentifier": "id"})
	}))
}

func (s *MeasurementSetBuilderSuite) TestSetUnexistingSliceField() {
	interval := 10 * time.Second
	bld := NewEventMeasurementSetBuilder(measTestFactory(interval))

	s.Len(bld.Measurements(), 0)
	now := time.Now()
	err := bld.Set("MemoryUsageArray.UnexistingField", "id", now, 12, map[string]string{"VMIdentifier": "id"})
	s.Error(err)
}

func (s *MeasurementSetBuilderSuite) TestSetAdditionalObject() {
	interval := 10 * time.Second
	tt := time.Now()
	bld := NewEventMeasurementSetBuilder(measTestFactory(interval))
	err := bld.SetAdditionalObject("", "oname1", "oinst1", tt, 12, MeasKeys{"k1": "v1", "k2": "v2"})
	s.Error(err)
	bld.SetAdditionalObject("id1", "oname1", "oinst1", tt, 12, MeasKeys{"k1": "v1", "k2": "v2"})
	bld.SetAdditionalObject("id1", "oname1", "oinst2", tt, 13, MeasKeys{"k1": "v1", "k2": "v2"})
	bld.SetAdditionalObject("id1", "oname1", "oinst2", tt, 14, MeasKeys{"k1": "v1", "k2": "v2"})
	bld.SetAdditionalObject("id1", "oname1", "oinst2", tt, 15, MeasKeys{"k1": "v1", "k2": "v3"})
	bld.SetAdditionalObject("id1", "oname1", "oinst2", tt, 16, MeasKeys{"k3": "v1", "k2": "v2"})
	bld.SetAdditionalObject("id1", "oname2", "oinst1", tt, 17, MeasKeys{"k3": "v1", "k2": "v2"})
	bld.SetAdditionalObject("id2", "oname1", "oinst1", tt, 18, MeasKeys{"k1": "v1", "k2": "v2"})

	meas := bld.Measurements()
	s.NotNil(meas)
	s.Len(meas, 2)
	s.Equal("id1", meas[0].SourceName)
	s.Len(meas[0].AdditionalObjects, 2)
	s.Equal("oname1", meas[0].AdditionalObjects[0].ObjectName)
	s.Len(meas[0].AdditionalObjects[0].ObjectInstances, 3)
	s.Len(meas[0].AdditionalObjects[0].ObjectInstances[0].ObjectInstance, 2)
	v, ok := meas[0].AdditionalObjects[0].ObjectInstances[0].ObjectInstance["oinst1"]
	s.True(ok)
	s.Equal(float64(12), v)
	v, ok = meas[0].AdditionalObjects[0].ObjectInstances[0].ObjectInstance["oinst2"]
	s.True(ok)
	s.Equal(float64(14), v)
}

func TestMeasurementSetToBatch(t *testing.T) {
	meas := EventMeasurementSet{
		evel.NewMeasurements("meas1", "id1", "source1", 10*time.Second, time.Now(), time.Now()),
		evel.NewMeasurements("meas2", "id2", "source2", 10*time.Second, time.Now(), time.Now()),
	}
	batch := meas.Batch()
	assert.Equal(t, len(meas), len(batch))
	for i := range meas {
		assert.Equal(t, meas[i], batch[i])
	}
}

func TestMeasKeyMatchJSONObjectKeys(t *testing.T) {
	pstr := func(s string) *string {
		return &s
	}
	k := MeasKeys{"k1": "v1", "k2": "v2"}
	o := []evel.Key{
		{KeyName: "k1", KeyValue: pstr("v1")},
		{KeyName: "k2", KeyValue: pstr("v2")},
	}
	assert.True(t, k.MatchJSONObjectKeys(o))

	k = MeasKeys{"k3": "v1", "k2": "v2"}
	assert.False(t, k.MatchJSONObjectKeys(o))
	k = MeasKeys{"k1": "v1", "k2": "v3"}
	assert.False(t, k.MatchJSONObjectKeys(o))
	k = MeasKeys{"k1": "v1", "k2": "v2", "k3": "v3"}
	assert.False(t, k.MatchJSONObjectKeys(o))
}
