package data

type Activity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PointValue  int    `json:"points"`
}

type ActivityDAL interface {
	GetActivites() ([]*Activity, error)
	GetActivity(*Activity) error
}

type mockActivityDal struct {
	all []*Activity
}

func NewMockActivityDal() ActivityDAL {
	return &mockActivityDal{
		all: []*Activity{
			{
				ID:          "1",
				Name:        "Work out",
				Description: "Go get some exercise",
				PointValue:  10,
			},
			{
				ID:          "2",
				Name:        "Apply for jobs",
				Description: "Fill out your resume in dropdown box form",
				PointValue:  12,
			},
			{
				ID:          "3",
				Name:        "Eat a whole pizza",
				Description: "Do it",
				PointValue:  -50,
			},
		},
	}
}

func (mock *mockActivityDal) GetActivites() ([]*Activity, error) {
	return mock.all, nil
}

func (mock *mockActivityDal) GetActivity(a *Activity) error {
	for _, act := range mock.all {
		if act.ID == a.ID {
			a.Description = act.Description
			a.Name = act.Name
			a.PointValue = act.PointValue
			return nil
		}
	}

	return ErrNotFound
}
