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

L8Parser includes built-in support for SNMP data parsing with pre-configured MIB-2 rules:

### System Information
- **Vendor Detection**: Automatically detects device vendors (Cisco, Ubuntu Linux, etc.)
- **System Name**: Extracts system names from SNMP sysName OID
- **System Description**: Processes system description strings

### Example SNMP Configuration

```go
func CreateSNMPBootPolls() *types.Pollaris {
    snmpPolaris := &types.Pollaris{
        Name:   "mib2",
        Groups: []string{"boot"},
        Polling: map[string]*types.Poll{
            "systemMib": {
                What:      ".1.3.6.1.2.1.1",
                Operation: types.Operation_OMap,
                Attributes: []*types.Attribute{
                    createVendorAttribute(),
                    createSysNameAttribute(),
                },
            },
        },
    }
    return snmpPolaris
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