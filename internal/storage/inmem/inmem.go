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
