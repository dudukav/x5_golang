package service

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hw5/internal/model"
)

var errRepository = errors.New("repository error")

type repositoryStub struct {
	order          *model.Order
	orders         []*model.Order
	createdOrder   *model.Order
	updatedOrder   *model.Order
	deletedID      int
	requestedID    int
	customer       string
	createErr      error
	getByIDErr     error
	listErr        error
	updateErr      error
	deleteErr      error
	getItemsErr    error
	getCustomerErr error
}

func (r *repositoryStub) Create(order *model.Order) (*model.Order, error) {
	r.createdOrder = order
	return order, r.createErr
}

func (r *repositoryStub) GetByID(id int) (*model.Order, error) {
	r.requestedID = id
	if r.order == nil {
		return &model.Order{ID: id}, r.getByIDErr
	}
	return cloneOrder(r.order), r.getByIDErr
}

func (r *repositoryStub) List() ([]*model.Order, error) {
	return r.orders, r.listErr
}

func (r *repositoryStub) Update(order *model.Order) (*model.Order, error) {
	r.updatedOrder = order
	return order, r.updateErr
}

func (r *repositoryStub) Delete(id int) error {
	r.deletedID = id
	return r.deleteErr
}

func (r *repositoryStub) GetItems(id int) ([]model.Item, error) {
	r.requestedID = id
	return []model.Item{{Name: "coffee", Quantity: 1, Price: 250}}, r.getItemsErr
}

func (r *repositoryStub) GetByCustomer(customer string) ([]*model.Order, error) {
	r.customer = customer
	return []*model.Order{{ID: 1, Customer: customer}}, r.getCustomerErr
}

func TestCreateAssignsFieldsAndCalculatesTotal(t *testing.T) {
	repository := &repositoryStub{}
	orders := New(repository)
	now := time.Date(2026, time.November, 27, 12, 0, 0, 0, time.FixedZone("test", 3*60*60))
	orders.now = func() time.Time { return now }
	input := validOrder()

	created, err := orders.Create(input)

	require.NoError(t, err)
	assert.Equal(t, 1, created.ID)
	assert.Equal(t, model.StatusNew, created.Status)
	assert.Equal(t, int64(800), created.TotalPrice)
	assert.Equal(t, now.UTC(), created.CreatedAt)
	require.Len(t, created.History, 1)
	assert.Equal(t, model.StatusNew, created.History[0].Status)
	assert.Equal(t, now.UTC(), created.History[0].ChangedAt)
	assert.Zero(t, input.ID, "Create must not mutate its input")
}

func TestCreateUsesSequentialIDs(t *testing.T) {
	orders := New(&repositoryStub{})
	first, err := orders.Create(validOrder())
	require.NoError(t, err)
	second, err := orders.Create(validOrder())
	require.NoError(t, err)
	assert.Equal(t, 1, first.ID)
	assert.Equal(t, 2, second.ID)
}

func TestCreateValidation(t *testing.T) {
	tests := []struct {
		name  string
		order *model.Order
	}{
		{name: "nil order"},
		{name: "empty items", order: &model.Order{}},
		{name: "empty item name", order: orderWithItem(model.Item{Quantity: 1, Price: 10})},
		{name: "zero quantity", order: orderWithItem(model.Item{Name: "item", Price: 10})},
		{name: "negative quantity", order: orderWithItem(model.Item{Name: "item", Quantity: -1, Price: 10})},
		{name: "zero price", order: orderWithItem(model.Item{Name: "item", Quantity: 1})},
		{name: "negative price", order: orderWithItem(model.Item{Name: "item", Quantity: 1, Price: -10})},
		{name: "total overflow", order: orderWithItem(model.Item{Name: "item", Quantity: 2, Price: math.MaxInt64})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &repositoryStub{}
			_, err := New(repository).Create(test.order)
			require.ErrorIs(t, err, ErrInvalidOrder)
			assert.Nil(t, repository.createdOrder)
		})
	}
}

