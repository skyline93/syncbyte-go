package policy

import (
	"fmt"
	"log"

	"github.com/skyline93/syncbyte-go/internal/engine/repository"
	"github.com/skyline93/syncbyte-go/internal/engine/scheduler"
	"github.com/skyline93/syncbyte-go/internal/pkg/types"
	"gorm.io/gorm"
)

type Policy struct {
	ID         uint
	ResourceID uint
	Retention  int
	IsCompress bool
	AgentID    uint
	Status     string
	Scheduler  *Scheduler
}

type Scheduler struct {
	PolicyID  uint
	Type      string
	Cron      string
	Interval  int
	StartTime types.LocalTime
	EndTime   types.LocalTime
	IsActive  bool
}

func LoadPolicy(db *gorm.DB, id uint) (*Policy, error) {
	var pl repository.BackupPolicy

	if result := db.Where("id = ?", id).First(&pl); result.Error != nil {
		return nil, result.Error
	}

	return &Policy{
		ID:         pl.ID,
		ResourceID: pl.ResourceID,
		Retention:  pl.Retention,
		IsCompress: pl.IsCompress,
		AgentID:    pl.AgentID,
		Status:     pl.Status,
		Scheduler: &Scheduler{
			PolicyID:  pl.ID,
			Type:      pl.ScheduleType,
			Cron:      pl.Cron,
			Interval:  pl.Frequency,
			StartTime: pl.StartTime,
			EndTime:   pl.EndTime,
			IsActive:  pl.IsActive,
		},
	}, nil
}

func GetPolicies(db *gorm.DB, policyIDs []uint) (pls []*Policy, err error) {
	var items []repository.BackupPolicy

	if result := db.Where("id IN ?", policyIDs).Find(&items); result.Error != nil {
		return nil, result.Error
	}

	for _, item := range items {
		pl := &Policy{
			ID:         item.ID,
			ResourceID: item.ResourceID,
			Retention:  item.Retention,
			IsCompress: item.IsCompress,
			AgentID:    item.AgentID,
			Status:     item.Status,
			Scheduler: &Scheduler{
				PolicyID:  item.ID,
				Type:      item.ScheduleType,
				Cron:      item.Cron,
				Interval:  item.Frequency,
				StartTime: item.StartTime,
				EndTime:   item.EndTime,
				IsActive:  item.IsActive,
			},
		}

		pls = append(pls, pl)
	}

	return pls, nil
}

func allocateAgent(agentID uint) uint {
	return agentID
}

func CreatePolicy(db *gorm.DB, resourceID, agentID uint, retention int, isCompress bool, schType string, cron string, interval int) (*Policy, error) {
	// TODO allocate agent
	agentID = allocateAgent(agentID)

	item := repository.BackupPolicy{
		ResourceID:   resourceID,
		Retention:    retention,
		ScheduleType: schType,
		Cron:         cron,
		Frequency:    interval,
		IsCompress:   isCompress,
		AgentID:      agentID,
	}

	if result := db.Create(&item); result.Error != nil {
		return nil, result.Error
	}

	pl := &Policy{
		ID:         item.ID,
		ResourceID: item.ResourceID,
		Retention:  item.Retention,
		IsCompress: item.IsCompress,
		AgentID:    item.AgentID,
		Status:     item.Status,
		Scheduler: &Scheduler{
			PolicyID:  item.ID,
			Type:      item.ScheduleType,
			Cron:      item.Cron,
			Interval:  item.Frequency,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			IsActive:  item.IsActive,
		},
	}

	scheduler.Sch.BackupScheduleChan <- pl.Scheduler
	return pl, nil
}

func (s *Scheduler) Activate(db *gorm.DB) error {
	if result := db.Model(&repository.BackupPolicy{}).Where("id = ?", s.PolicyID).Update("is_active", true); result.Error != nil {
		return result.Error
	}

	scheduler.Sch.BackupScheduleChan <- s

	return nil
}

func (s *Scheduler) Inactivate(db *gorm.DB) error {
	if result := db.Model(&repository.BackupPolicy{}).Where("id = ?", s.PolicyID).Update("is_active", false); result.Error != nil {
		return result.Error
	}

	scheduler.Sch.UnloadJobIDChan <- s.GetID()
	return nil
}

func (s *Scheduler) Run() {
	j := repository.ScheduledJob{
		JobType:        string(types.Backup),
		Status:         string(types.Queued),
		BackupPolicyID: s.PolicyID,
	}

	if result := repository.Db.Create(&j); result.Error != nil {
		log.Printf("create scheduled job error, msg: %v", result.Error)
		return
	}

	log.Printf("add scheduled job, <%d>", j.ID)
}

func (s *Scheduler) GetID() string {
	return fmt.Sprintf("backup-%d", s.PolicyID)
}

func (s *Scheduler) GetCron() string {
	return s.Cron
}

func LoadSchedulers() {
	var items []repository.BackupPolicy

	if result := repository.Db.Where("is_active = ?", true).Find(&items); result.Error != nil {
		return
	}

	for _, item := range items {
		pl := &Policy{
			ID:         item.ID,
			ResourceID: item.ResourceID,
			Retention:  item.Retention,
			IsCompress: item.IsCompress,
			AgentID:    item.AgentID,
			Status:     item.Status,
			Scheduler: &Scheduler{
				PolicyID:  item.ID,
				Type:      item.ScheduleType,
				Cron:      item.Cron,
				Interval:  item.Frequency,
				StartTime: item.StartTime,
				EndTime:   item.EndTime,
				IsActive:  item.IsActive,
			},
		}

		scheduler.Sch.BackupScheduleChan <- pl.Scheduler
	}
}
