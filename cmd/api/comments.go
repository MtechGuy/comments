package main

import (
	"errors"
	"fmt"
	"net/http"

	// import the data package which contains the definition for Comment
	"github.com/mtechguy/comments/internal/data"
	"github.com/mtechguy/comments/internal/validator"
)

func (a *applicationDependencies) createCommentHandler(w http.ResponseWriter, r *http.Request) {
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
	err = a.commentModel.Insert(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData) // delete this
	// Set a Location header. The path to the newly created comment
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/comments/%d", comment.ID))

	data := envelope{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}

func (a *applicationDependencies) displayCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL /v1/comments/:id so that we
	// can use it to query teh comments table. We will
	// implement the readIDParam() function later
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Call Get() to retrieve the comment with the specified id
	comment, err := a.commentModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// display the comment
	data := envelope{
		"comment": comment,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}
