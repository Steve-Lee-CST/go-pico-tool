package task_dagflow

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type Goods struct {
	ID    string
	Name  string
	Price int
}

var GoodsData = []Goods{
	{ID: "1", Name: "Apple", Price: 100},
	{ID: "2", Name: "Banana", Price: 50},
	{ID: "3", Name: "Cherry", Price: 200},
	{ID: "4", Name: "Date", Price: 150},
}

type Shop struct {
	ID       string
	Name     string
	GoodsIDs []string
}

var ShopsData = []Shop{
	{ID: "1", Name: "Shop A", GoodsIDs: []string{"1", "2"}},
	{ID: "2", Name: "Shop B", GoodsIDs: []string{"2", "3"}},
	{ID: "3", Name: "Shop C", GoodsIDs: []string{"1", "3", "4"}},
	{ID: "4", Name: "Shop D", GoodsIDs: []string{"2", "4"}},
}

type GoodsInShops struct {
	ShopToGoods map[string][]Goods
}

type GoodsInShopsCollection struct {
	shops        []Shop
	goods        []Goods
	goodsInShops GoodsInShops
}

func (c *GoodsInShopsCollection) InputTypes() []reflect.Type {
	return []reflect.Type{nil}
}

func (c *GoodsInShopsCollection) TargetTypes() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(c.goodsInShops)}
}

func (c *GoodsInShopsCollection) GetShops() []Shop {
	return c.shops
}

func (c *GoodsInShopsCollection) SetShops(shops []Shop) {
	c.shops = shops
}

func (c *GoodsInShopsCollection) GetGoods() []Goods {
	return c.goods
}

func (c *GoodsInShopsCollection) SetGoods(goods []Goods) {
	c.goods = goods
}

func (c *GoodsInShopsCollection) GetGoodsInShops() GoodsInShops {
	return c.goodsInShops
}

func (c *GoodsInShopsCollection) SetGoodsInShops(goodsInShops GoodsInShops) {
	c.goodsInShops = goodsInShops
}

type IGoods interface {
	ICollection
	SetGoods(goods []Goods)
}

var (
	_ IGoods                 = (*GoodsInShopsCollection)(nil)
	_ ICollection            = (*GoodsInShopsCollection)(nil)
	_ ITask[IGoods]          = (*GetGoodsTask[IGoods])(nil)
	_ TaskCreateFunc[IGoods] = NewGetGoodsTaskCreateFunc[IGoods]("GetGoodsTask", 50*time.Millisecond)
)

type GetGoodsTask[CT IGoods] struct {
	name    string
	timeout time.Duration
}

func NewGetGoodsTask[CT IGoods](name string, timeout time.Duration) *GetGoodsTask[CT] {
	return &GetGoodsTask[CT]{name: name, timeout: timeout}
}

func NewGetGoodsTaskCreateFunc[CT IGoods](name string, timeout time.Duration) TaskCreateFunc[CT] {
	return func() (ITask[CT], error) {
		return NewGetGoodsTask[CT](name, timeout), nil
	}
}

func (t *GetGoodsTask[CT]) Name() string {
	return t.name
}

func (t *GetGoodsTask[CT]) InputTypes() []reflect.Type {
	return []reflect.Type{nil}
}

func (t *GetGoodsTask[CT]) OutputType() reflect.Type {
	return reflect.TypeOf([]Goods{})
}

func (t *GetGoodsTask[CT]) Timeout() time.Duration {
	return t.timeout
}

func (t *GetGoodsTask[CT]) Execute(ctx context.Context, collection CT) error {
	goods := GoodsData
	if goods == nil {
		return fmt.Errorf("goods collection is nil")
	}
	time.Sleep(100 * time.Millisecond)
	collection.SetGoods(goods)
	return nil
}

type IShop interface {
	ICollection
	SetShops(shops []Shop)
}

var (
	_ IShop                 = (*GoodsInShopsCollection)(nil)
	_ ICollection           = (*GoodsInShopsCollection)(nil)
	_ ITask[IShop]          = (*GetShopsTask[IShop])(nil)
	_ TaskCreateFunc[IShop] = NewGetShopsTaskCreateFunc[IShop]("GetShopsTask", 50*time.Millisecond)
)

