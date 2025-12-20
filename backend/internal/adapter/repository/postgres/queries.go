package postgres

const (
	// queryCreateDrawing inserts a new drawing into the database
	queryCreateDrawing = `
		INSERT INTO drawings (id, slug, name, data, share_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// queryFindDrawingByID retrieves a drawing by its ID
	queryFindDrawingByID = `
		SELECT id, slug, name, data, share_token, created_at, updated_at
		FROM drawings
		WHERE id = $1
	`

	// queryFindDrawingBySlug retrieves a drawing by its slug
	queryFindDrawingBySlug = `
		SELECT id, slug, name, data, share_token, created_at, updated_at
		FROM drawings
		WHERE slug = $1
	`

	// queryFindDrawingByShareToken retrieves a drawing by its share token
	queryFindDrawingByShareToken = `
		SELECT id, slug, name, data, share_token, created_at, updated_at
		FROM drawings
		WHERE share_token = $1
	`

	// queryFindAllDrawings retrieves all drawings with pagination
	queryFindAllDrawings = `
		SELECT id, slug, name, data, share_token, created_at, updated_at
		FROM drawings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	// queryUpdateDrawing updates an existing drawing
	queryUpdateDrawing = `
		UPDATE drawings
		SET name = $1, data = $2, share_token = $3, updated_at = $4
		WHERE id = $5
	`

	// queryUpdateShareToken updates the share token for a drawing
	queryUpdateShareToken = `
		UPDATE drawings
		SET share_token = $1, updated_at = $2
		WHERE id = $3
	`

	// queryDeleteDrawing deletes a drawing by ID
	queryDeleteDrawing = `
		DELETE FROM drawings
		WHERE id = $1
	`

	// queryCountDrawings returns the total number of drawings
	queryCountDrawings = `
		SELECT COUNT(*)
		FROM drawings
	`
)
