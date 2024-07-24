package db

import (
	"database/sql"
	"errors"
	"log"
	"rtf/pkg/structs"
)

// this function will be reused in the functions below
func postExists(id int) bool {
	var status bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE ID = $1)", id).Scan(&status)
	if err != nil {
		return false
	}
	return status
}

func userExists(username string) bool {
	var status bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE username = $1)", username).Scan(&status)
	if err != nil {
		return false
	}
	return status
}

func isOwner(post int, username string) bool {
	if username == "" || post < 1 {
		return false
	}
	if !postExists(post) || !userExists(username) {
		return false
	}
	var status bool
	//TO DO: db.QueryRow
	return status
}

func CreatePost(title, content, username string) (int, error) {
	var postID int
	err := db.QueryRow(`
		INSERT INTO post (Title, Content, Username) 
		VALUES ($1, $2, $3) RETURNING PostID`, title, content, username).Scan(&postID)
	if err != nil {
		return 0, err
	}
	return postID, nil
}

func DeletePost(id int, user string) error {
	if !postExists(id) || !isOwner(id, user) {
		return errors.New("post does not exist")
	}
	_, err := db.Exec("DELETE FROM post WHERE ID = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func Interact(post int, username string, interaction int) error {
	if !postExists(post) || !userExists(username) {
		return errors.New("post does not exist")
	}
	//TO DO: Check if the user didn't already interact with the post
	_, err := db.Exec("INSERT INTO interaction (PostID, Username, Interaction) VALUES ($1, $2, $3)", post, username, interaction)
	if err != nil {
		return err
	}
	return nil
}

func GetPostsByUser(username string) ([]structs.Post, error) {
	var posts []structs.Post
	rows, err := db.Query("SELECT PostID, Title, Content, CreatedDate, username FROM Post WHERE username = $1", username)
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var post structs.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedDate, &post.Username); err != nil {
			return posts, err
		}
		categoryRows, err := db.Query("SELECT CategoryName FROM Category c INNER JOIN PostCategory pc ON c.CategoryID = pc.CategoryID WHERE pc.PostID = $1", post.ID)
		if err != nil {
			return posts, err
		}
		defer categoryRows.Close()
		for categoryRows.Next() {
			var category string
			if err := categoryRows.Scan(&category); err != nil {
				return posts, err
			}
			post.Categories = append(post.Categories, category)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetAllPosts() []structs.Post {
	var posts []structs.Post
	dbMutex.Lock()
	defer dbMutex.Unlock()
	rows, err := db.Query(`
	SELECT 
	p.PostID, p.Title, p.Content, p.Username, 
	IFNULL(likes.likes, 0) as likes, 
	IFNULL(dislikes.dislikes, 0) as dislikes, 
	IFNULL(comments.comments, 0) as comments
	FROM post p
	LEFT JOIN (SELECT PostID, COUNT(*) as likes FROM interaction WHERE Kind = 1 GROUP BY PostID) likes 
	ON p.PostID = likes.PostID
	LEFT JOIN (SELECT PostID, COUNT(*) as dislikes FROM interaction WHERE Kind = 0 GROUP BY PostID) dislikes 
	ON p.PostID = dislikes.PostID
	LEFT JOIN (SELECT PostID, COUNT(*) as comments FROM comment GROUP BY PostID) comments
	ON p.PostID = comments.PostID
	ORDER BY p.CreatedDate DESC 
    `)
	if err != nil {
		log.Printf("Query error: %s", err)
		return []structs.Post{}
	}
	defer rows.Close()
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Username, &post.Likes, &post.Dislikes, &post.Comments)
		if err != nil {
			log.Printf("Scan error: %s", err)
			continue
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %s", err)
		return []structs.Post{}
	}
	// fmt.Printf("%+v\n", posts)
	return posts
}

func GetPostDetails(postId int) (structs.Post, structs.User, []structs.Comment, []structs.Interaction) {
	var (
		thisPost          structs.Post
		thisUser          structs.User
		theseComments     []structs.Comment
		theseInteractions []structs.Interaction
	)
	err := db.QueryRow("SELECT * FROM post WHERE id = $1", postId).Scan(&thisPost.ID, &thisPost.Title, &thisPost.Content, &thisPost.Categories)
	if err != nil {
		return structs.Post{}, structs.User{}, []structs.Comment{}, []structs.Interaction{}
	}
	return thisPost, thisUser, theseComments, theseInteractions
}

func InsertOrUpdateInteraction(postID int, username string, kind int) error {
	var existingKind int
	err := db.QueryRow("SELECT Kind FROM Interaction WHERE PostID = ? AND Username = ?", postID, username).Scan(&existingKind)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing interaction, insert a new one
			_, err = db.Exec(
				"INSERT INTO Interaction (PostID, Username, Kind) VALUES (?, ?, ?)",
				postID, username, kind,
			)
			if err != nil {
				log.Printf("InsertInteraction error: %s", err)
				return err
			}
		} else {
			// Some other error occurred
			log.Printf("Query error: %s", err)
			return err
		}
	} else { //if no errors
		// Existing interaction found, update it if necessary
		if existingKind != kind {
			_, err = db.Exec(
				"UPDATE Interaction SET Kind = ? WHERE PostID = ? AND Username = ?",
				kind, postID, username,
			)
			if err != nil {
				log.Printf("UpdateInteraction error: %s", err)
				return err
			}
		} else {
			//if the interaction is the same
			_, err = db.Exec(
				"DELETE FROM Interaction WHERE PostID = ? AND Username = ?",
				postID, username,
			)
			if err != nil {
				log.Printf("UpdateInteraction error: %s", err)
				return err
			}
		}
	}
	return nil
}