type GetShopsTask[CT IShop] struct {
	name    string
	timeout time.Duration
}

func NewGetShopsTask[CT IShop](name string, timeout time.Duration) *GetShopsTask[CT] {
	return &GetShopsTask[CT]{name: name, timeout: timeout}
}

func NewGetShopsTaskCreateFunc[CT IShop](name string, timeout time.Duration) TaskCreateFunc[CT] {
	return func() (ITask[CT], error) {
		return NewGetShopsTask[CT](name, timeout), nil
	}
}

func (t *GetShopsTask[CT]) Name() string {
	return t.name
}

func (t *GetShopsTask[CT]) InputTypes() []reflect.Type {
	return []reflect.Type{nil}
}

func (t *GetShopsTask[CT]) OutputType() reflect.Type {
	return reflect.TypeOf([]Shop{})
}

func (t *GetShopsTask[CT]) Timeout() time.Duration {
	return t.timeout
}

func (t *GetShopsTask[CT]) Execute(ctx context.Context, collection CT) error {
	shops := ShopsData
	if shops == nil {
		return fmt.Errorf("shops collection is nil")
	}
	time.Sleep(200 * time.Millisecond)
	collection.SetShops(shops)
	return nil
}

type IGoodsInShops interface {
	ICollection
	GetShops() []Shop
	GetGoods() []Goods
	SetGoodsInShops(goodsInShops GoodsInShops)
}

type GoodsInShopsTask[CT IGoodsInShops] struct {
	name    string
	timeout time.Duration
}

func NewGoodsInShopsTask[CT IGoodsInShops](name string, timeout time.Duration) *GoodsInShopsTask[CT] {
	return &GoodsInShopsTask[CT]{name: name, timeout: timeout}
}

func NewGoodsInShopsTaskCreateFunc[CT IGoodsInShops](name string, timeout time.Duration) TaskCreateFunc[CT] {
	return func() (ITask[CT], error) {
		return NewGoodsInShopsTask[CT](name, timeout), nil
	}
}

func (t *GoodsInShopsTask[CT]) Name() string {
	return t.name
}

func (t *GoodsInShopsTask[CT]) InputTypes() []reflect.Type {
	return AutoInputTypes[IGoodsInShops]()
}

func (t *GoodsInShopsTask[CT]) OutputType() reflect.Type {
	types := AutoOutputType[IGoodsInShops]()
	return types
}

func (t *GoodsInShopsTask[CT]) Timeout() time.Duration {
	return t.timeout
}

func (t *GoodsInShopsTask[CT]) Execute(ctx context.Context, collection CT) error {
	shops := collection.GetShops()
	goods := collection.GetGoods()
	if shops == nil || goods == nil {
		return fmt.Errorf("shops or goods collection is nil")
	}

	goodsInShops := GoodsInShops{ShopToGoods: make(map[string][]Goods)}
	for _, shop := range shops {
		for _, goodsID := range shop.GoodsIDs {
			for _, good := range goods {
				if good.ID == goodsID {
					goodsInShops.ShopToGoods[shop.ID] = append(goodsInShops.ShopToGoods[shop.ID], good)
				}
			}
		}
	}

	time.Sleep(300 * time.Millisecond)
	collection.SetGoodsInShops(goodsInShops)
	return nil
}

