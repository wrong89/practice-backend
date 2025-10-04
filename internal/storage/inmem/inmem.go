package inmem

type Storage struct {
	EntryList
	UserList
}

func NewStorage() *Storage {
	return &Storage{
		EntryList: NewEntryList(),
		UserList:  NewUserList(),
	}
}

// func (s *Storage) GetEntryByID(ctx context.Context, id int) (entry.Entry, error) {}

// func (s *Storage) CreateEntry(
// 	ctx context.Context,
// 	course string,
// 	date time.Time,
// 	UserID int,
// 	PaymentID int,
// ) (entry.Entry, error) {
// 	newEntry := entry.NewEntry(course, date, UserID, PaymentID)
// }
