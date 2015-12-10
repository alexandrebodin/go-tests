package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

type Contact struct {
	Id                int        `json:"id"`
	Status            nullString `json:"status"`
	GroupName         nullString `json:"group_name"`
	AccountName       nullString `json:"account_name"`
	Lastname          nullString `json:"lastname"`
	Firstname         nullString `json:"firstname"`
	IndustryName      nullString `json:"industry_name"`
	MetierName        nullString `json:"metier_name"`
	JobName           nullString `json:"job_name"`
	Owner             nullString `json:"owner"`
	SiteCity          nullString `json:"site_city"`
	SitePostalCode    nullString `json:"site_postal_code"`
	ContactTelephone  nullString `json:"contact_telephone"`
	ContactTelephone2 nullString `json:"contact_telephone2"`
	TelephonePrefix   nullString `json:"telephone_prefix"`
	TelephonePrefix2  nullString `json:"telephone2_prefix"`
}

type nullString struct {
	sql.NullString
}

func (s nullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Request received")
	t0 := time.Now()

	db, err := sql.Open(
		"mysql",
		"root:vagrant@tcp(:3306)/Vamos")

	if err != nil {
		w.Write([]byte("Error1"))
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		w.Write([]byte("noPing"))
		return
	}

	s := `
        SELECT 
            contacts.id, 
            contacts.status, 
            groups.name as group_name, 
            accounts.usual_name as account_name, 
            humans.lastname, 
            humans.firstname, 
            industries.name as industry_name, 
            metiers.name as metier_name,
            jobs.name as job_name, 
            manager_human.trigram as owner, 
            sites.city as site_city, 
            sites.postal_code as site_postal_code, 
            humans.telephone as contact_telephone, 
            tel1.phonecode as telephone_prefix,
            tel2.phonecode as telephone2_prefix,
            humans.telephone2 as contact_telephone2
        FROM contacts
        LEFT JOIN humans ON (contacts.human_id=humans.id)
        LEFT JOIN telephone tel1 ON (humans.prefix_telephone = tel1.id)
        LEFT JOIN telephone tel2 ON (humans.prefix_telephone2 = tel2.id)
        LEFT JOIN telephone tel3 ON (humans.prefix_telephone3 = tel3.id)
        LEFT JOIN accounts ON (accounts.id=contacts.account_id)
        LEFT JOIN groups ON (groups.id=accounts.group_id)
        LEFT JOIN industries ON (industries.id=accounts.industry_id)
        LEFT JOIN sites ON (sites.id=contacts.site_id)
        LEFT JOIN users as manager ON (manager.id=contacts.manager_id)
        LEFT JOIN humans as manager_human ON (manager_human.id=manager.human_id)
        LEFT JOIN metiers ON (metiers.id=contacts.metier_id)
        LEFT JOIN link_contact_job ON (link_contact_job.contact_id = contacts.id AND link_contact_job.main=1)
        LEFT JOIN jobs ON (link_contact_job.job_id = jobs.id)
        LEFT JOIN (SELECT contact_id, MAX(updated_at) as updated_at FROM mails GROUP BY contact_id) as mails ON (contacts.id=mails.contact_id)
        LEFT JOIN (SELECT contact_id, MAX(updated_at) as updated_at FROM task GROUP BY contact_id) as task ON (contacts.id=task.contact_id)
        LEFT JOIN (SELECT contact_id, MAX(updated_at) as updated_at FROM prospection_meeting GROUP BY contact_id) as prospection_meeting ON (contacts.id=prospection_meeting.contact_id)
        LIMIT 100
    `
	rows, err := db.Query(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	contacts := make([]Contact, 0, 1000)

	for rows.Next() {

		var c Contact
		err := rows.Scan(
			&c.Id,
			&c.Status,
			&c.GroupName,
			&c.AccountName,
			&c.Lastname,
			&c.Firstname,
			&c.IndustryName,
			&c.MetierName,
			&c.JobName,
			&c.Owner,
			&c.SiteCity,
			&c.SitePostalCode,
			&c.ContactTelephone,
			&c.ContactTelephone2,
			&c.TelephonePrefix,
			&c.TelephonePrefix2)

		if err != nil {
			fmt.Println(err)
			return
		}
		contacts = append(contacts, c)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	js, err := json.Marshal(contacts)

	if err != nil {
		w.Write([]byte("Error1"))
		return

	} else {

		w.Write(js)
		t1 := time.Now()
		log.Printf("The call took %v to run.\n", t1.Sub(t0))
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

}

func main() {

	//http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", handler)

	fmt.Println("listening on port 9000")
	http.ListenAndServe(":9000", nil)
}
