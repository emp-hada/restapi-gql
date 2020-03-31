package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/emp/restapi-gql/data"
	"github.com/emp/restapi-gql/graph/generated"
	model1 "github.com/emp/restapi-gql/graph/model"
	"github.com/emp/restapi-gql/model"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model1.NewTodo) (*model1.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateBook(ctx context.Context, input model1.NewBook, author model1.NewAuthor) (*model.Book, error) {
	fmt.Println(input, author)
	book := new(model.Book)

	book.ID = strconv.Itoa(rand.Intn(10000000)) // Mock ID - not safe
	book.Title = input.Title
	book.Isbn = input.Isbn
	book.Author = &model.Author{Firstname: author.Firstname, Lastname: author.Lastname}

	data.Books = append(data.Books, book)

	return book, nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model1.Todo, error) {
	return nil, nil
}

func (r *queryResolver) Book(ctx context.Context, id string) (*model.Book, error) {
	book := new(model.Book)

	for _, b := range data.Books {
		if b.ID == id {
			book = b
			break
		}
	}

	if *book == *new(model.Book) {
		return nil, errors.New("Not found")
	}

	return book, nil
}

func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	return data.Books, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
