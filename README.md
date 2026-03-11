### Run 
go run
 
Saga Checkout Microservice
README
Overview

This is a single-microservice implementation of the Saga Pattern for an e-commerce checkout workflow.

The workflow contains three steps:

Payment

Inventory

Shipping

Each step supports:

Do() — performs the step

Compensate() — undoes the step if a later step fails

If any step fails, all previously completed steps are compensated in reverse order.

Design

A shared Context stores the order state.

Each saga step implements a common Step interface.

A Saga coordinator executes steps sequentially.

On failure, the coordinator runs compensation for completed steps in reverse order.

Why this design

Simple and easy to extend with new steps

Keeps business flow explicit
Demonstrates the core Saga idea without external brokers or distributed services
Run
go run .
Example behavior
Success path: payment captured → inventory reserved → shipping created

Failure path: if shipping fails, inventory and payment are compensated in reverse order
