package apis

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/cron_task"
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/dto"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
	"os"
	"path"
	"strconv"
)

// CreateScheduleTask
// 创建定时任务
func CreateScheduleTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	// 1. 从c中取出数据
	dto2 := dto.CronTaskCreateDTO{}
	err := c.ShouldBind(&dto2)
	if err != nil {
		return nil, err
	}
	// 2. 新增到数据库
	task, err := service.CreateScheduleTask(&dto2)
	if err != nil {
		return nil, err
	}
	// 3. 保存文件
	dir := path.Join("cron_assets", strconv.Itoa(int(task.ID)))
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return nil, err
	}
	filepath := path.Join(dir, dto2.File.Filename)
	err = c.SaveUploadedFile(dto2.File, filepath)
	if err != nil {
		return nil, err
	}
	// 4. 更新数据库
	updateTask := model.MCronTask{
		Command: filepath,
		WorkDir: dir,
	}
	updateTask.ID = task.ID

	db := sqlitedao.Sqlite.DB()
	tx := db.Updates(&updateTask)
	if tx.Error != nil {
		return nil, err
	}

	return 0, nil
}

func DeleteScheduleTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	dto := model.MCronTask{}
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}
	err = service.DeleteScheduleTask(dto.ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func PageScheduleTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	page := model.PageRequest{}
	var err error
	page.Current, err = strconv.Atoi(c.DefaultQuery("current", "1"))
	if err != nil {
		return nil, err
	}
	page.Size, err = strconv.Atoi(c.DefaultQuery("size", "25"))
	if err != nil {
		return nil, err
	}

	condition := model.MCronTask{}
	scheduleTask, err := service.PageScheduleTask(page, condition)
	return scheduleTask, err
}

func UpdateScheduleTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	dto2 := dto.CronTaskUpdateDTO{}
	err := c.ShouldBind(&dto2)
	if err != nil {
		return nil, err
	}

	task, err := service.UpdateScheduleTask(&dto2)
	if err != nil {
		return nil, err
	}
	if dto2.File != nil {
		dir := path.Join("cron_assets", strconv.Itoa(int(task.ID)))
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return nil, err
		}
		filepath := path.Join(dir, dto2.File.Filename)
		err = c.SaveUploadedFile(dto2.File, filepath)
		if err != nil {
			return nil, err
		}
		// 4. 更新数据库
		updateTask := model.MCronTask{
			Command: filepath,
			WorkDir: dir,
		}
		updateTask.ID = task.ID

		db := sqlitedao.Sqlite.DB()
		tx := db.Updates(&updateTask)
		if tx.Error != nil {
			return nil, err
		}
	}
	return task, err
}

func EnableTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	id, ok := c.GetQuery("id")
	if !ok {
		return nil, errors.New("id must not be null")
	}
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	// 0. 更新数据库
	db := sqlitedao.Sqlite.DB()
	task := model.MCronTask{}
	task.ID = uint(idNum)
	task.Enable = "1"
	tx := db.Updates(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	find := db.Find(&task)
	if find.Error != nil {
		return nil, find.Error
	}

	// 1. 调用cron的库进行调度
	manager := cron_task.GetCronManager()
	err = manager.AddTask(task)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func DisableTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	id, ok := c.GetQuery("id")
	if !ok {
		return nil, errors.New("id must not be null")
	}
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	// 0. 更新数据库
	db := sqlitedao.Sqlite.DB()
	task := model.MCronTask{}
	task.ID = uint(idNum)
	task.Enable = "0"
	tx := db.Updates(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 1. 调用cron的库进行调度
	manager := cron_task.GetCronManager()
	manager.DeleteTask(uint(idNum))
	return nil, nil
}

func ListRunningTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	// 1. 请求cronManager获取正在运行的列表
	manager := cron_task.GetCronManager()
	tasks := manager.ListRunningTask()
	return tasks, nil
}

func TerminateRunningTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	// 1. 请求cronManager进行停止正在运行的任务
	id, ok := c.GetQuery("id")
	if !ok {
		return nil, errors.New("id must not be null")
	}
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	manager := cron_task.GetCronManager()
	err = manager.KillTask(idNum)
	return nil, nil
}