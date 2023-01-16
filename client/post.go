package client

import (
	"net/http"
	"strconv"

	engine "github.com/JoanGTSQ/api"
	"github.com/gin-gonic/gin"
	"neft.web/models"
)

type Posts struct {
	pdb models.PostDB
}

func NewPosts(db models.PostDB) *Posts {
	return &Posts{
		pdb: db,
	}
}

// GetPost GET /post/:id
// Obtain the post id compare to the user, and return it
func GetPost(id int) (*models.Post, error) {

	post := &models.Post{
		NeftModel: models.NeftModel{
			ID: uint(id),
		},
	}

	err := post.ByID()
	if err != nil {
		engine.Warning.Println(err)
		// handleError(err, context)
		return &models.Post{}, err
	}
	ResponseMap["data"] = post
	response = engine.Response{
		ResponseCode: http.StatusOK,
		// Context:      context,
		Response: ResponseMap,
	}
	return post, nil
	// response.SendAnswer()
}

// CreatePost PUT /post
// Receive a JSON with a valid post and create it
func (pts *Posts) CreatePost(context *gin.Context) {

	var post models.Post
	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(&post); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	// Check and return user from token
	user, err := getUserFromContext(context)
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// Create post
	post.UserID = user.ID

	if err := post.Create(); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	engine.Debug.Println("New post created")

	ResponseMap["data"] = &post
	response = engine.Response{
		ResponseCode: 201,
		Context:      context,
		Response:     ResponseMap,
	}
	response.SendAnswer()
}

// RetrieveAllPost GET /posts
// Return all posts from the user checked in the neftAuth token
func (pts *Posts) RetrieveAllPost(context *gin.Context) {

	// Create pagination
	count, err := pts.pdb.Count()
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	pagination, links := GeneratePaginationFromRequest(context, count)

	// Check and return user from token
	user, err := getUserFromContext(context)
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// Retrieve all posts data
	posts, err := pts.pdb.AllPosts(pagination, user.ID)
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	ResponseMap["data"] = posts
	ResponseMap["links"] = links
	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     ResponseMap,
	}
	response.SendAnswer()
}

// DeletePost DELETE /post
// Receive a valida post and delete it
func (pts *Posts) DeletePost(context *gin.Context) {
	post := &models.Post{}

	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(post); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	err := post.ByID()
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
	}

	if err := post.Delete(); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     nil,
	}
	response.SendAnswer()
}

// UpdatePost PATCH /post
// Receive a valid post and update it
func (pts *Posts) UpdatePost(context *gin.Context) {
	postReceived := &models.Post{}

	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(postReceived); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	if err := postReceived.Update(); err != nil {
		engine.Warning.Println(engine.ERR_NOT_SAME_USER)
		handleError(engine.ERR_NOT_SAME_USER, context)
		return
	}

	if err := postReceived.ByID(); err != nil {
		engine.Warning.Println(engine.ERR_NOT_SAME_USER)
		handleError(engine.ERR_NOT_SAME_USER, context)
		return
	}

	ResponseMap["data"] = postReceived

	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     ResponseMap,
	}
	response.SendAnswer()
}

// Comment  POST /user/comment/:id
// Obtain the remember_hash from the JWT token and return it in JSON
func (pts *Posts) Comment(context *gin.Context) {
	user, err := getUserFromContext(context)
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	postToLikeID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	post := &models.Post{
		NeftModel: models.NeftModel{
			ID: uint(postToLikeID),
		},
	}
	commentReceived := &models.Comment{}

	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(commentReceived); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	commentReceived.CommentatorID = user.ID
	if err = post.Comment(commentReceived); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// Return token and 200 Code
	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     nil,
	}
	response.SendAnswer()
}

// Uncomment  POST /user/uncomment/:id
// Obtain the remember_hash from the JWT token and return it in JSON
func (pts *Posts) Uncomment(context *gin.Context) {

	postToLikeID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	post := &models.Post{
		NeftModel: models.NeftModel{
			ID: uint(postToLikeID),
		},
	}
	commentReceived := &models.Comment{}

	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(commentReceived); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	if err = post.Uncomment(commentReceived); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// Return token and 200 Code
	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     nil,
	}
	response.SendAnswer()
}

// Like POST /user/like/:id
// Obtain the remember_hash from the JWT token and return it in JSON
func (pts *Posts) Like(context *gin.Context) {

	// Search the user from the claims by remember hash
	user, err := getUserFromContext(context)
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	postToLikeID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	post := &models.Post{
		NeftModel: models.NeftModel{
			ID: uint(postToLikeID),
		},
	}

	if err = post.Like(user.ID); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// Return token and 200 Code
	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     nil,
	}
	response.SendAnswer()
}

// Unlike POST /user/like/:id
// Obtain the remember_hash from the JWT token and return it in JSON
func (pts *Posts) Unlike(context *gin.Context) {

	// Search the user from the claims by remember hash
	user, err := getUserFromContext(context)
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	postToLikeID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}
	post := &models.Post{
		NeftModel: models.NeftModel{
			ID: uint(postToLikeID),
		},
	}

	if err = post.Unlike(user.ID); err != nil {
		engine.Warning.Println(err)
		handleError(err, context)
		return
	}

	// Return token and 200 Code
	response = engine.Response{
		ResponseCode: http.StatusOK,
		Context:      context,
		Response:     nil,
	}
	response.SendAnswer()
}
