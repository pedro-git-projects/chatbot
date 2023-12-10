package users

type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleCollaborator UserRole = "collaborator"
	RoleUser         UserRole = "user"
)
