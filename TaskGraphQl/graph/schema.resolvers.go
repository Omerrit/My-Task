package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"gerrit-share.lan/go/errors"

	"gerrit-share.lan/go/graph/generated"
	"gerrit-share.lan/go/graph/model"
)

// по какой-то причине пытался отладить, но не появился mutation
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (bool, error) {
	newPost := &model.Post{ID: input.ID, Message: input.Message, Title: input.Title}
	r.posts = append(r.posts, newPost)
	return true, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	return r.posts, nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	// вообще говоря если есть несколько постов с одинаковым id, то могу лишь пожелать удачи
	for _, key := range r.posts {
		if key.ID == id {
			return key, nil
		}
	}
	return nil, errors.New("Not Find")
}

func (r *subscriptionResolver) Notification(ctx context.Context, id string) (<-chan *model.Post, error) {
	// не успел вызывается ли эта функция в потоке, на всякий случай напишу
	r.mutex.Lock()
	// я не успел посмотреть где можно проинициализировать эту мапу, но ее нужно сделать!
	 _, ok := r.users[id]
	 if ok {
	 	return nil, nil
	 }
	 messages := make(chan *model.Post)
	 r.users[id] = messages
	 r.mutex.Unlock()
	 // нужно где -то придумать насчет удаления данных из мапы, не успел сделать
	 return messages, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
