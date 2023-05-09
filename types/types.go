package types

import "database/sql"

type Country struct {
	Id          int
	CountryName string
}

type University struct {
	Id             int
	CountryName    string
	UniversityName string
}

type RankingCriteria struct {
	Id           int
	SystemID     int
	CriteriaName string
}

type ChangeStudentStaffRatio struct {
	UniversityName string
	Year           int
	NewStaffRatio  int
}

type AddUniversityRankingYear struct {
	UniversityName string
	CriteriaName   string
	Year           int
	Score          int
}

type User struct {
	Id              int64
	Name            string         `json:"name"`
	Username        string         `json:"username"`
	Password        string         `json:"password"`
	RefreshToken    sql.NullString `json:"refresh_token"`
	RefreshTokenEAT sql.NullInt64  `json:"refresh_token_eat"`
	Role            string         `json:"role"`
}
