# Saga-Checkout-Microservice
Implementation of the Saga Pattern in a single Go microservice for an e-commerce checkout workflow. The workflow includes Payment, Inventory, and Shipping steps, each supporting do and compensate actions. If any step fails, previously completed steps are compensated in reverse order to maintain consistency.
