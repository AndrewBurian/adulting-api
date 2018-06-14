package data

type Activity struct {
	Name        string
	Description string
	PointValue  int
}

type ActivityDAL interface {
	GetActivites() ([]*Activity, error)
}

type mockActivityDal struct {
	all []*Activity
}

func NewMockActivityDal() ActivityDAL {
	return &mockActivityDal{
		all: []*Activity{
			{
				Name:        "Work out",
				Description: "Go get some exercise",
				PointValue:  10,
			},
			{
				Name:        "Apply for jobs",
				Description: "Fill out your resume in dropdown box form",
				PointValue:  12,
			},
			{
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
