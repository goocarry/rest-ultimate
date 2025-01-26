package storage

type Dating struct {
	ID        int64
	UserID    int64
	PartnerID int64
	CreatedAt string
	Status    string
}

type DatingRepository interface {
	CreateDating(dating Dating) (int64, error)
	GetDatingByID(id int64) (*Dating, error)
	GetAllDatingByUserID(userID int64) ([]Dating, error)
	UpdateDating(dating Dating) error
	DeleteDating(id int64) error
}
