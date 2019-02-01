/*
	Copyright 2019 Nokia

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package convert

import (
	"github.com/nokia/onap-vespa/govel"
	"sync"

	log "github.com/sirupsen/logrus"
)

// FaultManagerState handles the alert internal state
type FaultManagerState interface {
	// NextFaultIndex return the Fault Index and increments it
	NextFaultIndex() (int32, error)
	// IncrementFaultSn increments the fault sequence number
	IncrementFaultSn(faultID int32) error
	// GetFaultSn return the fault sequence number
	GetFaultSn(faultID int32) int64
	// GetFaultStartEpoch return the startEpoch
	GetFaultStartEpoch(faultID int32) int64
	// SetFaultStartEpoch set the startEpoch
	SetFaultStartEpoch(faultID int32, epoch int64) error
	// InitAlertInfos update the alertInfos map
	InitAlertInfos(faultID int32) error
	// GetFaultInStorage checks if faultName already associated to an index
	GetFaultInStorage(faultName string) int32
	// StoreFaultInStorage stores the index associated to the faultName
	StoreFaultInStorage(faultName string, faultID int32) error
	// DeleteFaultInStorage delete Fault in storage
	DeleteFaultInStorage(faultName string) error
}

// AlertInfos struct used to store sequence and startepoch of the alert
type AlertInfos struct {
	Sequence   int64
	StartEpoch int64
}

type inMemState struct {
	index      int32
	storage    map[string]int32
	alertInfos map[int32]*AlertInfos
}

// NextFaultIndex return the Fault Index to use and increments it
func (mem *inMemState) NextFaultIndex() (int32, error) {
	i := mem.index
	mem.index++
	return i, nil
}

// GetFaultInStorage checks if faultName already exist
// return the associated index if exist else return 0
func (mem *inMemState) GetFaultInStorage(faultName string) int32 {
	if val, ok := mem.storage[faultName]; ok {
		return val
	}
	return 0
}

// StoreFaultInStorage stores the faultID associated to the faultName
func (mem *inMemState) StoreFaultInStorage(faultName string, faultID int32) error {
	mem.storage[faultName] = faultID
	return nil
}

// DeleteFaultInStorage delete the index associated to the faultName
func (mem *inMemState) DeleteFaultInStorage(faultName string) error {
	delete(mem.storage, faultName)
	return nil
}

// GetFaultSn return the sequence value of the faultID index
func (mem *inMemState) GetFaultSn(faultID int32) int64 {
	return mem.alertInfos[faultID].Sequence
}

// IncFaultSequenceNumber increment the sequence value of the faultID index
func (mem *inMemState) IncrementFaultSn(faultID int32) error {
	mem.alertInfos[faultID].Sequence++
	return nil
}

// GetFaultStartEpoch return the startEpoch value of the faultID index
func (mem *inMemState) GetFaultStartEpoch(faultID int32) int64 {
	return mem.alertInfos[faultID].StartEpoch
}

func (mem *inMemState) InitAlertInfos(faultID int32) error {
	mem.alertInfos[faultID] = &AlertInfos{1, 0}
	return nil
}

// SetFaultStartEpoch set the value epoch to the alert faultID
func (mem *inMemState) SetFaultStartEpoch(faultID int32, epoch int64) error {
	mem.alertInfos[faultID].StartEpoch = epoch
	return nil
}

// FaultManager struct used to manage and store fault
type FaultManager struct {
	state FaultManagerState
	lock  *sync.Mutex
	conf  *govel.EventConfiguration
}

// StatusResult describes the result of the operation on storage
type StatusResult int

// Possible values for StatusResult
const (
	InError      StatusResult = 0
	AlreadyExist StatusResult = 1
	Stored       StatusResult = 2
	Cleared      StatusResult = 3
	NotExist     StatusResult = 4
)

// NewFaultManagerWithState with state management
func NewFaultManagerWithState(conf *govel.EventConfiguration, state FaultManagerState) *FaultManager {
	return &FaultManager{
		//index:   0,
		//storage: make(map[string]int32),
		//alertInfos: make(map[int32]*AlertInfos),
		lock:  new(sync.Mutex),
		conf:  conf,
		state: state,
	}
}

// NewFaultManager create a FaultManager
// return a pointer on the FaultManager created
func NewFaultManager(conf *govel.EventConfiguration) *FaultManager {
	return NewFaultManagerWithState(conf, &inMemState{index: 1, storage: make(map[string]int32), alertInfos: make(map[int32]*AlertInfos)})
}

// GetFaultState returns a reference to the underlying fault manager state
func (fm *FaultManager) GetFaultState() FaultManagerState {
	return fm.state
}

// GetEventConf return a pointer on config.EventConfiguration
func (fm *FaultManager) GetEventConf() *govel.EventConfiguration {
	return fm.conf
}

// storeFault check if fault not already exist in the storage and store it;
// return status = stored and id = the new index in case of success;
// return status = alreadyExist if fault already exist;
func (fm *FaultManager) storeFault(faultName string) (StatusResult, int32) {
	var id int32
	var status StatusResult

	fm.lock.Lock()
	defer fm.lock.Unlock()
	//faultName := fault.EventName + "_" + fault.SourceName
	val := fm.state.GetFaultInStorage(faultName)
	if val != 0 {
		log.Warnf("fault name %s already exist with id %d", faultName, val)
		id = val
		status = AlreadyExist
		//fm.alertInfos[id].sequence++
		//return status, id
	} else {
		faultID, err := fm.state.NextFaultIndex()
		if err != nil {
			status = InError
			return status, 0
		}
		if err = fm.state.StoreFaultInStorage(faultName, faultID); err != nil {
			log.Error(err.Error())
			return InError, 0
		}
		log.Infof("store fault: %s with index %d \n", faultName, faultID)
		id = faultID
		status = Stored
		err = fm.state.InitAlertInfos(id)
		if err != nil {
			status = InError
			return status, 0
		}
	}
	return status, id
}

// clearFault remove the faultName in the storage;
// return status = inError if error;
// return status = cleared and id = the removed index in case of success;
func (fm *FaultManager) clearFault(faultName string) (StatusResult, int32) {
	var status StatusResult

	fm.lock.Lock()
	defer fm.lock.Unlock()
	id := fm.state.GetFaultInStorage(faultName)
	if id == 0 {
		log.Warning("clearFault: fault name not present in storage: " + faultName)
		status = NotExist
	} else {
		log.Infof("suppress fault name %s with id %d in storage \n ", faultName, id)
		//fm.state.DeleteFaultInStorage(faultName)
		status = Cleared
	}
	return status, id
}
