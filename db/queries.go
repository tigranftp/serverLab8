package db

const (
	SelectAllCountries = `select * from country `
	AddCountryQuery    = `INSERT INTO country (country_name)
						  VALUES ((?));`
	DeleteUniversityQuery = `DELETE FROM university_year WHERE university_id = (SELECT id FROM university WHERE university_name = (?));
						  DELETE FROM university_ranking_year WHERE university_id = (SELECT id FROM university WHERE university_name = (?));
						  DELETE FROM university WHERE university_name = (?);`
	AddUniversityQuery = `INSERT INTO university (country_id, university_name) 
							VALUES ((SELECT id from country where country_name == (?)), (?))`
	DeleteRankingCriteriaQuery = `DELETE FROM university_ranking_year WHERE ranking_criteria_id = (SELECT id FROM ranking_criteria WHERE ranking_system_id = (?) AND criteria_name = (?));
									DELETE FROM  ranking_criteria WHERE ranking_system_id = (?) AND criteria_name = (?)`
	ChangeUniversityYearStaffRatio = `UPDATE university_year
									SET student_staff_ratio = (?)
									WHERE university_id = (SELECT id from university where university_name = (?)) AND year = (?);`
	AddUniversityRankingYear = `INSERT INTO university_ranking_year (university_id, ranking_criteria_id, year, score) 
								VALUES ((SELECT id FROM university WHERE university_name = (?)),
										(SELECT id FROM ranking_criteria WHERE criteria_name = (?)),(?),(?))`
	CreateUserQuery = `INSERT INTO users (Name, Username, Password, Role)
						VALUES ((?),(?),(?), (?))`
	GetUserQuery = `SELECT * FROM users
						WHERE Username = (?) AND Password = (?)`
	UpdateRefreshQuery = `UPDATE users
						  SET RefreshToken = (?), RefreshTokenEAT = (?)
						  WHERE Id = (?);`
	GetUserByRefreshTokenQuery = `SELECT * FROM users
								  WHERE RefreshToken = (?)`
)
