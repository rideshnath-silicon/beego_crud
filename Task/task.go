package task

import (
	"CarCrudv2/helpers"
	"CarCrudv2/models"
	"context"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/task"
)

func CreateTask(taskName, schedule string, f task.TaskFunc) {
	tasks := task.NewTask(taskName, schedule, f)
	task.AddTask(taskName, tasks)
	task.StartTask()
}
func Demo(c context.Context) error {
	log.Print("hello")
	return nil
}

func SendPendingEmail(c context.Context) error {
	o := orm.NewOrm()
	// orm.Debug = true
	var emails []models.EmailLogs
	_, err := o.QueryTable(new(models.EmailLogs)).Filter("status", "pending").All(&emails, "LogId", "emailTo", "name", "subject", "body")
	if err != nil {
				return err
	}
	for _, email := range emails {
		sent, _ := helpers.SendMailOTp(email.To, email.Name, email.Subject, email.Body)
		if sent {
			var UpdateEmail = models.EmailLogs{Id: email.Id, Status: "success"}
			_, err := o.Update(&UpdateEmail, "status")
			if err != nil {
				return err
			}			
		}
	}
	return nil
}
