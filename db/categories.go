package db

// CreateCategory creates a new category in the database.
func CreateCategory(name string) error {
	_, err := DB.Exec("INSERT INTO categories (category) VALUES (?)", name)
	return err
}

// GetCategories retrieves all categories from the database.
func GetCategories() ([]Category, error) {
	rows, err := DB.Query("SELECT category_id, category FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetCategoryByID returns a category by its ID.
func GetCategoryByID(id int) (*Category, error) {
	var category Category
	err := DB.QueryRow("SELECT category_id, category FROM categories WHERE category_id = ?", id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
