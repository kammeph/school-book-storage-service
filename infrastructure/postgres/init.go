package postgres

var (
	createSchoolsTable = "CREATE TABLE IF NOT EXISTS schools (" +
		"id VARCHAR(100) NOT NULL," +
		"aggregate_id VARCHAR(100) NOT NULL," +
		"type VARCHAR(100) NOT NULL," +
		"version INTEGER NOT NULL," +
		"timestamp TIMESTAMP NOT NULL," +
		"data VARCHAR(255) NOT NULL," +
		"PRIMARY KEY (id))"
	createStoragesTable = "CREATE TABLE IF NOT EXISTS storages (" +
		"id VARCHAR(100) NOT NULL," +
		"aggregate_id VARCHAR(100) NOT NULL," +
		"type VARCHAR(100) NOT NULL," +
		"version INTEGER NOT NULL," +
		"timestamp TIMESTAMP NOT NULL," +
		"data VARCHAR(255) NOT NULL," +
		"PRIMARY KEY (id))"
	createSchoolClassesTable = "CREATE TABLE IF NOT EXISTS school_classes (" +
		"id VARCHAR(100) NOT NULL," +
		"aggregate_id VARCHAR(100) NOT NULL," +
		"type VARCHAR(100) NOT NULL," +
		"version INTEGER NOT NULL," +
		"timestamp TIMESTAMP NOT NULL," +
		"data VARCHAR(255) NOT NULL," +
		"PRIMARY KEY (id))"
	createBooksTable = "CREATE TABLE IF NOT EXISTS books (" +
		"id VARCHAR(100) NOT NULL," +
		"aggregate_id VARCHAR(100) NOT NULL," +
		"type VARCHAR(100) NOT NULL," +
		"version INTEGER NOT NULL," +
		"timestamp TIMESTAMP NOT NULL," +
		"data VARCHAR(255) NOT NULL," +
		"PRIMARY KEY (id))"
)

func CreatePostgresStoreTables() {
	db := NewDB()
	defer db.Close()
	_, err := db.Exec(createSchoolsTable)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(createStoragesTable)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(createSchoolClassesTable)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(createBooksTable)
	if err != nil {
		panic(err)
	}
}
