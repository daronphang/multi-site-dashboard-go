package domain

// Based on Clean Architecture, the Service/UseCase depends on a Repository interface,
// and the Repository will be implemented in the Data Layer.
// The Repository interface should contain all the methods required for the Use Case.
// However, since sqlc generates the Repository implementation together with its models,
// to use that as the interface to achieve dependency inversion.
// Ensure emit_interface is set to true in sqlc.yaml.

type Repository interface {}