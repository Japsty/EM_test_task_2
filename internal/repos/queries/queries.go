package queries

const (
	// USER QUERIES---------------------------------

	ExistCheck = `
		SELECT EXISTS(
		SELECT 1 
		FROM users 
		WHERE id = $1)
	`

	CreateUser = `
		INSERT INTO users (passport_number, surname, name, patronymic, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	FindUserByID = `
		SELECT id, passport_number, surname, name, patronymic, address
		FROM users
		WHERE id = $1;
	`

	UpdateUser = `
		UPDATE users
		SET surname= $2,name= $3,patronymic= $4,address= $5
		WHERE id = $1
		RETURNING id, passport_number, surname, name, patronymic, address;
	`

	DeleteUser = `
		DELETE FROM users
		WHERE id = $1;
	`

	//----------------------------------------------

	// TASKS QUERIES---------------------------------

	CreateTask = `
		INSERT INTO tasks (name, user_id,start_time,end_time)
        VALUES ($1, $2, NULL, NULL)
        RETURNING id, name, user_id;
	`

	FindTaskByID = `
		SELECT id, name, user_id, start_time, end_time
		FROM tasks
		WHERE id = $1;
	`

	DeleteTask = `
		DELETE FROM tasks
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

	GetAllTasks = `
		SELECT id, name, user_id, start_time, end_time
		FROM tasks
		ORDER BY user_id DESC;
	`

	//----------------------------------------------
)
