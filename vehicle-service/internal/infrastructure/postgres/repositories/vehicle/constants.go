package vehicle

const (
	// createVehicleQuery inserts a new vehicle and returns the generated fields.
	createVehicleQuery = `INSERT INTO vehicles (plate) VALUES ($1) RETURNING id, created_at, updated_at`

	// softDeleteVehicleQuery marks a vehicle as deleted by setting deleted_at and updated_at.
	softDeleteVehicleQuery = `UPDATE vehicles SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	// getVehicleByIDQuery selects an active vehicle by its ID.
	getVehicleByIDQuery = `SELECT id, plate, created_at, updated_at FROM vehicles WHERE id = $1 AND deleted_at IS NULL`

	// findVehicleByPlateQuery selects an active vehicle by its plate.
	findVehicleByPlateQuery = `SELECT id, plate, created_at, updated_at FROM vehicles WHERE plate = $1 AND deleted_at IS NULL`

	// listVehiclesQuery selects active vehicles with pagination ordered by id.
	listVehiclesQuery = `SELECT id, plate, created_at, updated_at FROM vehicles WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2`

	// countVehiclesQuery counts the total number of active vehicles.
	countVehiclesQuery = `SELECT COUNT(*) FROM vehicles WHERE deleted_at IS NULL`

	// uniqueViolationCode is the PostgreSQL error code for unique constraint violations.
	uniqueViolationCode = "23505"
)
