statement, _ := db.Prepare(`
		CREATE TABLE IF NOT EXIST "users" (
			"roll no"	INTEGER NOT NULL UNIQUE,
			"coins"	INTEGER,
			PRIMARY KEY("roll no")
		);
	`)

statement.Exec()