package controller

import (
	engine "github.com/JoanGTSQ/api"
	"neft.web/models"
)

// GetPost
// Obtain the post id compare to the user, and return it
func (client *Client) GetPost() {

	post, err := client.GetPostFromRequest()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	err = post.ByID()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["post"] = post
	client.SendMessage()
}

// CreatePost PUT /post
// Receive a JSON with a valid post and create it
func (client *Client) CreatePost() {

	post, err := client.GetPostFromRequest()

	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	// Create post
	post.UserID = client.User.ID

	if err := post.Create(); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	engine.Debug.Println("New post created")

	client.LastMessage.Data["post"] = post
	client.SendMessage()
}

// RetrieveAllPost GET /posts
// Return all posts from the user checked in the neftAuth token
//func (pts *Posts) RetrieveAllPost() {
//
//	// Create pagination
//	count, err := pts.pdb.Count()
//	if err != nil {
//		engine.Warning.Println(err)
//		handleError(err, context)
//		return
//	}
//	pagination, links := GeneratePaginationFromRequest(context, count)
//
//	// Check and return user from token
//	user, err := getUserFromContext(context)
//	if err != nil {
//		engine.Warning.Println(err)
//		handleError(err, context)
//		return
//	}
//
//	// Retrieve all posts data
//	posts, err := pts.pdb.AllPosts(pagination, user.ID)
//	if err != nil {
//		engine.Warning.Println(err)
//		handleError(err, context)
//		return
//	}
//
//	ResponseMap["data"] = posts
//	ResponseMap["links"] = links
//	response = engine.Response{
//		ResponseCode: http.StatusOK,
//		Context:      context,
//		Response:     ResponseMap,
//	}
//	response.SendAnswer()
//}

// DeletePost DELETE /post
// Receive a valida post and delete it
func (client *Client) DeletePost() {
	post := &models.Post{}

	// Obtain the body in the request and parse to the user
	if err := client.GetInterfaceFromMap("post", post); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	err := post.ByID()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := post.Delete(); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "Post deleted!"
	client.SendMessage()
}

// UpdatePost PATCH /post
// Receive a valid post and update it
func (client *Client) UpdatePost() {
	post := &models.Post{}

	// Obtain the body in the request and parse to the user
	if err := client.GetInterfaceFromMap("post", post); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := post.Update(); err != nil {
		engine.Warning.Println(engine.ERR_NOT_SAME_USER)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := post.ByID(); err != nil {
		engine.Warning.Println(engine.ERR_NOT_SAME_USER)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["error"] = post
	client.SendMessage()
}

// Comment  POST /user/comment/:id
// Obtain the remember_hash from the JWT token and return it in JSON
func (client *Client) CommentPost() {
	post := &models.Post{}
	err := client.GetInterfaceFromMap("post", post)
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	if err := post.ByID(); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	comment := &models.Comment{}

	// Obtain the body in the request and parse to the user
	if err := client.GetInterfaceFromMap("comment", comment); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	comment.CommentatorID = client.User.ID
	if err = post.Comment(comment); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "Comment created successfully!"
	client.SendMessage()
}

// Uncomment  POST /user/uncomment/:id
// Obtain the remember_hash from the JWT token and return it in JSON
func (client *Client) UncommentPost() {
	post := &models.Post{}

	if err := client.GetInterfaceFromMap("post", post); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	comment := &models.Comment{}

	// Obtain the body in the request and parse to the user
	if err := client.GetInterfaceFromMap("comment", comment); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	if err := post.Uncomment(comment); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "Comment deleted successfully"
	client.SendMessage()
}

// Like POST /user/like/:id
// Obtain the remember_hash from the JWT token and return it in JSON
func (client *Client) LikePost() {

	// Search the user from the claims by remember hash

	post := &models.Post{}

	if err := client.GetInterfaceFromMap("post", post); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := post.Like(client.User.ID); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "Like successfully"
	client.SendMessage()
}

// Unlike POST /user/like/:id
// Obtain the remember_hash from the JWT token and return it in JSON
func (client *Client) UnlikePost() {

	post := &models.Post{}

	if err := client.GetInterfaceFromMap("post", post); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := post.Unlike(client.User.ID); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "Post unliked successfully"
	client.SendMessage()
}
