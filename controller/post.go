package controller

import (
	engine "github.com/JoanGTSQ/api"
	"neft.web/models"
)

type PostFuncs struct{}

// GetPost Obtain the post id compare to the user, and return it
func (p PostFuncs) GetPost(client *Client) {

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

// CreatePost Receive a JSON with a valid post and create it
func (p PostFuncs) CreatePost(client *Client) {

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

// RetrieveAllPost Return all posts
func (p PostFuncs) RetrieveAllPost(client *Client) {
	var mp *models.MultiplePost
	// Create pagination
	err := mp.Count()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	var links Links
	mp.Pagination, links = client.GeneratePaginationFromRequest(int(mp.Quantity))

	// Retrieve all posts data
	err = mp.AllPosts(client.User.ID)
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["data"] = mp.Posts
	client.LastMessage.Data["links"] = links
	client.SendMessage()
	return
}

// DeletePost Receive a valida post and delete it
func (p PostFuncs) DeletePost(client *Client) {
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

// UpdatePost Receive a valid post and update it
func (p PostFuncs) UpdatePost(client *Client) {
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

// CommentPost  Obtain the remember_hash from the JWT token and return it in JSON
func (p PostFuncs) CommentPost(client *Client) {
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

// UncommentPost Obtain the remember_hash from the JWT token and return it in JSON
func (p PostFuncs) UncommentPost(client *Client) {
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

// LikePost Obtain the remember_hash from the JWT token and return it in JSON
func (p PostFuncs) LikePost(client *Client) {

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

// UnlikePost Obtain the remember_hash from the JWT token and return it in JSON
func (p PostFuncs) UnlikePost(client *Client) {

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
