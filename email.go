package main

import (
	"encoding/json"
	"net/http"
)

// Post new email templates
// [POST] /api/v1/email
func PostEmailTemplate(w http.ResponseWriter, r *http.Request) {
	e := &EmailTemplateRequest{}
	err := e.Deserialize(r)

	if err != nil {
		WriteResponse(w, 400, &Response{Error: "Invalid JSON"})
		return
	}

	if err := Validator.Struct(e.Template); err != nil {
		WriteValidationErrorResponse(w, err)
		return
	}

	tx := App.db.Begin()

	if err := tx.Create(&e.Template).Error; err != nil {
		WriteResponse(w, 422, &Response{Error: "Unable to save email template"})
		tx.Rollback()
		return
	}

	for k := range e.Tags {
		if err := tx.Where(&e.Tags[k]).FirstOrCreate(&e.Tags[k]).Error; err != nil {
			WriteResponse(w, 422, &Response{Error: "Error upserting tags"})
			tx.Rollback()
			return
		} else {
			e.Tags[k].EmailTemplate = e.Template
			tx.Save(e.Tags[k])
		}
	}

	tx.Commit()

	jsonStr, _ := json.Marshal(&e)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}

// Query Email template
// [GET] /api/v1/email
func GetEmailTemplate(w http.ResponseWriter, r *http.Request) {
	var e []EmailTemplate

	p := GetPaginationFromRequest(r)

	App.db.Limit(p.Limit).
		Offset(p.Offset).
		Find(&e)

	jsonStr, _ := json.Marshal(&e)
	WriteResponseHeader(w, 200)
	w.Write(jsonStr)
}
