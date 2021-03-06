package controllers

import (
	"errors"
	"strconv"

	"github.com/astaxie/beego/orm"

	"../core"
	"../helpers"
	"../models"
)

// ServiceRecordController ServiceRecordController
type ServiceRecordController struct {
	BaseController
}

// Get Get
func (c *ServiceRecordController) Get() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecords := []*models.ServiceRecord{}
	qs := core.GetModelQuerySeter(nil, new(models.ServiceRecord), true)
	qs = qs.Filter("congregation_id", user.Congregation.ID)
	qs, _, subLimit, _ := c.SetQuerySeterByURIParam(qs)
	qs.All(&serviceRecords)

	models.LoadRecordsDetails(serviceRecords)
	c.Success(subLimit, serviceRecords)
}

// Post Post
func (c *ServiceRecordController) Post() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = c.ParseJSONBodyStruct(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecord.Congregation = user.Congregation
	serviceRecord.Recorder = user

	if serviceRecord.Area == "" {
		c.Error(errors.New("ERROR_MISSING_AREA"))
		return
	}

	if serviceRecord.StartedAt == nil {
		c.Error(errors.New("ERROR_MISSING_STARTED_AT"))
		return
	}

	if serviceRecord.LeaderName == "" {
		c.Error(errors.New("ERROR_MISSING_LEADER_NAME"))
		return
	}

	if serviceRecord.NumOfRecordsInTheSameDayForTheSameArea() > 0 {
		c.Error(errors.New("RECORD_ALREADY_EXISTS"))
		return
	}

	_, err = core.InsertModel(nil, serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}

// Delete Delete
func (c *ServiceRecordController) Delete() {
	user, err := c.GetAuthUser()
	if err != nil || user.Role == "public" {
		c.Error(err)
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = core.GetModelQuerySeter(nil, new(models.ServiceRecord), true).Filter("id", id).One(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	err = core.DeleteModel(nil, serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}

// GetWithDayName GetWithDayName
func (c *ServiceRecordController) GetWithDayName() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	urlParam := "id"

	arg := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		urlParam = "day"
	}

	serviceRecord := new(models.ServiceRecord)
	serviceRecords := []*models.ServiceRecord{}
	if urlParam == "id" {
		err = core.GetModelQuerySeter(nil, new(models.ServiceRecord), true).Filter("id", id).One(serviceRecord)
		c.Success(1, serviceRecord)
		return
	} else if urlParam == "day" {
		var total int64
		days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		for _, d := range days {
			if d == arg {
				o := orm.NewOrm()
				_, err = o.Raw("SELECT id FROM service_records WHERE dayname(started_at) = ? AND congregation_id = ?", arg, user.Congregation.ID).QueryRows(&serviceRecords)
				break
			}
		}
		ids := []int64{}
		for _, r := range serviceRecords {
			ids = append(ids, r.ID)
		}
		qs := core.GetModelQuerySeter(nil, new(models.ServiceRecord), true).Filter("id__in", ids)
		qs, fields, subTotal, _ := c.SetQuerySeterByURIParam(qs)
		total, err := qs.OrderBy("-startedAt", "-area").All(&serviceRecords)
		if err != nil || total == 0 {
			serviceRecords = []*models.ServiceRecord{}
		}

		if helpers.ContainsString(fields, "filter") {
			total = subTotal
		} else {
			total = int64(len(ids))
		}

		c.Success(total, serviceRecords)
		return
	}
	c.Error(errors.New("UNSUPPORTED_URL_PARAM"))
	return
}

// Put Put
func (c *ServiceRecordController) Put() {
	user, err := c.GetAuthUser()
	if err != nil || user.Role == "public" {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = c.ParseJSONBodyStruct(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	existing := new(models.ServiceRecord)
	if err := core.GetModelQuerySeter(nil, existing, true).
		Filter("id", serviceRecord.ID).
		One(existing); err != nil {
		c.Error(err)
		return
	}

	if existing.Area != serviceRecord.Area {
		if serviceRecord.NumOfRecordsInTheSameDayForTheSameArea() > 0 {
			c.Error(errors.New("RECORD_ALREADY_EXISTS"))
			return
		}
	}

	err = c.PutModel(nil, serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}

// PostDetail -
func (c *ServiceRecordController) PostDetail() {
	user, err := c.GetAuthUser()
	if err != nil || user.Role == "public" {
		c.Error(err)
		return
	}

	recordDetail := new(models.RecordDetail)
	c.ParseJSONBodyStruct(recordDetail)
	_, err = core.InsertModel(nil, recordDetail)
	if err != nil {
		c.Error(err)
		return
	}
	c.Success(1, "success")
}

// PutDetail -
func (c *ServiceRecordController) PutDetail() {
	user, err := c.GetAuthUser()
	if err != nil || user.Role == "public" {
		c.Error(err)
		return
	}

	err = c.PutModel(nil, &models.RecordDetail{})
	if err != nil {
		c.Error(err)
		return
	}
	c.Success(1, "success")
}

// DeleteDetail -
func (c *ServiceRecordController) DeleteDetail() {
	user, err := c.GetAuthUser()
	if err != nil || user.Role == "public" {
		c.Error(err)
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	_, err = core.GetModelQuerySeter(nil, &models.RecordDetail{}, false).Filter("id", id).Delete()
	if err != nil {
		c.Error(err)
		return
	}
	c.Success(1, "success")
}
