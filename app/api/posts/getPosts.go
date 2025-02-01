package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/app/api/auth"
	"forum/app/api/comments"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
)

func GetPosts(resp http.ResponseWriter, req *http.Request, db *sql.DB) []byte {
	var posts []models.Post
	var categories []string
	var err error
	pageStr := ""
	config.Logger.Printf("Request URL Path: %s", req.URL.Path)
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) >= 4 {
		pageStr = pathParts[3]
		config.Logger.Printf("Requested Page: %s", pageStr)
	}
	page := 1

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if strings.Contains(req.URL.Path, "categories") {
		cat := strings.Split(req.URL.Path, "/categories=")[1]
		categories = strings.Split(cat, "&")
		config.Logger.Printf("Categories found in URL: %v", categories)
	}

	if len(categories) > 0 {
		config.Logger.Printf("Fetching posts by categories: %v for page %d", categories, page)
		err = fetchPostsByCategories(categories, &posts, page, db)
	} else if strings.Contains(req.URL.Path, "created") {
		if !auth.SessionCheck(resp, req, db) {
			models.SendErrorResponse(resp, http.StatusUnauthorized, "Access Unauthorized")
			return []byte("")
		}
		token, _ := utils.GetSessionToken(req)
		_, userName, _ := utils.GetUsernameByToken(token, db)
		err = fetchUserCreatedPosts(&posts, userName, db)
	} else if strings.Contains(req.URL.Path, "liked") {
		if !auth.SessionCheck(resp, req, db) {
			models.SendErrorResponse(resp, http.StatusUnauthorized, "Access Unauthorized")
			return []byte("null")
		}
		token, _ := utils.GetSessionToken(req)
		userID, _, _ := utils.GetUsernameByToken(token, db)
		fmt.Println("BIG XXXXX")
		err = fetchUserLikedPosts(&posts, userID, db)
	} else if req.URL.Path == "/api/posts" || strings.HasPrefix(req.URL.Path, "/api/posts/") {
		config.Logger.Printf("Fetching posts for page %d", page)
		err = fetchAllPosts(&posts, page, db)
	} else {
		return []byte("null")
	}

	if err != nil {
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return nil
	}
	var postJSON []byte
	if len(posts) == 0 {
		return []byte{}
	} else {
		postJSON, err = json.Marshal(posts)
		if err != nil {
			models.SendErrorResponse(resp, 500, "Error: Internal Server Error")
			config.Logger.Printf("Error marshaling posts: %v", err)
			return nil
		}

	}
	config.Logger.Printf("Number of posts returned: %d", len(posts))

	return postJSON
}

func fetchPostsByCategories(categories []string, posts *[]models.Post, page int, db *sql.DB) error {
	const limit = 10
	offset := (page - 1) * limit

	// Prepare the LIKE conditions for each category
	var likeConditions []string
	for _, category := range categories {
		config.Logger.Printf("Category Found: %s", category)
		likeConditions = append(likeConditions, "Categories LIKE ?")
	}

	condition := strings.Join(likeConditions, " OR ")

	query := fmt.Sprintf(`
		SELECT username, id, title, content, Categories, created_at, likes, dislikes 
		FROM posts 
		WHERE %s
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`, condition)

	params := make([]any, len(categories)+2)
	for i, category := range categories {
		params[i] = "%" + category + "%"
	}
	// add limit and offset still have to see ila ghytzado ola la, pagination o dakchi
	params[len(categories)] = limit
	params[len(categories)+1] = offset

	config.Logger.Println("Params: ", params)

	config.Logger.Printf("Executing query: %s with categories: %v, limit: %d, offset: %d", query, categories, limit, offset)

	rows, err := db.Query(query, params...)
	if err != nil {
		log.Printf("Error querying posts by category: %v", err)
		return fmt.Errorf("error querying posts by category: %v", err)
	}
	defer rows.Close()

	return rowsProcess(rows, posts, db)
}

func fetchAllPosts(posts *[]models.Post, page int, db *sql.DB) error {
	const limit = 10
	offset := (page - 1) * limit

	query := `SELECT username, id, title, content, Categories, created_at, likes, dislikes 
              FROM posts 
              ORDER BY created_at DESC 
              LIMIT ? OFFSET ?`

	config.Logger.Printf("Executing query: %s with limit %d and offset %d", query, limit, offset)

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		config.Logger.Printf("Error querying all posts: %v", err)
		return fmt.Errorf("error querying all posts: %v", err)
	}
	defer rows.Close()

	return rowsProcess(rows, posts, db)
}

func rowsProcess(rows *sql.Rows, posts *[]models.Post, db *sql.DB) error {
	config.Logger.Printf("Starting to process rows...")

	for rows.Next() {
		var post models.Post

		config.Logger.Printf("Scanning row for post...")
		var creationTime time.Time
		if err := rows.Scan(&post.Username, &post.ID, &post.Title, &post.Content, &post.Categories, &creationTime, &post.Likes, &post.Dislikes); err != nil {
			log.Printf("Error scanning post: %v", err)
			return fmt.Errorf("error scanning post: %v", err)
		}
		post.CreatedAt = utils.TimeAgo(creationTime)
		config.Logger.Println("SSSSS: ", creationTime, post.CreatedAt)
		config.Logger.Printf("Successfully scanned post with ID: %d", post.ID)

		var err error
		config.Logger.Printf("Fetching comments for post ID: %d", post.ID)
		post.Comments, err = comments.GetComments(post.ID, db)

		if err != nil {
			log.Printf("Error fetching comments for post ID %d: %v", post.ID, err)
			return fmt.Errorf("error fetching comments for post ID %d: %v", post.ID, err)
		}
		config.Logger.Printf("Successfully fetched comments for post ID: %d, Comments: %v", post.ID, post.Comments)

		config.Logger.Printf("Appending post with ID: %d to posts slice", post.ID)
		*posts = append(*posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return err
	}

	config.Logger.Printf("Finished processing rows. Total posts fetched: %d", len(*posts))
	return nil
}
