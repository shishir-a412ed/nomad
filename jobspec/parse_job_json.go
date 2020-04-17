package jobspec

import (
	"encoding/json"
	"fmt"
	"io"
	//	"reflect"

	"github.com/hashicorp/nomad/api"
)

func UnmarshalJSON(r io.Reader) error {
	var f interface{}
	if err := json.NewDecoder(r).Decode(&f); err != nil {
		return err
	}

	var job api.Job

	m := f.(map[string]interface{})
	f1 := m["Job"]
	m1 := f1.(map[string]interface{})
	for key, val := range m1 {
		if val == nil {
			continue
		}

		switch key {
		case "Stop":
			job.Stop = api.ConvertBoolToPtr(val.(bool))
		case "Region":
			job.Region = api.ConvertStringToPtr(val.(string))
		case "Namespace":
			job.Namespace = api.ConvertStringToPtr(val.(string))
		case "ID":
			job.ID = api.ConvertStringToPtr(val.(string))
		case "ParentID":
			job.ParentID = api.ConvertStringToPtr(val.(string))
		case "Name":
			job.Name = api.ConvertStringToPtr(val.(string))
		case "Type":
			job.Type = api.ConvertStringToPtr(val.(string))
		case "Dispatched":
			job.Dispatched = val.(bool)
		case "Meta":
			job.Meta = val.(map[string]string)
		case "ConsulToken":
			job.ConsulToken = api.ConvertStringToPtr(val.(string))
		case "VaultToken":
			job.VaultToken = api.ConvertStringToPtr(val.(string))
		case "Status":
			job.Status = api.ConvertStringToPtr(val.(string))
		case "StatusDescription":
			job.StatusDescription = api.ConvertStringToPtr(val.(string))
		case "Stable":
			job.Stable = api.ConvertBoolToPtr(val.(bool))
		case "AllAtOnce":
			job.AllAtOnce = api.ConvertBoolToPtr(val.(bool))
		case "Version":
			job.Version = api.ConvertUint64ToPtr(uint64(val.(float64)))
		case "SubmitTime":
			job.SubmitTime = api.ConvertInt64ToPtr(int64(val.(float64)))
		case "CreateIndex":
			job.CreateIndex = api.ConvertUint64ToPtr(uint64(val.(float64)))
		case "ModifyIndex":
			job.ModifyIndex = api.ConvertUint64ToPtr(uint64(val.(float64)))
		case "JobModifyIndex":
			job.JobModifyIndex = api.ConvertUint64ToPtr(uint64(val.(float64)))
		case "Priority":
			job.Priority = api.ConvertIntToPtr(val.(int))
		}
	}

	fmt.Printf("Job is %+v", job)
	return nil
}
