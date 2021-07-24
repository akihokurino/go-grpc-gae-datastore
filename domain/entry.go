package domain

import (
	"time"
)

type Entry struct {
	CustomerID CustomerID
	ProjectID  ProjectID
	CreatedAt  time.Time
}

func newEntry(customerID CustomerID, projectID ProjectID, now time.Time) *Entry {
	return &Entry{
		CustomerID: customerID,
		ProjectID:  projectID,
		CreatedAt:  now,
	}
}

func NewEntryID(customerID CustomerID, projectID ProjectID) EntryID {
	return (&Entry{CustomerID: customerID, ProjectID: projectID}).ID()
}

func (e *Entry) ID() EntryID {
	return EntryID(string(e.CustomerID) + "-" + string(e.ProjectID))
}