func TestStatusTransitions(t *testing.T) {
	tests := []struct {
		name string
		from model.Status
		to   model.Status
	}{
		{name: "new to processing", from: model.StatusNew, to: model.StatusProcessing},
		{name: "new to cancelled", from: model.StatusNew, to: model.StatusCancelled},
		{name: "processing to completed", from: model.StatusProcessing, to: model.StatusCompleted},
		{name: "processing to cancelled", from: model.StatusProcessing, to: model.StatusCancelled},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			previous := time.Date(2026, 11, 27, 10, 0, 0, 0, time.UTC)
			now := previous.Add(time.Hour)
			repository := &repositoryStub{order: &model.Order{
				ID:      1,
				Status:  test.from,
				History: []model.StatusChange{{Status: test.from, ChangedAt: previous}},
			}}
			orders := New(repository)
			orders.now = func() time.Time { return now }

			updated, err := orders.UpdateStatus(1, test.to)

			require.NoError(t, err)
			assert.Equal(t, test.to, updated.Status)
			require.Len(t, updated.History, 2)
			assert.Equal(t, test.to, updated.History[1].Status)
			assert.True(t, updated.History[0].ChangedAt.Before(updated.History[1].ChangedAt))
		})
	}
}

func TestInvalidStatusTransitions(t *testing.T) {
	tests := []struct {
		name string
		from model.Status
		to   model.Status
	}{
		{name: "new to completed", from: model.StatusNew, to: model.StatusCompleted},
		{name: "completed is final", from: model.StatusCompleted, to: model.StatusCancelled},
		{name: "cancelled is final", from: model.StatusCancelled, to: model.StatusProcessing},
		{name: "same status", from: model.StatusProcessing, to: model.StatusProcessing},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &repositoryStub{order: &model.Order{ID: 1, Status: test.from}}
			_, err := New(repository).UpdateStatus(1, test.to)
			require.ErrorIs(t, err, ErrInvalidStatusTransition)
			assert.Nil(t, repository.updatedOrder)
		})
	}
}

func TestReadListAndHistory(t *testing.T) {
	history := []model.StatusChange{{Status: model.StatusNew, ChangedAt: time.Now()}}
	repository := &repositoryStub{
		order:  &model.Order{ID: 4, History: history},
		orders: []*model.Order{{ID: 1}, {ID: 2}},
	}
	orders := New(repository)

	got, err := orders.GetByID(4)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
	list, err := orders.List()
	require.NoError(t, err)
	assert.Len(t, list, 2)
	gotHistory, err := orders.History(4)
	require.NoError(t, err)
	assert.Equal(t, history, gotHistory)
}

func TestServiceWrapsRepositoryErrors(t *testing.T) {
	repository := &repositoryStub{
		order:          validOrder(),
		createErr:      errRepository,
		getByIDErr:     errRepository,
		listErr:        errRepository,
		updateErr:      errRepository,
		deleteErr:      errRepository,
		getItemsErr:    errRepository,
		getCustomerErr: errRepository,
	}
	orders := New(repository)
	update := validOrder()
	update.ID = 1
	update.Status = model.StatusNew

	operations := []struct {
		name string
		run  func() error
	}{
		{name: "create", run: func() error { _, err := orders.Create(validOrder()); return err }},
		{name: "get by id", run: func() error { _, err := orders.GetByID(1); return err }},
		{name: "list", run: func() error { _, err := orders.List(); return err }},
		{name: "update", run: func() error { _, err := orders.Update(update); return err }},
		{name: "update status", run: func() error { _, err := orders.UpdateStatus(1, model.StatusProcessing); return err }},
		{name: "delete", run: func() error { return orders.Delete(1) }},
		{name: "get items", run: func() error { _, err := orders.GetItems(1); return err }},
		{name: "get by customer", run: func() error { _, err := orders.GetByCustomer("Alice"); return err }},
	}

	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			require.ErrorIs(t, operation.run(), errRepository)
		})
	}
}

func TestInvalidIdentifiersAndInputs(t *testing.T) {
	orders := New(&repositoryStub{})
	_, err := orders.GetByID(0)
	require.ErrorIs(t, err, ErrInvalidOrder)
	_, err = orders.UpdateStatus(1, model.Status("unknown"))
	require.ErrorIs(t, err, ErrInvalidOrder)
	_, err = orders.GetByCustomer(" ")
	require.ErrorIs(t, err, ErrInvalidOrder)
}

func validOrder() *model.Order {
	return &model.Order{Items: []model.Item{
		{Name: "coffee", Quantity: 2, Price: 250},
		{Name: "cake", Quantity: 1, Price: 300},
	}}
}

func orderWithItem(item model.Item) *model.Order {
	return &model.Order{Items: []model.Item{item}}
}
