package controller

import (
	engine "github.com/JoanGTSQ/api"
	"neft.web/models"
)

type SpotFuncs struct{}

// GetSpot Obtain the spot id compare to the user, and return it
func (p SpotFuncs) GetSpot(client *Client) {
	var spot models.Spot
	err := client.GetInterfaceFromMap("spot", &spot)
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	err = spot.ByID()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["spot"] = spot
	client.SendMessage()
}

// CreateSpot Receive a JSON with a valid spot and create it
func (p SpotFuncs) CreateSpot(client *Client) {

	var spot models.Spot
	err := client.GetInterfaceFromMap("spot", &spot)
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	// Create spot
	spot.UserID = client.User.ID

	if err := spot.Create(); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	engine.Debug.Println("New spot created")

	client.LastMessage.Data["spot"] = spot
	client.LastMessage.Data["message"] = "New spot created!"
	client.SendMessage()
}

// RetrieveAllSpot Return all spots
func (p SpotFuncs) RetrieveAllSpot(client *Client) {
	mp := &models.MultipleSpots{}
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

	// Retrieve all spots data
	err = mp.AllSpots()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["data"] = mp.Spots
	client.LastMessage.Data["links"] = links
	client.SendMessage()
	return
}

// DeleteSpot Receive a valida spot and delete it
func (p SpotFuncs) DeleteSpot(client *Client) {
	spot := &models.Spot{}

	// Obtain the body in the request and parse to the user
	if err := client.GetInterfaceFromMap("spot", spot); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	err := spot.ByID()
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := spot.Delete(); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "Spot deleted!"
	client.SendMessage()
}

// UpdateSpot Receive a valid spot and update it
func (p SpotFuncs) UpdateSpot(client *Client) {
	spot := &models.Spot{}

	// Obtain the body in the request and parse to the user
	if err := client.GetInterfaceFromMap("spot", spot); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := spot.Update(); err != nil {
		engine.Warning.Println(engine.ERR_NOT_SAME_USER)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := spot.ByID(); err != nil {
		engine.Warning.Println(engine.ERR_NOT_SAME_USER)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["error"] = spot
	client.SendMessage()
}

// CommentSpot  Obtain the remember_hash from the JWT token and return it in JSON
func (p SpotFuncs) CommentSpot(client *Client) {
	spot := &models.Spot{}
	err := client.GetInterfaceFromMap("spot", spot)
	if err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	if err := spot.ByID(); err != nil {
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
	if err = spot.Comment(comment); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "Comment created successfully!"
	client.SendMessage()
}

// UncommentSpot Obtain the remember_hash from the JWT token and return it in JSON
func (p SpotFuncs) UncommentSpot(client *Client) {
	spot := &models.Spot{}

	if err := client.GetInterfaceFromMap("spot", spot); err != nil {
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
	if err := spot.Uncomment(comment); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "Comment deleted successfully"
	client.SendMessage()
}

// LikeSpot Obtain the remember_hash from the JWT token and return it in JSON
func (p SpotFuncs) LikeSpot(client *Client) {

	// Search the user from the claims by remember hash

	spot := &models.Spot{}

	if err := client.GetInterfaceFromMap("spot", spot); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := spot.Like(client.User.ID); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["spot"] = spot
	client.LastMessage.Data["message"] = "Like successfully"
	client.SendMessage()
}

func (p SpotFuncs) RetrieveSingle(client *Client) {
	// Search the user from the claims by remember hash

	spot := &models.Spot{}

	if err := client.GetInterfaceFromMap("spot", spot); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := spot.ByID(); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}
	client.LastMessage.Data["spot"] = spot
	client.SendMessage()
}

// UnlikeSpot Obtain the remember_hash from the JWT token and return it in JSON
func (p SpotFuncs) UnlikeSpot(client *Client) {

	spot := &models.Spot{}

	if err := client.GetInterfaceFromMap("spot", spot); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	if err := spot.Unlike(client.User.ID); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return
	}

	client.LastMessage.Data["message"] = "Spot unliked successfully"
	client.SendMessage()
}
