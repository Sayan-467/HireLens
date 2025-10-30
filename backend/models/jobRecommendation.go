package models

import "time"

type JobRecommendation struct {
	Id          uint      `gorm:"primaryKey" json:"id"`
	ResumeId    uint      `json:"resume_id"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	Description string    `gorm:"type:text" json:"description"`
	Salary      string    `json:"salary"`
	JobUrl      string    `json:"job_url"`
	PostedDate  string    `json:"posted_date"`
	JobType     string    `json:"job_type"`
	CreatedAt   time.Time `json:"created_at"`
	Resume      Resume    `gorm:"foreignKey:ResumeId"`
}
