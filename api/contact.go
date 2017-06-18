package api

type Contact struct {
	Id        int      `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Phones    []string `json:"phones"`
}

func (c *Contact) delete() {
	stmt, err := conn.Prepare("DELETE FROM contacts WHERE id = ?")
	checkErr(err)
	defer stmt.Close()

	stmt.Exec(c.Id)
}

func (c *Contact) create() int {
	stmt, err := conn.Prepare("INSERT INTO contacts(first_name, last_name) VALUES(?, ?)")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(c.FirstName, c.LastName)
	if err != nil {
		//error handling
	}
	id, _ := res.LastInsertId()
	return int(id)
}

func (c *Contact) update() {
	stmt, err := conn.Prepare("UPDATE contacts SET first_name = ?, last_name = ? WHERE id = ?")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(c.FirstName, c.LastName, c.Id)
	if err != nil {
		//error handling
	}
	_ = res
}
