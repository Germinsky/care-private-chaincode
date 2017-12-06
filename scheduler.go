package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"encoding/json"
	"errors"
)

type Scheduler struct {
	logger *shim.ChaincodeLogger
}

func NewSchedulerDefault() Scheduler {
	var logger = shim.NewLogger("scheduler_default")
	return Scheduler{logger}
}

func (s Scheduler) ConstructScheduleKey(ownerId string) string {
	return "schedule:ownerId:" + ownerId
}

func (s Scheduler) GetByOwnerId(stub shim.ChaincodeStubInterface, ownerId string) (*Schedule, error) {
	scheduleId := s.ConstructScheduleKey(ownerId)
	scheduleBytes, err := stub.GetState(scheduleId)

	if err != nil {
		return nil, err
	}

	var schedule Schedule
	if scheduleBytes == nil {
		return nil, errors.New(fmt.Sprintf("Schedule with key '%v' not found", scheduleId))
	} else {
		json.Unmarshal(scheduleBytes, &schedule)
		s.logger.Infof("Retrieve schedule: %v", schedule)
	}

	return &schedule, nil
}

func (s Scheduler) Save(stub shim.ChaincodeStubInterface, schedule Schedule) (*Schedule, error) {
	scheduleKey := s.ConstructScheduleKey(schedule.OwnerId)

	scheduleBytes, err := stub.GetState(scheduleKey)

	if scheduleBytes != nil {
		errorMsg := fmt.Sprintf("Schedule with key '%v' already exists", scheduleKey)
		s.logger.Errorf(errorMsg)
		errors.New(errorMsg)
	}

	s.logger.Infof("Creating new schedule for owner %v", schedule.OwnerId)

	schedule.ScheduleId = scheduleKey;

	jsonSchedule, err := json.Marshal(schedule)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(scheduleKey, jsonSchedule)
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}