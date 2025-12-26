# Layer 8 Model Agnostic Parsing (L8Parser)

A powerful model-agnostic parsing service for processing l8collector payloads in the Layer8 ecosystem. L8Parser provides a sophisticated rule-based parsing framework that processes collected data from network devices and transforms it into structured objects with advanced PropertyId injection and protocol buffer integration.

## Overview

L8Parser is part of the Layer8 network management ecosystem and serves as the data transformation layer between data collection and inventory management. It processes raw data collected by l8collector and applies configurable parsing rules to extract meaningful information that can be stored in the inventory.

## Features

- **Rule-Based Parsing**: Flexible parsing engine with pluggable rules
- **SNMP Support**: Built-in support for SNMP data parsing with MIB-2 compatibility and multi-vendor device support
- **Kubernetes Integration**: Complete K8s cluster monitoring with pods, nodes, deployments, and services
- **Protocol Buffer Support**: Enhanced data serialization with protobuf definitions for inventory structures
- **Protocol Agnostic**: Extensible architecture supporting multiple collection protocols
- **Real-time Processing**: Event-driven processing of collector jobs
- **Service Integration**: Seamless integration with Layer8 service architecture
- **Comprehensive Testing**: Extensive test suite with physical device simulations and persistent test data
- **Documentation**: Complete vendor limitation documentation and polling configuration guides

## Architecture

### Core Components

- **Parser Service** (`go/parser/service/Parser.go`): Main parsing engine that processes jobs
- **Parsing Service** (`go/parser/service/ParsingService.go`): Service wrapper implementing Layer8 service interface
- **Parsing Center** (`go/parser/service/ParsingCenter.go`): Job completion handler and inventory integration
- **Rule Engine** (`go/parser/rules/`): Collection of parsing rules for data transformation

### Project Structure

```
l8parser/
├── go/                           # Go module root
│   ├── parser/
│   │   ├── boot/                 # Boot configuration and vendor-specific parsers
│   │   │   ├── K8s.go           # Kubernetes resource monitoring
│   │   │   ├── SNMP.go          # Common SNMP utilities
│   │   │   ├── cisco.go         # Cisco device configurations
│   │   │   ├── juniper.go       # Juniper device configurations
│   │   │   └── ...              # Other vendor configurations
│   │   ├── service/             # Core parsing services
│   │   └── rules/               # Parsing rule implementations
│   ├── tests/                   # Comprehensive test suite
│   │   ├── jobsPersistency/     # Persistent test data
│   │   ├── clusters.json        # Cluster test configurations
│   │   └── *.go                 # Test files
│   ├── go.mod                   # Go module definition
│   └── test.sh                  # Test execution script
├── proto/                       # Protocol buffer definitions
│   └── inventory.proto          # Device inventory structures
├── web.html                     # Web interface
├── notsupported.md             # Vendor limitation documentation
└── README.md                    # This file
```

### Parsing Rules

The parser supports multiple rule types:

1. **Contains** (`Contains.go`): Searches for specific text patterns in data
2. **Set** (`Set.go`): Directly sets values to object properties
3. **ToTable** (`ToTable.go`): Converts data into tabular format
4. **TableToMap** (`TableToMap.go`): Transforms tables into key-value maps

## Installation

### Prerequisites

- Go 1.24.0 or higher (current: Go 1.24.9 toolchain)
- Access to Layer8 ecosystem modules
- Protocol Buffer support for enhanced data structures
- Kubernetes API access for K8s monitoring features

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

## Recent Updates (December 2025)

### Latest Enhancements

- **Apache License 2.0**: Project now fully licensed under Apache License 2.0 with comprehensive copyright headers
- **Full Kubernetes Support**: Complete K8s resource monitoring with pods, nodes, deployments, and services
- **Protocol Buffer Integration**: Added inventory.proto for enhanced data serialization and type safety
- **Extended Testing Framework**: New test files for GetStringInput, GetValueInput, and cluster testing
- **Documentation Improvements**: Added comprehensive documentation for unsupported attributes by vendor
- **Dependency Updates**: Updated all Layer8 dependencies to latest versions (December 2025)
- **Go Toolchain Update**: Updated to Go 1.24.9 toolchain for enhanced performance and security

### October 2025 Updates

- **Kubernetes Integration**: Enhanced k8s interval configuration for cloud-native deployments
- **Performance Optimizations**: Enlarged timeout settings and improved status cadence handling
- **Test Framework Improvements**: Comprehensive testing enhancements with better physical device testing
- **Repository Restructuring**: Updated repository organization and dependency management

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

## Kubernetes Integration

### Overview