// GetNewestPosts retrieves posts ordered by creation date (newest first)
func GetNewestPosts() ([]structs.Post, error) {
	var posts []structs.Post
	rows, err := db.Query(`
		SELECT 
			p.PostID, p.Title, p.Content, p.Username, 
			IFNULL(likes.likes, 0) as likes, 
			IFNULL(dislikes.dislikes, 0) as dislikes, 
			IFNULL(comments.comments, 0) as comments
		FROM post p
		LEFT JOIN (SELECT PostID, COUNT(*) as likes FROM interaction WHERE Kind = 1 GROUP BY PostID) likes 
		ON p.PostID = likes.PostID
		LEFT JOIN (SELECT PostID, COUNT(*) as dislikes FROM interaction WHERE Kind = 0 GROUP BY PostID) dislikes 
		ON p.PostID = dislikes.PostID
		LEFT JOIN (SELECT PostID, COUNT(*) as comments FROM comment GROUP BY PostID) comments
		ON p.PostID = comments.PostID
		ORDER BY p.CreatedDate DESC
	`)
	if err != nil {
		log.Printf("Query error: %s", err)
		return []structs.Post{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Username, &post.Likes, &post.Dislikes, &post.Comments)
		if err != nil {
			log.Printf("Scan error: %s", err)
			continue
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %s", err)
		return []structs.Post{}, err
	}
	return posts, nil
}

func GetLikedPostsByUser(username string) ([]structs.Post, error) {
	var posts []structs.Post
	rows, err := db.Query(`
        SELECT p.PostID, p.Title, p.Content, p.CreatedDate, p.username
        FROM Post p
        INNER JOIN Interaction i ON p.PostID = i.PostID
        WHERE i.username = $1 AND i.Kind = 1
    `, username)
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var post structs.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedDate, &post.Username); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func GetPost(postID int) (structs.Post, error) {
	var post structs.Post
	row := db.QueryRow(`
		SELECT 
			p.PostID, p.Title, p.Content, p.CreatedDate, p.username,
			IFNULL(likes.likes, 0) as likes, 
			IFNULL(dislikes.dislikes, 0) as dislikes
		FROM Post p
		LEFT JOIN (SELECT PostID, COUNT(*) as likes FROM Interaction WHERE Kind = 1 GROUP BY PostID) likes 
		ON p.PostID = likes.PostID
		LEFT JOIN (SELECT PostID, COUNT(*) as dislikes FROM Interaction WHERE Kind = 0 GROUP BY PostID) dislikes 
		ON p.PostID = dislikes.PostID
		WHERE p.PostID = $1`, postID)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedDate, &post.Username, &post.Likes, &post.Dislikes)
	if err != nil {
		return post, err
	}
	rows, err := db.Query("SELECT CategoryName FROM Category c INNER JOIN PostCategory pc ON c.CategoryID = pc.CategoryID WHERE pc.PostID = $1", postID)
	if err != nil {
		return post, err
	}
	defer rows.Close()
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return post, err
		}
		post.Categories = append(post.Categories, category)
	}
	return post, nil
}

func AddPostCategories(postID int, categories []string) error {
	for _, category := range categories {
		var categoryID int
		err := db.QueryRow(`
			INSERT INTO Category (CategoryName)
			VALUES ($1)
			ON CONFLICT(CategoryName) DO UPDATE SET CategoryName=excluded.CategoryName
			RETURNING CategoryID`, category).Scan(&categoryID)
		if err != nil {
			log.Printf("Error inserting/fetching category: %s\n", err.Error()) // Debug print
			return err
		}

		log.Printf("Linking category %d to post %d\n", categoryID, postID) // Debug print
		_, err = db.Exec(`
			INSERT INTO PostCategory (PostID, CategoryID)
			VALUES ($1, $2)`, postID, categoryID)
		if err != nil {
			log.Printf("Error linking category to post: %s\n", err.Error()) // Debug print
			return err
		}
	}
	return nil
}
