package dtos

import contracts "simulator/packages/contracts/gateways"

// Gateway and GatewayInput alias the canonical contracts.
type Gateway = contracts.Gateway

type GatewayInput = contracts.GatewayInput

// GatewayIDParam binds the :id path parameter.
type GatewayIDParam = contracts.GatewayIDParam
