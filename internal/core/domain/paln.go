package domain

type PlanType string

const (
	PlanTypeGeo = "geo"
	PlanTypeAi  = "ai"
)

type Plan struct {
	Base
	UserID      string `gorm:"unique"`
	Name        string `gorm:"size:255;not null;unique" json:"name"`
	CreditsUsed int    `json:"creditsUsed"` // Requests per hour
	CreditLimit int    `json:"creditLimit"` // Maximum credits per user
}

func NewPlan(userID, name, planType string) Plan {
	plan := Plan{
		UserID: userID,
		Name:   name,
	}

	return plan
}
