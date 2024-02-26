package domain

import (
	"context"
	"fmt"
	mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

func (i *Impl) CreateTodo(ctx context.Context, title string, description string) (*entities.Todo, error) {
	id := i.db.NewId()
	todo := entities.Todo{
		Id:          id,
		Title:       title,
		Description: description,
	}
	_, err := i.db.InsertRecord(ctx, TodoTable, todo)
	//fmt.Println("user", nUser)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (i *Impl) UpdateTodo(ctx context.Context, id string, title string, description string) (*entities.Todo, error) {
	fmt.Println(id)
	err := i.db.UpdateMany(
		ctx,
		TodoTable,
		mongo_db.Filter{"id": id},
		mongo_db.Filter{
			"title":       title,
			"description": description,
		},
	)

	var todo entities.Todo
	err = i.db.FindOne(ctx, TodoTable, &todo, mongo_db.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	fmt.Println("todo", todo)
	return &todo, nil
}

func (i *Impl) GetTodo(ctx context.Context, id string) (*entities.Todo, error) {
	var todo entities.Todo
	err := i.db.FindOne(ctx, TodoTable, &todo, mongo_db.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	fmt.Println("todo", todo)
	return &todo, nil
}

func (i *Impl) GetTodos(ctx context.Context) ([]*entities.Todo, error) {
	todos := make([]*entities.Todo, 0)
	cursor, err := i.db.Find(ctx, TodoTable)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cer := cursor.Close(ctx); cer != nil {
			fmt.Println(cer)
		}
	}()
	for cursor.Next(ctx) {
		var todo entities.Todo
		if err := cursor.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (i *Impl) DeleteTodo(ctx context.Context, id string) (string, error) {
	var todo entities.Todo
	err := i.db.FindOne(ctx, TodoTable, &todo, mongo_db.Filter{"id": id})
	if err != nil {
		return "", err
	}

	err = i.db.DeleteRecord(ctx, TodoTable, mongo_db.Filter{"id": id})
	if err != nil {
		return "", err
	}
	return "todo deleted successfully", nil
}
