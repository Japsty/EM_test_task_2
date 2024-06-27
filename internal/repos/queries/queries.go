package queries

const (
	// USER QUERRIES---------------------------------

	GetAllUsers = `
		SELECT id, passport_hash, surname, name, patronymic, address
		FROM users;
	`

	CreateUser = `
		INSERT INTO users (passport_hash, surname, name, patronymic, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	FindUserByID = `
		SELECT id, passport_hash, surname, name, patronymic, address
		FROM users
		WHERE id = $1;
	`

	UpdateUser = `
		UPDATE user
		SET passport_hash = $2,surname= $3,name= $4,patronymic= $5,address= $6
		WHERE id = $1
		RETURNING id, passport_hash,surname,name,patronymic,address;
	`

	DeleteUser = `
		DELETE FROM users
		WHERE id = $1;
	`

	//----------------------------------------------
)
