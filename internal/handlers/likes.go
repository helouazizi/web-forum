package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/database"
)

// LikePost handles liking/disliking a post
func LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	postID, err1 := strconv.Atoi(r.FormValue("post_id"))
	userID, err2 := strconv.Atoi(r.FormValue("user_id"))
	reaction, err3 := strconv.Atoi(r.FormValue("reaction")) // 1 for like, -1 for dislike

	if err1 != nil || err2 != nil || err3 != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	// Check if user already reacted
	var existingReaction int
	err = database.Database.QueryRow("SELECT reaction FROM likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&existingReaction)

	if err == sql.ErrNoRows {
		// Insert new like/dislike
		_, err = database.Database.Exec("INSERT INTO likes (post_id, user_id, reaction) VALUES (?, ?, ?)", postID, userID, reaction)
	} else if err == nil && existingReaction != reaction {
		// Update reaction
		_, err = database.Database.Exec("UPDATE likes SET reaction = ? WHERE post_id = ? AND user_id = ?", reaction, postID, userID)
	} else {
		// Remove reaction if user clicks again
		_, err = database.Database.Exec("DELETE FROM likes WHERE post_id = ? AND user_id = ?", postID, userID)
	}

	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Count likes & dislikes
	var likes int
	_ = database.Database.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND reaction = 1", postID).Scan(&likes)
	//_ = database.Database.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND reaction = -1", postID).Scan(&dislikes)

	// Update `total_likes` and `total_dislikes` in `posts`
	_, _ = database.Database.Exec("UPDATE posts SET total_likes = ?,  WHERE id = ?", likes, postID)

	// Return updated counts
	/*w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"likes":    likes,
		"dislikes": dislikes,
	})*/
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Like_Comment(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	// lets extract the post id
	post_id := r.URL.Query().Get("id")

	result, err := database.Database.Exec("UPDATE comments SET total_likes = total_likes + 1 WHERE post_id = $1", post_id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("could not retrieve rows affected: %v", err)
		return
	}
	if rowsAffected == 0 {
		fmt.Printf("no post found with id %v", post_id)
		return
	}

	fmt.Printf("Successfully updated totallikes for post ID %v\n", post_id)
	http.Redirect(w, r, "/", http.StatusFound)
}
