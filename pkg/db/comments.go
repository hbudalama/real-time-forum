package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"rtf/pkg/structs"
)

func GetCommentDetails(commentId int) (structs.Comment, structs.User, []structs.CommentInteraction) {
	var (
		thisComment              structs.Comment
		thisUser                 structs.User
		theseCommentInteractions []structs.CommentInteraction
	)
	err := db.QueryRow("SELECT * FROM comment WHERE id = $1", commentId).Scan(&thisComment.ID, &thisComment.Content, &thisComment.PostID)
	if err != nil {
		return structs.Comment{}, structs.User{}, []structs.CommentInteraction{}
	}
	return thisComment, thisUser, theseCommentInteractions
}

func GetComments(postID int) ([]structs.Comment, error) {
	var comments []structs.Comment
	rows, err := db.Query("SELECT CommentID, Content, CreatedDate, PostID, username FROM Comment WHERE PostID = $1", postID)
	if err != nil {
		return comments, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment structs.Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedDate, &comment.PostID, &comment.Username); err != nil {
			return comments, err
		}
		likes, dislikes, err := GetCommentInteractions(comment.ID)
		if err != nil {
			return comments, err
		}
		comment.Likes = likes
		comment.Dislikes = dislikes
		comments = append(comments, comment)
	}
	return comments, nil
}

func AddComment(PostID int, username string, Content string) error {
	fmt.Println("i am here 1")
	// if !postExists(PostID) {
	// 	return errors.New("post does not exist")
	// }   this  vis not working
	fmt.Println("i am here nowww ")
	_, err := db.Exec("INSERT INTO Comment (PostID, username, Content) VALUES ($1, $2, $3)", PostID, username, Content)
	fmt.Println("i am here at end")
	if err != nil {
		return err
	}
	return nil
}

func AddCommentInteraction(commentID int, username string, kind int) error {
	if kind != 0 && kind != 1 {
		return errors.New("invalid interaction kind")
	}
	var existingKind int
	err := db.QueryRow("SELECT Kind FROM CommentInteractions WHERE CommentID = ? AND Username = ?", commentID, username).Scan(&existingKind)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing interaction, insert a new one
			_, err = db.Exec(
				"INSERT INTO CommentInteractions (CommentID, Username, Kind) VALUES (?, ?, ?)",
				commentID, username, kind,
			)
			if err != nil {
				log.Printf("InsertCommentInteraction error: %s", err)
				return err
			}
		} else {
			// Some other error occurred
			log.Printf("Query error: %s", err)
			return err
		}
	} else {
		// Existing interaction found, update it if necessary
		if existingKind != kind {
			_, err = db.Exec(
				"UPDATE CommentInteractions SET Kind = ? WHERE CommentID = ? AND Username = ?",
				kind, commentID, username,
			)
			if err != nil {
				log.Printf("UpdateCommentInteraction error: %s", err)
				return err
			}
		} else {
			// If the interaction is the same, delete it
			_, err = db.Exec(
				"DELETE FROM CommentInteractions WHERE CommentID = ? AND Username = ?",
				commentID, username,
			)
			if err != nil {
				log.Printf("DeleteCommentInteraction error: %s", err)
				return err
			}
		}
	}
	return nil
}

func GetCommentInteractions(commentID int) (likes int, dislikes int, err error) {
	err = db.QueryRow("SELECT COUNT(*) FROM CommentInteractions WHERE CommentID = $1 AND Kind = 1", commentID).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}
	err = db.QueryRow("SELECT COUNT(*) FROM CommentInteractions WHERE CommentID = $1 AND Kind = 0", commentID).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}