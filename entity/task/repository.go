package task

import "context"

type Repository interface {
	Get(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context) ([]*Task, error)

	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}
