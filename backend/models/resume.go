package models

import "time"

type Resume struct {
	Id             uint      `gorm:"primaryKey" json:"id"`
	UserId         uint      `json:"user_id"`
	Title          string    `json:"title"`
	FileUrl        string    `json:"file_url"` // appwrite storage url
	AnalysisResult string    `gorm:"type:jsonb" json:"analysis_result"`
	AtsScore       int       `gorm:"default:0" json:"ats_score"`
	JdMatchScore   int       `gorm:"default:0" json:"jd_match_score"`
	MatchingSkills string    `gorm:"type:jsonb" json:"matching_skills"` // JSON array of strings
	MissingSkills  string    `gorm:"type:jsonb" json:"missing_skills"`  // JSON array of strings
	UploadedAt     time.Time `json:"uploaded_at"`
	User           User      `gorm:"foreignKey:UserId"`
}
