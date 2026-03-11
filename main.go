# Saga Checkout Microservice

## README

### Overview

This is a single-microservice implementation of the **Saga Pattern** for an e-commerce checkout workflow.

The workflow contains three steps:

1. **Payment**
2. **Inventory**
3. **Shipping**

Each step supports:

* `Do()` — performs the step
* `Compensate()` — undoes the step if a later step fails

If any step fails, all previously completed steps are compensated in **reverse order**.

### Design

* A shared `Context` stores the order state.
* Each saga step implements a common `Step` interface.
* A `Saga` coordinator executes steps sequentially.
* On failure, the coordinator runs compensation for completed steps in reverse order.

### Why this design

* Simple and easy to extend with new steps
* Keeps business flow explicit
* Demonstrates the core Saga idea without external brokers or distributed services

### Run

```bash
go run .
```

### Example behavior

* Success path: payment captured → inventory reserved → shipping created
* Failure path: if shipping fails, inventory and payment are compensated in reverse order

---

## go.mod

```go
module saga-checkout

go 1.22
```

---

## main.go

```go
package main

import (
	"errors"
	"fmt"
)

type CheckoutContext struct {
	OrderID             string
	PaymentAuthorized   bool
	InventoryReserved   bool
	ShippingScheduled   bool
	FailPayment         bool
	FailInventory       bool
	FailShipping        bool
}

type Step interface {
	Name() string
	Do(*CheckoutContext) error
	Compensate(*CheckoutContext) error
}

type PaymentStep struct{}

func (s PaymentStep) Name() string { return "Payment" }

func (s PaymentStep) Do(ctx *CheckoutContext) error {
	fmt.Println("[Payment] authorizing payment")
	if ctx.FailPayment {
		return errors.New("payment authorization failed")
	}
	ctx.PaymentAuthorized = true
	fmt.Println("[Payment] payment authorized")
	return nil
}

func (s PaymentStep) Compensate(ctx *CheckoutContext) error {
	if !ctx.PaymentAuthorized {
		return nil
	}
	fmt.Println("[Payment] refunding payment")
	ctx.PaymentAuthorized = false
	fmt.Println("[Payment] payment refunded")
	return nil
}

type InventoryStep struct{}

func (s InventoryStep) Name() string { return "Inventory" }

func (s InventoryStep) Do(ctx *CheckoutContext) error {
	fmt.Println("[Inventory] reserving inventory")
	if ctx.FailInventory {
		return errors.New("inventory reservation failed")
	}
	ctx.InventoryReserved = true
	fmt.Println("[Inventory] inventory reserved")
	return nil
}

func (s InventoryStep) Compensate(ctx *CheckoutContext) error {
	if !ctx.InventoryReserved {
		return nil
	}
	fmt.Println("[Inventory] releasing inventory")
	ctx.InventoryReserved = false
	fmt.Println("[Inventory] inventory released")
	return nil
}

type ShippingStep struct{}

func (s ShippingStep) Name() string { return "Shipping" }

func (s ShippingStep) Do(ctx *CheckoutContext) error {
	fmt.Println("[Shipping] creating shipment")
	if ctx.FailShipping {
		return errors.New("shipping creation failed")
	}
	ctx.ShippingScheduled = true
	fmt.Println("[Shipping] shipment created")
	return nil
}

func (s ShippingStep) Compensate(ctx *CheckoutContext) error {
	if !ctx.ShippingScheduled {
		return nil
	}
	fmt.Println("[Shipping] canceling shipment")
	ctx.ShippingScheduled = false
	fmt.Println("[Shipping] shipment canceled")
	return nil
}

type Saga struct {
	steps []Step
}

func NewSaga(steps ...Step) *Saga {
	return &Saga{steps: steps}
}

func (s *Saga) Execute(ctx *CheckoutContext) error {
	completed := make([]Step, 0, len(s.steps))

	for _, step := range s.steps {
		fmt.Printf("--> Executing step: %s\n", step.Name())
		if err := step.Do(ctx); err != nil {
			fmt.Printf("!! Step failed: %s: %v\n", step.Name(), err)
			s.compensate(ctx, completed)
			return err
		}
		completed = append(completed, step)
	}

	fmt.Println("✅ Checkout completed successfully")
	return nil
}

func (s *Saga) compensate(ctx *CheckoutContext, completed []Step) {
	fmt.Println("↩ Starting compensation in reverse order")
	for i := len(completed) - 1; i >= 0; i-- {
		step := completed[i]
		fmt.Printf("<-- Compensating step: %s\n", step.Name())
		if err := step.Compensate(ctx); err != nil {
			fmt.Printf("!! Compensation failed for %s: %v\n", step.Name(), err)
		}
	}
}

func main() {
	saga := NewSaga(
		PaymentStep{},
		InventoryStep{},
		ShippingStep{},
	)

	fmt.Println("=== Successful checkout ===")
	successCtx := &CheckoutContext{OrderID: "order-1001"}
	_ = saga.Execute(successCtx)

	fmt.Println()
	fmt.Println("=== Failing checkout (shipping fails) ===")
	failureCtx := &CheckoutContext{
		OrderID:      "order-1002",
		FailShipping: true,
	}
	_ = saga.Execute(failureCtx)
}
```
 
