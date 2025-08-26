# L8Parser

A powerful parsing service for processing l8collector payloads in the Layer8 ecosystem. L8Parser provides a rule-based parsing framework that processes collected data from network devices and transforms it into structured objects.

## Overview

L8Parser is part of the Layer8 network management ecosystem and serves as the data transformation layer between data collection and inventory management. It processes raw data collected by l8collector and applies configurable parsing rules to extract meaningful information that can be stored in the inventory.

## Features

- **Rule-Based Parsing**: Flexible parsing engine with pluggable rules
- **SNMP Support**: Built-in support for SNMP data parsing with MIB-2 compatibility
- **Protocol Agnostic**: Extensible architecture supporting multiple collection protocols
- **Real-time Processing**: Event-driven processing of collector jobs
- **Service Integration**: Seamless integration with Layer8 service architecture

## Architecture

### Core Components

- **Parser Service** (`go/parser/service/Parser.go`): Main parsing engine that processes jobs
- **Parsing Service** (`go/parser/service/ParsingService.go`): Service wrapper implementing Layer8 service interface
- **Parsing Center** (`go/parser/service/ParsingCenter.go`): Job completion handler and inventory integration
- **Rule Engine** (`go/parser/rules/`): Collection of parsing rules for data transformation

### Parsing Rules

The parser supports multiple rule types:

1. **Contains** (`Contains.go`): Searches for specific text patterns in data
2. **Set** (`Set.go`): Directly sets values to object properties
3. **ToTable** (`ToTable.go`): Converts data into tabular format
4. **TableToMap** (`TableToMap.go`): Transforms tables into key-value maps

## Installation

### Prerequisites

- Go 1.23.8 or higher
- Access to Layer8 ecosystem modules

### Dependencies

The project depends on several Layer8 modules:
- `l8collector` - Data collection framework
- `l8inventory` - Inventory management
- `l8pollaris` - Polling configuration
- `l8types` - Common types and interfaces
- `l8services` - Service framework

### Build

```bash
# Clone the repository
git clone https://github.com/saichler/l8parser

# Navigate to the Go module
cd l8parser/go

# Run the build script
./test.sh
```

The `test.sh` script will:
- Initialize Go modules
- Fetch dependencies
- Run security checks
- Execute unit tests with coverage
- Generate coverage report

## Usage

### Basic Integration

```go
import (
    "github.com/saichler/l8parser/go/parser/service"
    "github.com/saichler/l8types/go/ifs"
)

// Create and activate parsing service
parsingService := &service.ParsingService{}
err := parsingService.Activate(
    "parsing-service",
    serviceArea,
    resources,
    vnic,
    &NetworkDevice{}, // Target object type
    "Id",            // Primary key field
)
```

### Configuration Example

The parser uses Pollaris configurations to define parsing rules:

```go
// SNMP polling configuration with parsing rules
poll := &types.Poll{
    Name:      "systemMib",
    What:      ".1.3.6.1.2.1.1", // SNMP OID
    Operation: types.Operation_OMap,
    Attributes: []*types.Attribute{
        {
            PropertyId: "networkdevice.equipmentinfo.vendor",
            Rules: []*types.Rule{
                {
                    Name: "Contains",
                    Params: map[string]*types.Parameter{
                        "what":   {Value: "cisco"},
                        "from":   {Value: ".1.3.6.1.2.1.1.1.0"},
                        "output": {Value: "Cisco"},
                    },
                },
            },
        },
    },
}
```

### Processing Flow

1. **Job Reception**: Parser receives completed collection jobs from l8collector
2. **Rule Application**: Applies configured parsing rules based on Pollaris configuration
3. **Data Transformation**: Transforms raw data into structured objects
4. **Inventory Update**: Sends parsed objects to inventory service via PATCH operations

## SNMP Support

L8Parser includes comprehensive support for SNMP data parsing with vendor-specific configurations:

### Multi-Vendor Support

The parser now supports vendor-specific device configurations with automatic OID-based detection:

- **Cisco**: Switches and Routers with enterprise OID `.1.3.6.1.4.1.9.`
- **Juniper**: Routers with enterprise OID `.1.3.6.1.4.1.2636.`
- **Palo Alto**: Firewalls with enterprise OID `.1.3.6.1.4.1.25461.`
- **Fortinet**: Firewalls with enterprise OID `.1.3.6.1.4.1.12356.`
- **Arista**: Switches with enterprise OID `.1.3.6.1.4.1.30065.`
- **Nokia**: Routers with enterprise OID `.1.3.6.1.4.1.6527.`
- **Huawei**: Routers with enterprise OID `.1.3.6.1.4.1.2011.`
- **Dell**: Servers with enterprise OID `.1.3.6.1.4.1.674.`
- **HPE**: Servers with enterprise OID `.1.3.6.1.4.1.232.`

