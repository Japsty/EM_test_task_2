package queries

const (
	// USER QUERIES---------------------------------

	GetAllUsers = `
		SELECT id, passport_hash, surname, name, patronymic, address
		FROM users
		WHERE 1=1
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
		UPDATE users
		SET surname= $2,name= $3,patronymic= $4,address= $5
		WHERE id = $1
		RETURNING id, passport_hash, surname, name, patronymic, address;
	`

	DeleteUser = `
		DELETE FROM users
		WHERE id = $1;
	`

	//----------------------------------------------

	// TASKS QUERIES---------------------------------

	CreateTask = `
		INSERT INTO tasks (name, user_id)
        VALUES ($1, $2)
        RETURNING id, name, user_id;
	`

	FindTaskByID = `
		SELECT id, name, user_id, start_time, end_time
		FROM tasks
		WHERE id = $1;
	`

	FindTasksByUserID = `
		SELECT id, name, user_id, start_time, end_time
		FROM tasks
		WHERE user_id = $1;
	`

	FindTasksByUserIDWithPeriod = `
		SELECT id, name, user_id, start_time, end_time
		FROM tasks
		WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
		ORDER BY (end_time - start_time) DESC;
	`

	DeleteTask = `
		DELETE FROM users
		WHERE id = $1;
	`

	StartTimeTracker = `
        UPDATE tasks
        SET start_time = $1
        WHERE id = $2 AND user_id = $3;
	`

	StopTimeTracker = `
        UPDATE tasks
        SET end_time = $1
        WHERE id = $2 AND user_id = $3;
	`

	//----------------------------------------------
)
