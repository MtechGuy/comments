package main

import (
	"fmt"
	"net/http"

	// import the data package which contains the definition for Comment
	"github.com/mtechguy/comments/internal/data"
	"github.com/mtechguy/comments/internal/validator"
)

func (a *applicationDependencies) createCommentHandler(w http.ResponseWriter,
	r *http.Request) {
	// create a struct to hold a comment
	// we use struct tags to make the names display in lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	// perform the decoding
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	comment := &data.Comment{
		Content: incomingData.Content,
		Author:  incomingData.Author,
	}
	// Initialize a Validator instance
	v := validator.New()

	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors) // implemented later
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}
