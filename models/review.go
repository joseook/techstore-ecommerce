package models

type Review struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	User      string  `json:"user"`
	Rating    float64 `json:"rating"`
	Comment   string  `json:"comment"`
	Date      string  `json:"date"`
}

var Reviews = map[string][]Review{
	"1": {
		{
			ID:        "1",
			ProductID: "1",
			User:      "John Doe",
			Rating:    5.0,
			Comment:   "Amazing phone! The camera quality is outstanding and the battery life is impressive.",
			Date:      "2024-03-15",
		},
		{
			ID:        "2",
			ProductID: "1",
			User:      "Jane Smith",
			Rating:    4.5,
			Comment:   "Great phone overall, but a bit expensive. The titanium frame feels premium though.",
			Date:      "2024-03-10",
		},
	},
	"2": {
		{
			ID:        "3",
			ProductID: "2",
			User:      "Mike Johnson",
			Rating:    5.0,
			Comment:   "Best laptop I've ever owned. The M3 Pro chip is incredibly fast and the display is stunning.",
			Date:      "2024-03-12",
		},
		{
			ID:        "4",
			ProductID: "2",
			User:      "Sarah Williams",
			Rating:    4.8,
			Comment:   "Perfect for my work as a designer. The battery life is amazing and it never gets hot.",
			Date:      "2024-03-08",
		},
	},
	// Add more reviews for other products as needed
} 