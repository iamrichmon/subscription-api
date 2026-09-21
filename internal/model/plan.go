package model

type Plan struct {
	ID       SubscriptionStatus `json:"id"`
	Name     string             `json:"name"`
	Price    int                `json:"price"`
	Features []string           `json:"features"`
}

var Plans = []Plan{
	{ID: FreeSubscription,
		Name:  "Free",
		Price: 0,
		Features: []string{
			"Access to basic features",
		},
	},
	{ID: ProSubscription,
		Name:  "Pro",
		Price: 10,
		Features: []string{
			"Access to all features",
		},
	},
	{ID: PremiumSubscription,
		Name:  "Premium",
		Price: 20,
		Features: []string{
			"Access to all features",
			"Priority support",
		}},
}