func TestNormal(t *testing.T) {
	factory := NewFactory[*GoodsInShopsCollection]()
	if err := factory.RegisterTask(NewGetGoodsTaskCreateFunc[*GoodsInShopsCollection](
		"GetGoodsTask", 500*time.Millisecond)); err != nil {
		t.Fatalf("failed to register GetGoodsTask: %v", err)
	}
	if err := factory.RegisterTask(NewGetShopsTaskCreateFunc[*GoodsInShopsCollection](
		"GetShopsTask", 500*time.Millisecond)); err != nil {
		t.Fatalf("failed to register GetShopsTask: %v", err)
	}
	if err := factory.RegisterTask(NewGoodsInShopsTaskCreateFunc[*GoodsInShopsCollection](
		"GoodsInShopsTask", 500*time.Millisecond)); err != nil {
		t.Fatalf("failed to register GoodsInShopsTask: %v", err)
	}

	factory.CreateGraph()

	collection := &GoodsInShopsCollection{}
	taskDagflow, err := factory.CreateTaskDagflow(collection)
	if err != nil {
		t.Fatalf("failed to create TaskDagflow: %v", err)
	}

	ctx := context.Background()
	if err := taskDagflow.Execute(ctx, 2*time.Second); err != nil {
		t.Fatalf("task dagflow execution failed: %v", err)
	}

	goodsInShops := collection.GetGoodsInShops()
	if len(goodsInShops.ShopToGoods) == 0 {
		t.Fatal("expected goods in shops, but got none")
	}
	fmt.Printf("task dagflow cost: %v\n", taskDagflow.timeCost)
}

func TestTaskTimeout(t *testing.T) {
	factory := NewFactory[*GoodsInShopsCollection]()
	if err := factory.RegisterTask(NewGetGoodsTaskCreateFunc[*GoodsInShopsCollection](
		"GetGoodsTask", 150*time.Millisecond)); err != nil {
		t.Fatalf("failed to register GetGoodsTask: %v", err)
	}
	if err := factory.RegisterTask(NewGetShopsTaskCreateFunc[*GoodsInShopsCollection](
		"GetShopsTask", 150*time.Millisecond)); err != nil {
		t.Fatalf("failed to register GetShopsTask: %v", err)
	}
	if err := factory.RegisterTask(NewGoodsInShopsTaskCreateFunc[*GoodsInShopsCollection](
		"GoodsInShopsTask", 500*time.Millisecond)); err != nil {
		t.Fatalf("failed to register GoodsInShopsTask: %v", err)
	}

	factory.CreateGraph()

	collection := &GoodsInShopsCollection{}
	taskDagflow, err := factory.CreateTaskDagflow(collection)
	if err != nil {
		t.Fatalf("failed to create TaskDagflow: %v", err)
	}

	ctx := context.Background()
	if err := taskDagflow.Execute(ctx, 2*time.Second); err != nil {
		t.Fatalf("task dagflow execution failed: %v", err)
	}

	goodsInShops := collection.GetGoodsInShops()
	if len(goodsInShops.ShopToGoods) == 0 {
		t.Fatal("expected goods in shops, but got none")
	}
	fmt.Printf("task dagflow cost: %v\n", taskDagflow.timeCost)
}

func TestFlowTimeout(t *testing.T) {
	factory := NewFactory[*GoodsInShopsCollection]()
	if err := factory.RegisterTask(NewGetGoodsTaskCreateFunc[*GoodsInShopsCollection](
		"GetGoodsTask", 500*time.Millisecond)); err != nil {
		t.Fatalf("failed to register GetGoodsTask: %v", err)
	}
	if err := factory.RegisterTask(NewGetShopsTaskCreateFunc[*GoodsInShopsCollection](
		"GetShopsTask", 500*time.Millisecond)); err != nil {
		t.Fatalf("failed to register GetShopsTask: %v", err)
	}
	if err := factory.RegisterTask(NewGoodsInShopsTaskCreateFunc[*GoodsInShopsCollection](
		"GoodsInShopsTask", 500*time.Millisecond)); err != nil {
		t.Fatalf("failed to register GoodsInShopsTask: %v", err)
	}

	factory.CreateGraph()

	collection := &GoodsInShopsCollection{}
	taskDagflow, err := factory.CreateTaskDagflow(collection)
	if err != nil {
		t.Fatalf("failed to create TaskDagflow: %v", err)
	}

	ctx := context.Background()
	if err := taskDagflow.Execute(ctx, 300*time.Millisecond); err != nil {
		t.Fatalf("task dagflow execution failed: %v", err)
	}

	goodsInShops := collection.GetGoodsInShops()
	if len(goodsInShops.ShopToGoods) == 0 {
		t.Fatal("expected goods in shops, but got none")
	}
	fmt.Printf("task dagflow cost: %v\n", taskDagflow.timeCost)
}
