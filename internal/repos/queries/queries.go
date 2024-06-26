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

	FindUserByHash = `
		SELECT id, passport_hash, surname, name, patronymic, address
		FROM users
		WHERE passport_hash = $1;
	`

	FindUserByID = `
		SELECT id, passport_hash, surname, name, patronymic, address
		FROM users
		WHERE id = $1;
	`

	UpdateUser = `
		UPDATE user
		SET passport_hash,surname,name,patronymic,address
	`

	DeleteUser = `
		DELETE FROM users
		WHERE id = $1;
	`

	//----------------------------------------------
)