L8Parser now includes comprehensive Kubernetes cluster monitoring and parsing capabilities through the new K8s module. This integration enables real-time monitoring of Kubernetes resources and workloads.

### Supported Kubernetes Resources

The parser can collect and process data from the following Kubernetes resources:

- **Nodes**: Cluster node information, status, and specifications
- **Pods**: Pod details, status, containers, and resource usage
- **Deployments**: Deployment configurations and rollout status
- **StatefulSets**: StatefulSet specifications and pod management
- **DaemonSets**: DaemonSet deployments across cluster nodes
- **Services**: Service endpoints, ports, and load balancing
- **Namespaces**: Namespace isolation and resource organization
- **Network Policies**: Network security rules and traffic control

### Kubernetes Polling Configuration

```go
func CreateK8sBootPolls() *l8tpollaris.L8Pollaris {
    k8sPollaris := &l8tpollaris.L8Pollaris{
        Name:   "kubernetes",
        Groups: []string{common.BOOT_STAGE_00},
        Polling: map[string]*l8tpollaris.L8Poll{
            // Node monitoring
            "nodes": {
                What:      "get nodes -o wide",
                Operation: l8tpollaris.L8C_Operation_L8C_Table,
            },
            // Pod monitoring
            "pods": {
                What:      "get pods -A -o wide",
                Operation: l8tpollaris.L8C_Operation_L8C_Table,
            },
        },
    }
    return k8sPollaris
}
```

### Features

- **Real-time Monitoring**: Continuous polling of Kubernetes resources
- **Log Collection**: Automated log retrieval from pods and containers
- **Resource Details**: Deep inspection of resource configurations
- **Multi-cluster Support**: Monitoring across multiple Kubernetes clusters
- **kubectl Integration**: Native kubectl command execution and parsing

### Configuration

The Kubernetes module is configured through the boot stage initialization and supports various polling intervals and detail levels for different resource types.

## Protocol Buffer Integration

### Enhanced Data Structures

The project now includes Protocol Buffer definitions for improved data serialization and type safety:

- **inventory.proto**: Defines comprehensive network device inventory structures
- **NetworkDevice**: Core device representation with equipment info, physical components, and logical interfaces
- **EquipmentInfo**: Detailed device metadata including vendor, model, firmware, and location
- **DeviceType**: Enumerated device classifications (routers, switches, firewalls, servers)
- **NetworkTopology**: Topology discovery and network relationship mapping

### Proto Structure Example

```protobuf
message NetworkDevice {
  string id = 1;
  EquipmentInfo equipmentinfo = 2;
  map<string, Physical> physicals = 3;
  map<string, Logical> logicals = 4;
  NetworkTopology topology = 5;
  repeated NetworkLink network_links = 6;
  NetworkHealth network_health = 7;
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
- Physical device testing with persistent job data
- Cluster configuration testing with clusters.json
- Lab environment simulation with .kubeadm-lab and lab.conf

### New Test Components

- **GetStringInput_test.go**: Comprehensive string parsing validation
- **GetValueInput_test.go**: Value extraction and type conversion testing
- **ClusterTest_test.go**: Kubernetes cluster configuration testing
- **PhysicalTestFromPersistency_test.go**: Testing with real device data from jobsPersistency directory
- **DevicesParsing_test.go**: Vendor-specific device parsing validation

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

## Documentation

### Unsupported Attributes Documentation

The project includes comprehensive documentation of vendor-specific limitations in `notsupported.md`:

- **Vendor-Specific Limitations**: Detailed list of attributes not available via SNMP/SSH for each vendor
- **Protocol Limitations**: SNMP protocol constraints and version dependencies
- **Alternative Collection Methods**: NETCONF and REST API options for extended data collection
- **Workarounds**: Practical solutions for gathering unsupported metrics

Key unsupported areas include:
- Real-time network health metrics requiring continuous monitoring
- Hardware-specific details not exposed via SNMP (CPU architecture, memory frequency)
- Topology discovery requiring network-wide visibility
- Historical performance data requiring time-series collection

### Available Documentation

- **README.md**: Project overview and usage instructions
- **notsupported.md**: Comprehensive list of vendor-specific limitations
- **COMPREHENSIVE_POLLING_DOCUMENTATION.md**: Detailed polling configuration guide
- **web.html**: Interactive web interface for parsing service visualization

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

© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Related Projects

- [l8collector](https://github.com/saichler/l8collector) - Data collection framework
- [l8inventory](https://github.com/saichler/l8inventory) - Inventory management
- [l8pollaris](https://github.com/saichler/l8pollaris) - Polling configuration
- [layer8](https://github.com/saichler/layer8) - Core Layer8 framework

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/saichler/l8parser).