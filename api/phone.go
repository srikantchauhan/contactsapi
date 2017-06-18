package api

import "strconv"

type Phone struct {
	Id        int    `json:"id"`
	Number    string `json:"number"`
	ContactId int    `json:"-"`
}

func (p *Phone) create() {
	//TODO - validation on phone number
	stmt, err := conn.Prepare("INSERT INTO phones(number, id_contact) VALUES(?, ?)")
	checkErr(err)
	defer stmt.Close()
	stmt.Exec(p.Number, p.ContactId)
}

func getPhoneByContactId(id int) []string {
	res, err := conn.Query("SELECT number FROM phones WHERE id_contact = " + strconv.Itoa(id))
	if err != nil {
		//error handling
	}
	phones := make([]string, 0)
	for res.Next() {
		var number string

		res.Scan(&number)
		phones = append(phones, number)
	}
	return phones
}