### Device Detection and Polling

```go
// Automatic device detection by sysOID
polaris := GetPollarisByOid(".1.3.6.1.4.1.9.1.122.0") // Returns Cisco switch configuration

// Get all available vendor models
models := GetAllPolarisModels() // Returns slice of all vendor-specific Pollaris models
```

### Vendor-Specific Features

Each vendor configuration includes device-specific polling for:
- **System Information**: Vendor, version, serial numbers
- **Interface Monitoring**: Status, speed, MTU, names
- **Hardware Monitoring**: Modules, power supplies, fans, temperature
- **Performance Metrics**: CPU, memory, utilization
- **Device-Specific Attributes**: Routing tables, firewall sessions, etc.

### File Organization

SNMP configurations are organized by vendor in separate files:
- `cisco.go` - Cisco switches and routers
- `juniper.go` - Juniper routers  
- `paloalto.go` - Palo Alto firewalls
- `fortinet.go` - Fortinet firewalls
- `arista.go` - Arista switches
- `nokia.go` - Nokia routers
- `huawei.go` - Huawei routers
- `dell.go` - Dell servers
- `hpe.go` - HPE servers
- `SNMP.go` - Common utilities and generic polling

### Group Classification

Each vendor model uses specific group classifications:
- Generic SNMP: `["boot"]` (common.BOOT_GROUP)
- Vendor models: `[vendor_name, vendor-device_type]` (e.g., `["cisco", "cisco-switch"]`)

### Example Cisco Switch Configuration

```go
func CreateCiscoSwitchBootPolls() *types.Pollaris {
    polaris := &types.Pollaris{
        Name:   "cisco-switch",
        Groups: []string{"cisco", "cisco-switch"},
        Polling: map[string]*types.Poll{
            "ciscoSystem": {
                What:      ".1.3.6.1.2.1.1",
                Operation: types.Operation_OMap,
                Attributes: []*types.Attribute{
                    createCiscoVendor(),
                    createSysName(),
                    createCiscoVersion(),
                },
            },
        },
    }
    return polaris
}
```

## Testing

### Unit Tests

Run the comprehensive test suite:

```bash
cd go
./test.sh
```

### Test Coverage

The test script generates detailed coverage reports. Coverage includes:
- Parser engine functionality
- Rule execution
- Service integration
- Error handling

### Integration Testing

The test suite includes integration tests that:
- Simulate SNMP device responses
- Test end-to-end parsing workflows
- Validate inventory integration

## Configuration

### Service Configuration

The parsing service requires:
- **Target Object Type**: The Go struct type for parsed data
- **Primary Key Field**: Field name for object identification
- **Service Dependencies**: Connections to inventory and collector services

### Rule Configuration

Parsing rules are configured through Pollaris definitions:
- **Attributes**: Define target object properties
- **Rules**: Specify transformation logic
- **Parameters**: Provide rule-specific configuration

## Error Handling

L8Parser implements comprehensive error handling:
- **Validation Errors**: Missing parameters and invalid configurations
- **Parsing Errors**: Data transformation failures
- **Service Errors**: Communication issues with dependent services
- **Recovery**: Graceful degradation and error reporting

## Performance

### Optimization Features
- **Concurrent Processing**: Parallel rule execution
- **Memory Efficiency**: Streaming data processing
- **Caching**: Rule compilation and object reflection caching

### Monitoring
- **Metrics**: Processing time and throughput metrics
- **Logging**: Detailed operation logging
- **Health Checks**: Service health monitoring

## Contributing

### Development Setup

1. Fork the repository
2. Create a feature branch
3. Make changes following Go best practices
4. Run tests and ensure coverage
5. Submit a pull request

### Code Standards
- Follow Go formatting conventions
- Include comprehensive unit tests
- Document public APIs
- Maintain backward compatibility

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.

## Related Projects

- [l8collector](https://github.com/saichler/l8collector) - Data collection framework
- [l8inventory](https://github.com/saichler/l8inventory) - Inventory management
- [l8pollaris](https://github.com/saichler/l8pollaris) - Polling configuration
- [layer8](https://github.com/saichler/layer8) - Core Layer8 framework

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/saichler/l8parser).