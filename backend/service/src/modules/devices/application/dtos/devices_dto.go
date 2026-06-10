package dtos

import contracts "simulator/packages/contracts/devices"

// Device and DeviceInput are aliases to the canonical contracts; the module never
// redefines the wire shapes.
type Device = contracts.Device

type DeviceInput = contracts.DeviceInput

// DeviceIDParam binds the :id path parameter.
type DeviceIDParam = contracts.DeviceIDParam
